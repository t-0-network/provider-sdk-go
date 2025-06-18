package service

import (
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/pkg/gen/proto/network/networkconnect"
	"github.com/t-0-network/provider-sdk-go/pkg/internal/crypto"
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

		networkPublicKey, err := crypto.HexToECDSAPublicKey(handler.networkHexedPublicKey)
		if err != nil {
			return nil, fmt.Errorf("invalid network public key: %w", err)
		}

		handler.verifySignatureFn = newVerifyEthereumSignature(networkPublicKey)
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
