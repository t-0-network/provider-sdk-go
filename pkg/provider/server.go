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

// Default timeout values
const (
	DefaultAddr              = ":8080"
	DefaultReadTimeout       = 10 * time.Second
	DefaultWriteTimeout      = 10 * time.Second
	DefaultReadHeaderTimeout = 10 * time.Second
	DefaultShutdownTimeout   = 15 * time.Second
	ServerStartupTimeout     = 5 * time.Second
)

// ServerOption configures server options using the functional options pattern
type ServerOption func(*serverOptions)

type serverOptions struct {
	addr              string
	readTimeout       time.Duration
	writeTimeout      time.Duration
	readHeaderTimeout time.Duration
	tlsConfig         *tls.Config
	shutdownTimeout   time.Duration // applies only to started server
	http2Config       *http2.Server
}

// WithAddr sets the server's address to listen on (host:port format)
// If an empty string is provided, the default ":8080" will be used
func WithAddr(addr string) ServerOption {
	return func(opts *serverOptions) {
		if addr != "" {
			opts.addr = addr
		}
	}
}

// WithReadTimeout sets the maximum duration for reading the entire request
// including the body. A timeout of 0 means no timeout.
// Negative timeouts are ignored and the default will be used.
func WithReadTimeout(timeout time.Duration) ServerOption {
	return func(opts *serverOptions) {
		if timeout < 0 {
			// Negative timeouts don't make sense, ignore and keep default
			return
		}
		opts.readTimeout = timeout
	}
}

// WithWriteTimeout sets the maximum duration before timing out writes of the response
// A timeout of 0 means no timeout.
// Negative timeouts are ignored and the default will be used.
func WithWriteTimeout(timeout time.Duration) ServerOption {
	return func(opts *serverOptions) {
		if timeout < 0 {
			// Negative timeouts don't make sense, ignore and keep default
			return
		}
		opts.writeTimeout = timeout
	}
}

// WithReadHeaderTimeout sets the amount of time allowed to read request headers
// A timeout of 0 means no timeout.
// Negative timeouts are ignored and the default will be used.
func WithReadHeaderTimeout(timeout time.Duration) ServerOption {
	return func(opts *serverOptions) {
		if timeout < 0 {
			// Negative timeouts don't make sense, ignore and keep default
			return
		}
		opts.readHeaderTimeout = timeout
	}
}

// WithTLSConfig sets the TLS configuration for the server
func WithTLSConfig(tlsConfig *tls.Config) ServerOption {
	return func(opts *serverOptions) {
		opts.tlsConfig = tlsConfig
	}
}

// WithShutdownTimeout sets the maximum duration to wait for the server to shutdown gracefully
// If the timeout is <= 0, the default timeout will be used.
func WithShutdownTimeout(timeout time.Duration) ServerOption {
	return func(opts *serverOptions) {
		if timeout > 0 {
			opts.shutdownTimeout = timeout
		}
	}
}

// WithHTTP2Config sets custom HTTP/2 server configuration
func WithHTTP2Config(config *http2.Server) ServerOption {
	return func(opts *serverOptions) {
		if config != nil {
			opts.http2Config = config
		}
	}
}

var defaultServerOptions = serverOptions{
	addr:              DefaultAddr,
	readTimeout:       DefaultReadTimeout,
	writeTimeout:      DefaultWriteTimeout,
	readHeaderTimeout: DefaultReadHeaderTimeout,
	tlsConfig:         nil,
	shutdownTimeout:   DefaultShutdownTimeout,
	http2Config:       &http2.Server{},
}

// ServerShutdownFn is a function that gracefully shuts down the server.
// It blocks until the server is shut down or the context is cancelled.
// It is safe to call concurrently, but only the first call is guaranteed to succeed.
type ServerShutdownFn func(ctx context.Context) error

// NewServer returns a ready-to-use *http.Server with the provided handler registered.
// The server is not started - you need to call ListenAndServe or similar methods.
func NewServer(handler http.Handler, serverOptions ...ServerOption) *http.Server {
	if handler == nil {
		panic("handler cannot be nil")
	}

	server, _ := createServer(handler, serverOptions)
	return server
}

