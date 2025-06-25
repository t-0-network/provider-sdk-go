package network

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/t-0-network/provider-sdk-go/pkg/constant"
	"github.com/t-0-network/provider-sdk-go/pkg/crypto"
)

func newSigningTransport(signFn crypto.SignFn) *signingTransport {
	return &signingTransport{
		transport: http.DefaultTransport,
		sign:      signFn,
	}
}

// EthereumSigningTransport is an HTTP transport that signs requests with a given signing function.
// It reads the request body, computes its digest, signs it, and adds the signature and public key
// to the request headers before forwarding the request to the underlying transport.
type signingTransport struct {
	transport http.RoundTripper
	sign      crypto.SignFn
}

func (t *signingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Read and restore request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("reading request body: %w", err)
	}
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(body))

	signature, pubKeyBytes, err := t.sign(crypto.LegacyKeccak256(body))
	if err != nil {
		return nil, fmt.Errorf("signing request body: %w", err)
	}

	// Set headers
	req.Header.Set(constant.PublicKeyHeader, "0x"+hex.EncodeToString(pubKeyBytes))
	req.Header.Set(constant.SignatureHeader, "0x"+hex.EncodeToString(signature))

	return t.transport.RoundTrip(req)
}
