package provider

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/api/gen/proto/tzero/v1/payment/paymentconnect"
	paymentintent "github.com/t-0-network/provider-sdk-go/api/gen/proto/tzero/v1/payment_intent/provider/providerconnect"
)

type BuildHandler func(options ...connect.HandlerOption) (path string, handler http.Handler)

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
	buildHandler BuildHandler,
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

	path, providerServiceHandler := buildHandler(connectHandlerOpts...)

	mux := http.NewServeMux()
	mux.Handle(path, providerServiceHandler)

	return newSignatureVerifierMiddleware(
		handler.verifySignatureFn, handler.verifySignatureMaxBodySize,
	)(mux), nil
}

func WithProviderServiceHandler(
	providerHandler paymentconnect.ProviderServiceHandler,
) BuildHandler {
	return func(options ...connect.HandlerOption) (string, http.Handler) {
		path, handler := paymentconnect.NewProviderServiceHandler(providerHandler, options...)
		return path, handler
	}
}

func WithPaymentIntentProviderServiceHandler(
	paymentIntentProviderHandler paymentintent.ProviderServiceHandler,
) BuildHandler {
	return func(options ...connect.HandlerOption) (string, http.Handler) {
		path, handler := paymentintent.NewProviderServiceHandler(paymentIntentProviderHandler, options...)
		return path, handler
	}
}
