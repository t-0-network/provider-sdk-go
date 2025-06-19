package service

import (
	"context"
	"errors"

	"connectrpc.com/connect"
)

// signatureErrorInterceptor checks for a signature error in the context.
// this error is propagated from the signature verification middleware.
func signatureErrorInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			sigErr, ok := getSignatureErrorFromContext(ctx)
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
