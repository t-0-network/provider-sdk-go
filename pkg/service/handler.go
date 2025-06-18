package service

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/pkg/gen/proto/network/networkconnect"
)

func NewHandler(
	providerHandler networkconnect.ProviderServiceHandler,
) http.Handler {
	path, provideServiceHandler := networkconnect.NewProviderServiceHandler(
		providerHandler,
		connect.WithInterceptors(
			signatureErrorInterceptor(),
		),
	)

	mux := http.NewServeMux()
	mux.Handle(path, provideServiceHandler)

	signatureVerifier := newSignatureVerifierMiddleware(
		newVerifyEthereumSignature(),
		verifySignatureMaxBodySize(1024*1024), // 1 MB
	)

	return signatureVerifier(mux)
}
