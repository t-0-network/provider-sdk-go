package provider

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type ServerOption func(*serverOptions)

type serverOptions struct {
	addr              string
	readTimeout       time.Duration
	writeTimeout      time.Duration
	readHeaderTimeout time.Duration
	tlsConfig         *tls.Config
	shutdownTimeout   time.Duration // applies only to started server
}

func WithAddr(addr string) ServerOption {
	return func(opts *serverOptions) {
		opts.addr = addr
	}
}

func WithReadTimeout(timeout time.Duration) ServerOption {
	return func(opts *serverOptions) {
		opts.readTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) ServerOption {
	return func(opts *serverOptions) {
		opts.writeTimeout = timeout
	}
}

func WithReadHeaderTimeout(timeout time.Duration) ServerOption {
	return func(opts *serverOptions) {
		opts.readHeaderTimeout = timeout
	}
}

func WithTLSConfig(tlsConfig *tls.Config) ServerOption {
	return func(opts *serverOptions) {
		opts.tlsConfig = tlsConfig
	}
}

var defaultServerOptions = serverOptions{
	addr:              ":8080",
	readTimeout:       10 * time.Second,
	writeTimeout:      10 * time.Second,
	readHeaderTimeout: 10 * time.Second,
	tlsConfig:         nil,
	shutdownTimeout:   15 * time.Second, // default shutdown timeout for started server
}

type ServerShutdownFn func(ctx context.Context) error

// NewServer returns a ready-to-use *http.Server with the provided handler registered.
func NewServer(handler http.Handler, serverOptions ...ServerOption) *http.Server {
	opts := defaultServerOptions
	for _, opt := range serverOptions {
		opt(&opts)
	}

	server := http.Server{
		Addr:              opts.addr,
		ReadTimeout:       opts.readTimeout,
		ReadHeaderTimeout: opts.readHeaderTimeout,
		WriteTimeout:      opts.writeTimeout,
		Handler:           h2c.NewHandler(handler, &http2.Server{}),
	}

	return &server
}

// StartServer creates and starts a new HTTP server with the provided handler
// registered and ready to handle requests.
func StartServer(handler http.Handler, serverOptions ...ServerOption) ServerShutdownFn {
	server := NewServer(handler, serverOptions...)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		listener, err := net.Listen("tcp", server.Addr)
		if err != nil {
			panic(fmt.Sprintf("failed to create listener: %v", err))
		}
		wg.Done()

		// At this point, server is ready to accept connections
		if err := server.Serve(listener); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Sprintf("failed to start provider server: %v", err))
		}
	}()

	// Wait for server to be ready
	wg.Wait()

	serverShutdown := func(ctx context.Context) error {
		timeoutCtx, cancel := context.WithTimeout(ctx, defaultServerOptions.shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(timeoutCtx); err != nil {
			return fmt.Errorf("http server shutdown: %w", err)
		}

		return nil
	}

	return serverShutdown
}
