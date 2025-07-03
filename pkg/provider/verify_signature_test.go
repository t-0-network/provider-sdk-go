package provider

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"github.com/t-0-network/provider-sdk-go/pkg/constant"
)

func TestNewSignatureVerifierMiddleware(t *testing.T) {
	// Mock verify signature function for testing
	mockVerifySignature := func(returnError bool) verifySignature {
		return func(publicKey, message, signature []byte) error {
			if returnError {
				return fmt.Errorf("signature verification failed")
			}
			return nil
		}
	}

	// Helper to create valid headers
	createValidHeaders := func() http.Header {
		headers := http.Header{}
		headers.Set(constant.PublicKeyHeader, "0x"+hex.EncodeToString([]byte("validpublickey")))
		headers.Set(constant.SignatureHeader, "0x"+hex.EncodeToString([]byte("validsignature")))

		timestamp := time.Now().UnixMilli()
		headers.Set(constant.SignatureTimestampHeader, strconv.FormatInt(timestamp, 10))

		return headers
	}

	tests := []struct {
		name                string
		setupHeaders        func() http.Header
		requestBody         string
		verifySignatureFunc verifySignature
		expectedError       *SignatureError
	}{
		{
			name:                "valid request with all headers",
			setupHeaders:        createValidHeaders,
			requestBody:         "test body",
			verifySignatureFunc: mockVerifySignature(false),
			expectedError:       nil,
		},
		{
			name: "missing public key header",
			setupHeaders: func() http.Header {
				headers := createValidHeaders()
				headers.Del(constant.PublicKeyHeader)
				return headers
			},
			requestBody:         "test body",
			verifySignatureFunc: mockVerifySignature(false),
			expectedError: &SignatureError{
				ConnectCode: connect.CodeInvalidArgument,
				Message:     fmt.Sprintf("%s: %s", ErrMissingRequiredHeader.Error(), constant.PublicKeyHeader),
			},
		},
		{
			name: "missing signature header",
			setupHeaders: func() http.Header {
				headers := createValidHeaders()
				headers.Del(constant.SignatureHeader)
				return headers
			},
			requestBody:         "test body",
			verifySignatureFunc: mockVerifySignature(false),
			expectedError: &SignatureError{
				ConnectCode: connect.CodeInvalidArgument,
				Message:     fmt.Sprintf("%s: %s", ErrMissingRequiredHeader.Error(), constant.SignatureHeader),
			},
		},
		{
			name: "missing timestamp header",
			setupHeaders: func() http.Header {
				headers := createValidHeaders()
				headers.Del(constant.SignatureTimestampHeader)
				return headers
			},
			requestBody:         "test body",
			verifySignatureFunc: mockVerifySignature(false),
			expectedError: &SignatureError{
				ConnectCode: connect.CodeInvalidArgument,
				Message:     fmt.Sprintf("%s: %s", ErrMissingRequiredHeader.Error(), constant.SignatureTimestampHeader),
			},
		},
		{
			name: "invalid public key header encoding",
			setupHeaders: func() http.Header {
				headers := createValidHeaders()
				headers.Set(constant.PublicKeyHeader, "0xINVALIDHEX")
				return headers
			},
			requestBody:         "test body",
			verifySignatureFunc: mockVerifySignature(false),
			expectedError: &SignatureError{
				ConnectCode: connect.CodeInvalidArgument,
				Message:     fmt.Sprintf("%s: %s", ErrInvalidHeaderEncoding.Error(), constant.PublicKeyHeader),
			},
		},
		{
			name: "invalid signature header encoding",
			setupHeaders: func() http.Header {
				headers := createValidHeaders()
				headers.Set(constant.SignatureHeader, "0xINVALIDHEX")
				return headers
			},
			requestBody:         "test body",
			verifySignatureFunc: mockVerifySignature(false),
			expectedError: &SignatureError{
				ConnectCode: connect.CodeInvalidArgument,
				Message:     fmt.Sprintf("%s: %s", ErrInvalidHeaderEncoding.Error(), constant.SignatureHeader),
			},
		},
		{
			name: "public key header too short",
			setupHeaders: func() http.Header {
				headers := createValidHeaders()
				headers.Set(constant.PublicKeyHeader, "0")
				return headers
			},
			requestBody:         "test body",
			verifySignatureFunc: mockVerifySignature(false),
			expectedError: &SignatureError{
				ConnectCode: connect.CodeInvalidArgument,
				Message:     fmt.Sprintf("%s: %s", ErrInvalidHeaderEncoding.Error(), constant.PublicKeyHeader),
			},
		},
		{
			name: "invalid timestamp format",
			setupHeaders: func() http.Header {
				headers := createValidHeaders()
				headers.Set(constant.SignatureTimestampHeader, "invalid-timestamp")
				return headers
			},
			requestBody:         "test body",
			verifySignatureFunc: mockVerifySignature(false),
			expectedError: &SignatureError{
				ConnectCode: connect.CodeInvalidArgument,
				Message:     "invalid timestamp",
			},
		},
		{
			name: "timestamp outside allowed window (too old)",
			setupHeaders: func() http.Header {
				headers := createValidHeaders()
				oldTimestamp := time.Now().Add(-2 * time.Minute).UnixMilli()
				headers.Set(constant.SignatureTimestampHeader, strconv.FormatInt(oldTimestamp, 10))
				return headers
			},
			requestBody:         "test body",
			verifySignatureFunc: mockVerifySignature(false),
			expectedError: &SignatureError{
				ConnectCode: connect.CodeInvalidArgument,
				Message:     "timestamp is outside the allowed time window",
			},
		},
		{
			name: "timestamp outside allowed window (too new)",
			setupHeaders: func() http.Header {
				headers := createValidHeaders()
				futureTimestamp := time.Now().Add(2 * time.Minute).UnixMilli()
				headers.Set(constant.SignatureTimestampHeader, strconv.FormatInt(futureTimestamp, 10))
				return headers
			},
			requestBody:         "test body",
			verifySignatureFunc: mockVerifySignature(false),
			expectedError: &SignatureError{
				ConnectCode: connect.CodeInvalidArgument,
				Message:     "timestamp is outside the allowed time window",
			},
		},
		{
			name:                "signature verification fails",
			setupHeaders:        createValidHeaders,
			requestBody:         "test body",
			verifySignatureFunc: mockVerifySignature(true),
			expectedError: &SignatureError{
				ConnectCode: connect.CodeUnauthenticated,
				Message:     "signature verification failed",
			},
		},
		{
			name:                "empty body success",
			setupHeaders:        createValidHeaders,
			requestBody:         "",
			verifySignatureFunc: mockVerifySignature(false),
			expectedError:       nil,
		},
		{
			name: "timestamp exactly at boundary (valid)",
			setupHeaders: func() http.Header {
				headers := createValidHeaders()
				boundaryTimestamp := time.Now().Add(-59 * time.Second).UnixMilli()
				headers.Set(constant.SignatureTimestampHeader, strconv.FormatInt(boundaryTimestamp, 10))
				return headers
			},
			requestBody:         "test body",
			verifySignatureFunc: mockVerifySignature(false),
			expectedError:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create middleware
			middleware := newSignatureVerifierMiddleware(tt.verifySignatureFunc, 1024*1024)

			// Create test handler that checks for signature errors
			var capturedError *SignatureError
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sigErr, exists := getSignatureErrorFromContext(r.Context())
				if exists {
					capturedError = sigErr
				}
				w.WriteHeader(http.StatusOK)
			})

			// Wrap handler with middleware
			wrappedHandler := middleware(testHandler)

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte(tt.requestBody)))
			req.Header = tt.setupHeaders()

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute request
			wrappedHandler.ServeHTTP(rr, req)

			// Verify results
			if tt.expectedError == nil {
				require.Nil(t, capturedError)
			} else {
				require.NotNil(t, capturedError)
				require.Equal(t, tt.expectedError.ConnectCode, capturedError.ConnectCode)
				require.Contains(t, capturedError.Message, tt.expectedError.Message)
			}
		})
	}
}
