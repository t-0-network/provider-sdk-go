package error

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/pkg/internal/crypto"
)

const (
	SignatureHeader = "X-Signature"
	PublicKeyHeader = "X-Public-Key"
)

type (
	VerifySignatureMiddleware  func(http.Handler) http.Handler
	VerifySignatureMaxBodySize int64
)

type SignatureError struct {
	ConnectCode connect.Code
	Message     string
}

type SignatureErrorContextKey struct{}

func GetSignatureErrorFromContext(ctx context.Context) (*SignatureError, bool) {
	sigErr, ok := ctx.Value(SignatureErrorContextKey{}).(*SignatureError)
	return sigErr, ok
}

func NewSignatureVerifierMiddleware(
	verifySignature VerifySignature,
	maxBodySizeOpt VerifySignatureMaxBodySize,
) VerifySignatureMiddleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			setErrorAndContinue := func(req *http.Request, statusCode connect.Code, message string) {
				errObj := &SignatureError{
					ConnectCode: statusCode,
					Message:     message,
				}
				ctx := context.WithValue(req.Context(), SignatureErrorContextKey{}, errObj)
				handler.ServeHTTP(writer, req.WithContext(ctx))
			}

			publicKey, err := parseRequiredHexedHeader(PublicKeyHeader, req.Header)
			if err != nil {
				setErrorAndContinue(req, connect.CodeInvalidArgument, err.Error())
				return
			}

			signature, err := parseRequiredHexedHeader(SignatureHeader, req.Header)
			if err != nil {
				setErrorAndContinue(req, connect.CodeInvalidArgument, err.Error())
				return
			}

			body, err := readBodyWithCap(req, int64(maxBodySizeOpt))
			if err != nil {
				setErrorAndContinue(req, connect.CodeInvalidArgument, err.Error())
				return
			}

			// Restore body for downstream handlers
			_ = req.Body.Close()
			req.Body = io.NopCloser(bytes.NewReader(body))

			if err := verifySignature(publicKey, body, signature); err != nil {
				setErrorAndContinue(req, connect.CodeUnauthenticated, err.Error())
				return
			}

			ctx := context.WithValue(req.Context(), SignatureErrorContextKey{}, (*SignatureError)(nil))
			handler.ServeHTTP(writer, req.WithContext(ctx))
		})
	}
}

func parseRequiredHexedHeader(headerName string, headers http.Header) ([]byte, error) {
	encodedHeader := headers.Get(headerName)
	if encodedHeader == "" {
		return nil, fmt.Errorf("%w: %s", ErrMissingRequiredHeader, headerName)
	}

	if len(encodedHeader) < 2 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidHeaderEncoding, headerName)
	}

	decodedHeader, err := hex.DecodeString(encodedHeader[2:])
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidHeaderEncoding, headerName)
	}

	return decodedHeader, nil
}

func readBodyWithCap(r *http.Request, cap int64) ([]byte, error) {
	contentLenHeader := r.Header.Get("Content-Length")
	contentLen, err := strconv.ParseInt(contentLenHeader, 10, 64)
	if err == nil && contentLen > cap {
		return nil, fmt.Errorf("max payload size of %d bytes exceeded", cap)
	}
	// The Content-Length header is optional, and we shouldn't trust it anyway. It's just an optimization.
	// Let's also put a cap while reading the body to avoid memory overload.
	var body bytes.Buffer
	w := bufio.NewWriter(&body)
	_, err = io.CopyN(w, r.Body, cap)
	if err == nil || !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("max payload size of %d bytes exceeded", cap)
	}

	return body.Bytes(), nil
}

type VerifySignature func(publicKey []byte, message []byte, signature []byte) error

func newVerifyEthereumSignature() VerifySignature {
	return func(publicKey []byte, message []byte, signature []byte) error {
		if len(signature) < 64 || len(signature) > 65 {
			return ErrInvalidSignature
		}

		digestHash := crypto.LegacyKeccak256(message)

		pk, err := crypto.GetPublicKeyFromBytes(publicKey)
		if err != nil {
			return fmt.Errorf("invalid public key: %w", err)
		}

		if !crypto.VerifySignature(pk, digestHash, signature[:64]) {
			return ErrSignatureVerificationFailed
		}

		return nil
	}
}
