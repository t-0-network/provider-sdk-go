package provider

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/api/gen/proto/network/networkconnect"
)

// T-ZERO Network Public Key, required for signature verification.
type NetworkPublicKeyHexed string

// NewProviderHandler returns a ready-to-use *http.ServeMux with the
// networkconnect.ProviderServiceHandler registered.
//
// It creates a new HTTP mux, registers the provided ProviderServiceHandler on the appropriate path,
// and returns the mux for immediate use in your HTTP server.
//
// Parameters:
//   - service: An implementation of the networkconnect.ProviderServiceHandler interface.
//
// Returns:
//   - *http.ServeMux: An HTTP mux with the provider service handler registered.
func NewProviderHandler(
	networkPublicKey NetworkPublicKeyHexed,
	providerHandler networkconnect.ProviderServiceHandler,
	option ...HandlerOption,
) (http.Handler, error) {
	handler := defaultProviderHandlerOptions
	for _, opt := range option {
		opt(&handler)
	}

	if handler.verifySignatureFn == nil {
		if networkPublicKey == "" {
			return nil, ErrNetworkPublicKeyIsRequired
		}

		verifySignatureFn, err := newVerifySignature(string(networkPublicKey))
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
