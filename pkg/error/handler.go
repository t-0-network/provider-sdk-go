package error

import (
	"context"
	"errors"
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

	signatureVerifier := NewSignatureVerifierMiddleware(
		newVerifyEthereumSignature(),
		VerifySignatureMaxBodySize(1024*1024), // 1 MB
	)

	return signatureVerifier(mux)
}

func signatureErrorInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			sigErr, ok := GetSignatureErrorFromContext(ctx)
			if !ok {
				return nil, connect.NewError(connect.CodeInternal, ErrNoSignatureResult)
			}

			if sigErr != nil {
				return nil, connect.NewError(sigErr.ConnectCode, errors.New(sigErr.Message))
			}

			return next(ctx, req)
		}
	}
}