// StartServer creates and starts a new HTTP server with the provided handler.
//
// The server starts asynchronously and this function returns immediately after
// confirming the server is ready to accept connections or after ServerStartupTimeout.
//
// Example:
//
//	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	    w.WriteHeader(http.StatusOK)
//	})
//
//	shutdown, err := StartServer(handler,
//	    WithAddr(":8080"),
//	    WithReadTimeout(30*time.Second),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer shutdown(context.Background())
//
// Returns:
//   - ServerShutdownFn: Safe for concurrent use, only first call performs shutdown
//   - error: Non-nil if server failed to start or bind to address
func StartServer(handler http.Handler, serverOptions ...ServerOption) (ServerShutdownFn, error) {
	if handler == nil {
		return nil, fmt.Errorf("handler cannot be nil")
	}

	server, opts := createServer(handler, serverOptions)
	listener, err := createListener(server.Addr)
	if err != nil {
		return nil, err
	}

	// Channels to communicate server startup status
	startupErr := make(chan error, 1)
	startupReady := make(chan struct{})

	// Wait group for graceful shutdown
	var wg sync.WaitGroup
	// Once to ensure server shutdown is only executed once
	var shutdownOnce sync.Once

	wg.Add(1)

	go func() {
		defer wg.Done()

		// Signal that server is ready to accept connections
		close(startupReady)

		var err error
		if opts.tlsConfig != nil {
			err = server.ServeTLS(listener, "", "")
		} else {
			err = server.Serve(listener)
		}

		// Only report startup errors, ignore shutdown errors
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			select {
			case startupErr <- err:
			default:
				// Startup phase is over, ignore the error
				// The server.Shutdown() call will handle cleanup
			}
		}
	}()

	// Wait for server to be ready, error out, or timeout
	select {
	case <-startupReady:
		// Server is ready to accept connections
	case err := <-startupErr:
		if err != nil {
			// Server failed to start, close listener and return error
			listener.Close()
			return nil, fmt.Errorf("failed to start provider server on %s: %w", server.Addr, err)
		}
	case <-time.After(ServerStartupTimeout):
		// Timeout waiting for server to be ready
		listener.Close()
		return nil, fmt.Errorf("server startup timeout after %v on %s", ServerStartupTimeout, server.Addr)
	}

	// Create a reusable shutdown function that can be called concurrently
	serverShutdown := func(ctx context.Context) error {
		// Check if context is already cancelled
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("shutdown context already done: %w", err)
		}

		// Variable to collect shutdown errors
		var shutdownErr error

		// Ensure shutdown only happens once
		shutdownOnce.Do(func() {
			// Determine appropriate timeout context
			var timeoutCtx context.Context
			var cancel context.CancelFunc

			// Respect both the caller's context and our shutdown timeout
			deadline, hasDeadline := ctx.Deadline()
			shutdownDeadline := time.Now().Add(opts.shutdownTimeout)

			if hasDeadline && deadline.Before(shutdownDeadline) {
				timeoutCtx, cancel = context.WithDeadline(ctx, deadline)
			} else {
				timeoutCtx, cancel = context.WithTimeout(ctx, opts.shutdownTimeout)
			}
			defer cancel()

			// Shutdown the server gracefully
			if err := server.Shutdown(timeoutCtx); err != nil {
				shutdownErr = fmt.Errorf("http server shutdown: %w", err)
			}

			// Always ensure listener is closed
			if listener != nil {
				listener.Close()
			}

			// Wait for the server goroutine to finish with timeout
			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()

			select {
			case <-done:
				// Server goroutine finished normally
			case <-timeoutCtx.Done():
				// Timeout occurred while waiting for server goroutine
				// We can't force kill the goroutine, but we can continue with shutdown
				// The goroutine will eventually finish when the server stops
			}

			// No need to check server errors - shutdown error is more important
		})

		return shutdownErr
	}

	// Close startup error channel to prevent goroutine leaks
	close(startupErr)

	return serverShutdown, nil
}

// createListener creates a TCP listener for the given address
func createListener(addr string) (net.Listener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener on %s: %w", addr, err)
	}
	return listener, nil
}

// createServer creates a new http.Server with the provided handler and options
// This is an internal helper to avoid code duplication
func createServer(handler http.Handler, options []ServerOption) (*http.Server, *serverOptions) {
	// Process options once and store them for later use
	opts := defaultServerOptions
	for _, opt := range options {
		opt(&opts)
	}

	return &http.Server{
		Addr:              opts.addr,
		ReadTimeout:       opts.readTimeout,
		ReadHeaderTimeout: opts.readHeaderTimeout,
		WriteTimeout:      opts.writeTimeout,
		TLSConfig:         opts.tlsConfig,
		Handler:           h2c.NewHandler(handler, opts.http2Config),
	}, &opts
}
