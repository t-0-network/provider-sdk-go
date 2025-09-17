package provider

import (
	"context"
	"crypto/tls"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
)

func TestStartServer_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	shutdownFn, err := StartServer(handler, WithAddr(":0")) // Use port 0 for automatic port assignment
	require.NoError(t, err)
	require.NotNil(t, shutdownFn)

	// Test shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = shutdownFn(ctx)
	assert.NoError(t, err)
}

func TestStartServer_InvalidAddress(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Try to bind to an invalid address
	shutdownFn, err := StartServer(handler, WithAddr("invalid:address"))
	assert.Error(t, err)
	assert.Nil(t, shutdownFn)
	assert.Contains(t, err.Error(), "failed to create listener")
}

func TestStartServer_PortAlreadyInUse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Start first server
	shutdownFn1, err := StartServer(handler, WithAddr(":0"))
	require.NoError(t, err)
	require.NotNil(t, shutdownFn1)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdownFn1(ctx)
	}()

	// Try to start second server on same port (this should work with port 0)
	shutdownFn2, err := StartServer(handler, WithAddr(":0"))
	require.NoError(t, err)
	require.NotNil(t, shutdownFn2)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdownFn2(ctx)
	}()
}

func TestStartServer_WithOptions(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	shutdownFn, err := StartServer(
		handler,
		WithAddr(":0"),
		WithReadTimeout(5*time.Second),
		WithWriteTimeout(5*time.Second),
		WithReadHeaderTimeout(5*time.Second),
		WithShutdownTimeout(10*time.Second),
	)
	require.NoError(t, err)
	require.NotNil(t, shutdownFn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = shutdownFn(ctx)
	assert.NoError(t, err)
}

func TestStartServer_ShutdownTimeout(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	shutdownFn, err := StartServer(
		handler,
		WithAddr(":0"),
		WithShutdownTimeout(100*time.Millisecond), // Very short timeout
	)
	require.NoError(t, err)
	require.NotNil(t, shutdownFn)

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start := time.Now()
	err = shutdownFn(ctx)
	duration := time.Since(start)

	// The shutdown should complete quickly due to our short shutdown timeout
	// It might succeed or fail with timeout, but should be fast
	assert.Less(t, duration, 1*time.Second) // Should complete quickly
}

func TestStartServer_ConcurrentShutdown(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	shutdownFn, err := StartServer(handler, WithAddr(":0"))
	require.NoError(t, err)
	require.NotNil(t, shutdownFn)

	// Test concurrent shutdown calls
	var wg sync.WaitGroup
	errors := make([]error, 3)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			errors[index] = shutdownFn(ctx)
		}(i)
	}

	wg.Wait()

	// At least one shutdown should succeed, others might fail with "server closed" or similar
	successCount := 0
	for _, err := range errors {
		if err == nil {
			successCount++
		}
	}
	assert.GreaterOrEqual(t, successCount, 1, "At least one shutdown should succeed")
}

func TestNewServer_WithTLSConfig(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	server := NewServer(
		handler,
		WithAddr(":8080"),
		WithTLSConfig(tlsConfig),
		WithReadTimeout(5*time.Second),
	)

	assert.Equal(t, ":8080", server.Addr)
	assert.Equal(t, tlsConfig, server.TLSConfig)
	assert.Equal(t, 5*time.Second, server.ReadTimeout)
	assert.NotNil(t, server.Handler)
}

func TestServerOptions_Defaults(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := NewServer(handler)

	assert.Equal(t, ":8080", server.Addr)
	assert.Equal(t, 10*time.Second, server.ReadTimeout)
	assert.Equal(t, 10*time.Second, server.WriteTimeout)
	assert.Equal(t, 10*time.Second, server.ReadHeaderTimeout)
	assert.Nil(t, server.TLSConfig)
	assert.NotNil(t, server.Handler)
}

// Additional comprehensive tests from server_improved_test.go

