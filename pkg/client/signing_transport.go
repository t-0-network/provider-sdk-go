package client

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/t-0-network/provider-sdk-go/pkg/constant"
	"github.com/t-0-network/provider-sdk-go/pkg/internal/crypto"
)

const (
	ethereumSignatureLength int = 65 // 32 bytes r + 32 bytes s + 1 byte recovery ID
)

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

	digest := crypto.LegacyKeccak256(body)

	signature, pubKeyBytes, err := t.sign(digest)
	if err != nil {
		return nil, fmt.Errorf("signing request body: %w", err)
	}

	// Set headers
	req.Header.Set(constant.PublicKeyHeader, "0x"+hex.EncodeToString(pubKeyBytes))
	req.Header.Set(constant.SignatureHeader, "0x"+hex.EncodeToString(signature))

	return t.transport.RoundTrip(req)
}
