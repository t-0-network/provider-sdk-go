package service

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/pkg/gen/proto/network/networkconnect"
)

func NewProviderHandler(
	providerHandler networkconnect.ProviderServiceHandler,
	option ...HandlerOption,
) (http.Handler, error) {
	handler := defaultProviderHandlerOptions
	for _, opt := range option {
		opt(handler)
	}

	if handler.verifySignatureFn == nil {
		if handler.networkHexedPublicKey == "" {
			return nil, ErrNetworkPublicKeyIsRequired
		}

		verifySignatureFn, err := newVerifyEthereumSignature(handler.networkHexedPublicKey)
		if err != nil {
			return nil, err
		}
		handler.verifySignatureFn = verifySignatureFn
	}

	connectHandlerOpts := append([]connect.HandlerOption{
		connect.WithInterceptors(signatureErrorInterceptor()),
	}, handler.connectHandlerOptions...)

	path, provideServiceHandler := networkconnect.NewProviderServiceHandler(
		providerHandler,
		connectHandlerOpts...,
	)

	mux := http.NewServeMux()
	mux.Handle(path, provideServiceHandler)

	return newSignatureVerifierMiddleware(
		handler.verifySignatureFn, handler.verifySignatureMaxBodySize,
	)(mux), nil
}