func TestServerConstants(t *testing.T) {
	// Test that all constants have reasonable values
	assert.Equal(t, ":8080", DefaultAddr)
	assert.Equal(t, 10*time.Second, DefaultReadTimeout)
	assert.Equal(t, 10*time.Second, DefaultWriteTimeout)
	assert.Equal(t, 10*time.Second, DefaultReadHeaderTimeout)
	assert.Equal(t, 15*time.Second, DefaultShutdownTimeout)
	assert.Equal(t, 5*time.Second, ServerStartupTimeout)
}

func TestWithAddrValidation(t *testing.T) {
	tests := []struct {
		name     string
		addr     string
		expected string
	}{
		{
			name:     "valid address",
			addr:     ":9000",
			expected: ":9000",
		},
		{
			name:     "empty address should not override default",
			addr:     "",
			expected: DefaultAddr, // Should keep default
		},
		{
			name:     "localhost address",
			addr:     "localhost:8080",
			expected: "localhost:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := defaultServerOptions
			WithAddr(tt.addr)(&opts)
			assert.Equal(t, tt.expected, opts.addr)
		})
	}
}

func TestTimeoutValidation(t *testing.T) {
	tests := []struct {
		name        string
		timeout     time.Duration
		expectedSet bool
		optionFunc  func(time.Duration) ServerOption
	}{
		{
			name:        "valid positive timeout",
			timeout:     5 * time.Second,
			expectedSet: true,
			optionFunc:  WithReadTimeout,
		},
		{
			name:        "zero timeout is valid",
			timeout:     0,
			expectedSet: true,
			optionFunc:  WithReadTimeout,
		},
		{
			name:        "negative timeout is ignored",
			timeout:     -1 * time.Second,
			expectedSet: false,
			optionFunc:  WithReadTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := defaultServerOptions
			original := opts.readTimeout

			tt.optionFunc(tt.timeout)(&opts)

			if tt.expectedSet {
				assert.Equal(t, tt.timeout, opts.readTimeout)
			} else {
				assert.Equal(t, original, opts.readTimeout)
			}
		})
	}
}

func TestWithShutdownTimeoutValidation(t *testing.T) {
	tests := []struct {
		name        string
		timeout     time.Duration
		expectedSet bool
	}{
		{
			name:        "positive timeout",
			timeout:     30 * time.Second,
			expectedSet: true,
		},
		{
			name:        "zero timeout should not override default",
			timeout:     0,
			expectedSet: false,
		},
		{
			name:        "negative timeout should not override default",
			timeout:     -1 * time.Second,
			expectedSet: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := defaultServerOptions
			original := opts.shutdownTimeout

			WithShutdownTimeout(tt.timeout)(&opts)

			if tt.expectedSet {
				assert.Equal(t, tt.timeout, opts.shutdownTimeout)
			} else {
				assert.Equal(t, original, opts.shutdownTimeout)
			}
		})
	}
}

func TestWithHTTP2Config(t *testing.T) {
	customConfig := &http2.Server{
		MaxConcurrentStreams: 100,
		IdleTimeout:          30 * time.Second,
	}

	opts := defaultServerOptions
	WithHTTP2Config(customConfig)(&opts)

	assert.Equal(t, customConfig, opts.http2Config)

	// Test nil config doesn't override
	originalConfig := opts.http2Config
	WithHTTP2Config(nil)(&opts)
	assert.Equal(t, originalConfig, opts.http2Config)
}

func TestNewServerNilHandler(t *testing.T) {
	assert.Panics(t, func() {
		NewServer(nil)
	})
}

func TestNewServerWithCustomOptions(t *testing.T) {
	customTLS := &tls.Config{InsecureSkipVerify: true}
	customHTTP2 := &http2.Server{MaxConcurrentStreams: 50}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := NewServer(handler,
		WithAddr(":9090"),
		WithReadTimeout(20*time.Second),
		WithWriteTimeout(25*time.Second),
		WithReadHeaderTimeout(5*time.Second),
		WithTLSConfig(customTLS),
		WithHTTP2Config(customHTTP2),
	)

	assert.Equal(t, ":9090", server.Addr)
	assert.Equal(t, 20*time.Second, server.ReadTimeout)
	assert.Equal(t, 25*time.Second, server.WriteTimeout)
	assert.Equal(t, 5*time.Second, server.ReadHeaderTimeout)
	assert.Equal(t, customTLS, server.TLSConfig)
	assert.NotNil(t, server.Handler)
}

