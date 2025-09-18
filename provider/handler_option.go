package provider

import (
	"connectrpc.com/connect"
)

const (
	defaultMaxBodySize = 1024 * 1024 // 1 MB
)

type providerHandlerOptions struct {
	verifySignatureFn          VerifySignature
	verifySignatureMaxBodySize int64
	connectHandlerOptions      []connect.HandlerOption
}

func newDefaultHandlerOptions(verifySignatureFn VerifySignature) (providerHandlerOptions, error) {
	return providerHandlerOptions{
		verifySignatureMaxBodySize: defaultMaxBodySize,
		connectHandlerOptions: []connect.HandlerOption{
			connect.WithInterceptors(signatureErrorInterceptor()),
		},
		verifySignatureFn: verifySignatureFn,
	}, nil
}

type HandlerOption func(*providerHandlerOptions)

func WithVerifySignatureFn(fn VerifySignature) HandlerOption {
	return func(h *providerHandlerOptions) {
		h.verifySignatureFn = fn
	}
}

func WithConnectHandlerOptions(opts ...connect.HandlerOption) HandlerOption {
	return func(h *providerHandlerOptions) {
		h.connectHandlerOptions = append(h.connectHandlerOptions, opts...)
	}
}

// WithMaxBodySize sets the maximum allowed request body size for signature verification.
// If size is <= 0, the default size will be used.
func WithMaxBodySize(size int64) HandlerOption {
	return func(h *providerHandlerOptions) {
		if size > 0 {
			h.verifySignatureMaxBodySize = size
		}
	}
}
