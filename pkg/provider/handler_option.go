package provider

import (
	"errors"

	"connectrpc.com/connect"
)

const (
	defaultMaxBodySize = 1024 * 1024 // 1 MB
)

var ErrNetworkPublicKeyIsRequired = errors.New("network public key is not set")

type providerHandlerOptions struct {
	verifySignatureFn          verifySignature
	verifySignatureMaxBodySize int64
	connectHandlerOptions      []connect.HandlerOption
}

var defaultProviderHandlerOptions = providerHandlerOptions{
	verifySignatureMaxBodySize: defaultMaxBodySize,
	connectHandlerOptions:      []connect.HandlerOption{},
	verifySignatureFn:          nil,
}

type HandlerOption func(*providerHandlerOptions)

func WithVerifySignatureFn(fn verifySignature) HandlerOption {
	return func(h *providerHandlerOptions) {
		h.verifySignatureFn = fn
	}
}

func WithConnectHandlerOptions(opts ...connect.HandlerOption) HandlerOption {
	return func(h *providerHandlerOptions) {
		h.connectHandlerOptions = append(h.connectHandlerOptions, opts...)
	}
}