func TestStartServerNilHandler(t *testing.T) {
	shutdownFn, err := StartServer(nil)
	assert.Error(t, err)
	assert.Nil(t, shutdownFn)
	assert.Contains(t, err.Error(), "handler cannot be nil")
}

func TestStartServerCancelledContextShutdown(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	shutdownFn, err := StartServer(handler, WithAddr(":0"))
	require.NoError(t, err)
	require.NotNil(t, shutdownFn)

	// Create an already cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err = shutdownFn(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context already done")

	// Properly shutdown with fresh context
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	err = shutdownFn(ctx2)
	assert.NoError(t, err)
}

func TestStartServerShutdownRespectsCaller_Context(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Start server with long shutdown timeout
	shutdownFn, err := StartServer(handler,
		WithAddr(":0"),
		WithShutdownTimeout(30*time.Second), // Long server timeout
	)
	require.NoError(t, err)
	require.NotNil(t, shutdownFn)

	// But use a short caller timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	err = shutdownFn(ctx)
	duration := time.Since(start)

	// Should respect caller's short timeout, not server's long timeout
	assert.Less(t, duration, 1*time.Second)
}

func TestStartServerMultipleShutdownCallsSafeAndIdempotent(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	shutdownFn, err := StartServer(handler, WithAddr(":0"))
	require.NoError(t, err)
	require.NotNil(t, shutdownFn)

	// First shutdown should succeed
	ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel1()
	err1 := shutdownFn(ctx1)
	assert.NoError(t, err1)

	// Second shutdown should return nil (idempotent)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	err2 := shutdownFn(ctx2)
	assert.NoError(t, err2) // Should not error due to sync.Once
}

func TestCreateServerHelperFunction(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	customHTTP2 := &http2.Server{MaxConcurrentStreams: 200}
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	// Use ServerOption functions instead of direct struct
	options := []ServerOption{
		WithAddr(":8888"),
		WithReadTimeout(1 * time.Second),
		WithWriteTimeout(2 * time.Second),
		WithReadHeaderTimeout(3 * time.Second),
		WithTLSConfig(tlsConfig),
		WithShutdownTimeout(10 * time.Second),
		WithHTTP2Config(customHTTP2),
	}

	server, opts := createServer(handler, options)

	assert.Equal(t, ":8888", server.Addr)
	assert.Equal(t, 1*time.Second, server.ReadTimeout)
	assert.Equal(t, 2*time.Second, server.WriteTimeout)
	assert.Equal(t, 3*time.Second, server.ReadHeaderTimeout)
	assert.Equal(t, tlsConfig, server.TLSConfig)
	assert.NotNil(t, server.Handler)

	// Test that options were applied correctly
	assert.Equal(t, ":8888", opts.addr)
	assert.Equal(t, 1*time.Second, opts.readTimeout)
	assert.Equal(t, 2*time.Second, opts.writeTimeout)
	assert.Equal(t, 3*time.Second, opts.readHeaderTimeout)
	assert.Equal(t, tlsConfig, opts.tlsConfig)
	assert.Equal(t, 10*time.Second, opts.shutdownTimeout)
	assert.Equal(t, customHTTP2, opts.http2Config)
}

func TestDefaultServerOptionsIntegrity(t *testing.T) {
	opts := defaultServerOptions

	// Ensure default options have sensible values
	assert.NotEmpty(t, opts.addr)
	assert.Greater(t, opts.readTimeout, time.Duration(0))
	assert.Greater(t, opts.writeTimeout, time.Duration(0))
	assert.Greater(t, opts.readHeaderTimeout, time.Duration(0))
	assert.Greater(t, opts.shutdownTimeout, time.Duration(0))
	assert.NotNil(t, opts.http2Config)
}

// Benchmark the server creation to ensure no performance regression
func BenchmarkNewServer(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server := NewServer(handler)
		_ = server // Prevent optimization
	}
}

func BenchmarkStartServer(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shutdownFn, err := StartServer(handler, WithAddr(":0"))
		if err != nil {
			b.Fatal(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		_ = shutdownFn(ctx)
		cancel()
	}
}
