package network

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/t-0-network/provider-sdk-go/pkg/constant"
	"github.com/t-0-network/provider-sdk-go/pkg/crypto"
)

func newSigningTransport(signFn crypto.SignFn, timeNow func() time.Time) *signingTransport {
	return &signingTransport{
		transport: http.DefaultTransport,
		sign:      signFn,
		timeNow:   timeNow,
	}
}

// EthereumSigningTransport is an HTTP transport that signs requests with a given signing function.
// It reads the request body, computes its digest, signs it, and adds the signature and public key
// to the request headers before forwarding the request to the underlying transport.
type signingTransport struct {
	transport http.RoundTripper
	sign      crypto.SignFn
	timeNow   func() time.Time
}

func (t *signingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Read and restore request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("reading request body: %w", err)
	}
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(body))

	// Get current timestamp in milliseconds
	timestamp := t.timeNow().UnixMilli()

	// Convert timestamp to little-endian (8 bytes for int64)
	timestampBytes := [8]byte{}
	binary.LittleEndian.PutUint64(timestampBytes[:], uint64(timestamp))

	// Prepend timestamp bytes to the body and compute the digest
	digest := crypto.LegacyKeccak256(append(body, timestampBytes[:]...))

	signature, pubKeyBytes, err := t.sign(digest)
	if err != nil {
		return nil, fmt.Errorf("signing request body: %w", err)
	}

	// Set headers
	req.Header.Set(constant.PublicKeyHeader, "0x"+hex.EncodeToString(pubKeyBytes))
	req.Header.Set(constant.SignatureHeader, "0x"+hex.EncodeToString(signature))
	req.Header.Set(constant.SignatureTimestampHeader, strconv.FormatInt(timestamp, 10))

	return t.transport.RoundTrip(req)
}
