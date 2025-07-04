package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/pkg/constant"
	"github.com/t-0-network/provider-sdk-go/pkg/crypto"
)

type middleware func(http.Handler) http.Handler

type SignatureError struct {
	ConnectCode connect.Code
	Message     string
}

type signatureErrorContextKey struct{}

func getSignatureErrorFromContext(ctx context.Context) (*SignatureError, bool) {
	sigErr, ok := ctx.Value(signatureErrorContextKey{}).(*SignatureError)
	return sigErr, ok
}

func newSignatureVerifierMiddleware(
	verifySignature verifySignature,
	maxBodySizeOpt int64,
) middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			setErrorAndContinue := func(req *http.Request, statusCode connect.Code, message string) {
				errObj := &SignatureError{
					ConnectCode: statusCode,
					Message:     message,
				}
				ctx := context.WithValue(req.Context(), signatureErrorContextKey{}, errObj)
				handler.ServeHTTP(writer, req.WithContext(ctx))
			}

			publicKey, err := parseRequiredHexedHeader(constant.PublicKeyHeader, req.Header)
			if err != nil {
				setErrorAndContinue(req, connect.CodeInvalidArgument, err.Error())
				return
			}

			signature, err := parseRequiredHexedHeader(constant.SignatureHeader, req.Header)
			if err != nil {
				setErrorAndContinue(req, connect.CodeInvalidArgument, err.Error())
				return
			}

			timestamp, timestampBytes, err := parseTimestamp(req.Header)
			if err != nil {
				setErrorAndContinue(req, connect.CodeInvalidArgument, err.Error())
				return
			}

			if !timesWithinDelta(timestamp, time.Now(), time.Minute) {
				setErrorAndContinue(req, connect.CodeInvalidArgument, "timestamp is outside the allowed time window")
				return
			}

			body, err := readBodyWithCap(req, maxBodySizeOpt)
			if err != nil {
				setErrorAndContinue(req, connect.CodeInvalidArgument, err.Error())
				return
			}

			// Restore body for downstream handlers
			_ = req.Body.Close()
			req.Body = io.NopCloser(bytes.NewReader(body))

			if err := verifySignature(publicKey, append(timestampBytes, body...), signature); err != nil {
				setErrorAndContinue(req, connect.CodeUnauthenticated, err.Error())
				return
			}

			ctx := context.WithValue(req.Context(), signatureErrorContextKey{}, (*SignatureError)(nil))
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

// parseTimestamp extracts the timestamp from the request headers, and returns
// the parsed time and its byte representation in little-endian format.
func parseTimestamp(headers http.Header) (time.Time, []byte, error) {
	timestampValue := headers.Get(constant.SignatureTimestampHeader)
	if timestampValue == "" {
		return time.Time{}, nil, fmt.Errorf("%w: %s", ErrMissingRequiredHeader, constant.SignatureTimestampHeader)
	}

	timestamp, err := strconv.ParseInt(timestampValue, 10, 64)
	if err != nil {
		return time.Time{}, nil, fmt.Errorf("invalid timestamp header: %s", err.Error())
	}

	timestampBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestampBytes, uint64(timestamp))

	return time.UnixMilli(timestamp), timestampBytes, nil
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

// verifySignature accepts a public key, a message, and a signature, hashes the
// message, and verifies the signature against the public key.
type verifySignature func(publicKey, message, signature []byte) error

func newVerifySignature(networkPublicKeyHexed string) (verifySignature, error) {
	networkPublicKey, err := crypto.GetPublicKeyFromHex(networkPublicKeyHexed)
	if err != nil {
		return nil, fmt.Errorf("invalid network public key: %w", err)
	}

	return func(publicKey, message, signature []byte) error {
		if len(signature) < 64 || len(signature) > 65 {
			return ErrInvalidSignature
		}

		signerPublicKey, err := crypto.GetPublicKeyFromBytes(publicKey)
		if err != nil {
			return fmt.Errorf("invalid public key: %w", err)
		}

		if !signerPublicKey.IsEqual(networkPublicKey) {
			return ErrUnknownPublicKey
		}

		digestHash := crypto.LegacyKeccak256(message)
		if !crypto.VerifySignature(signerPublicKey, digestHash, signature[:64]) {
			return ErrSignatureVerificationFailed
		}

		return nil
	}, nil
}

func timesWithinDelta(t1, t2 time.Time, delta time.Duration) bool {
	diff := t1.Sub(t2)
	if diff < 0 {
		diff = -diff
	}

	return diff <= delta
}
