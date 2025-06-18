package service

import (
	"errors"

	"connectrpc.com/connect"
)

const (
	defaultMaxBodySize = 1024 * 1024 // 1 MB
)

var ErrPrivateKeyNotSet = errors.New("private key not set")

type providerHandlerOptions struct {
	networkHexedPublicKey        string
	verifySignatureFn            verifySignature
	verifySignatureMaxBodySize   int64
	connectHandlerOptions        []connect.HandlerOption
	disableSignatureVerification bool
}

var defaultProviderHandlerOptions = &providerHandlerOptions{
	verifySignatureMaxBodySize:   defaultMaxBodySize,
	connectHandlerOptions:        []connect.HandlerOption{},
	disableSignatureVerification: false,
	networkHexedPublicKey:        "",
	verifySignatureFn:            nil,
}

type HandlerOption func(*providerHandlerOptions)

func WithNetworkPublicKey(publicKey string) HandlerOption {
	return func(h *providerHandlerOptions) {
		h.networkHexedPublicKey = publicKey
	}
}

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
