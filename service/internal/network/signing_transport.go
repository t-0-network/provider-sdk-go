package network

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	btcec_ecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"

	"golang.org/x/crypto/sha3"
)

const (
	signatureHeader = "X-Signature"
	publicKeyHeader = "X-Public-Key"

	ethereumSignatureLength int = 65 // 32 bytes r + 32 bytes s + 1 byte recovery ID
)

type signingTransport struct {
	transport  http.RoundTripper
	privateKey *ecdsa.PrivateKey
}

func NewSigningTransport(transport http.RoundTripper, ecdsaKey *ecdsa.PrivateKey) *signingTransport {
	return &signingTransport{
		transport:  transport,
		privateKey: ecdsaKey,
	}
}

func (t *signingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Read and restore request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("reading request body: %w", err)
	}
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(body))

	hash := sha3.NewLegacyKeccak256()
	hash.Write(body)
	hashBytes := hash.Sum(nil)

	// Sign the hash
	// signature, err := t.sign(hashBytes)
	signature, err := sign(hashBytes, t.privateKey)
	if err != nil {
		return nil, fmt.Errorf("signing request body: %w", err)
	}

	// Get public key in Ethereum format
	pubKeyBytes := getPublicKey(t.privateKey)

	// Set headers
	req.Header.Set(publicKeyHeader, "0x"+hex.EncodeToString(pubKeyBytes))
	req.Header.Set(signatureHeader, "0x"+hex.EncodeToString(signature))

	return t.transport.RoundTrip(req)
}

// getPublicKey returns the public key in Ethereum format.
func getPublicKey(privateKey *ecdsa.PrivateKey) []byte {
	// Ethereum uses uncompressed public key format: 0x04 + x + y
	pubKeyBytes := make([]byte, 65)
	pubKeyBytes[0] = 0x04

	// Pad coordinates to 32 bytes each
	xBytes := privateKey.PublicKey.X.Bytes()
	yBytes := privateKey.PublicKey.Y.Bytes()

	copy(pubKeyBytes[33-len(xBytes):33], xBytes)
	copy(pubKeyBytes[65-len(yBytes):65], yBytes)

	return pubKeyBytes
}

func sign(digest []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	if len(digest) != 32 {
		return nil, errors.New("digest must be 32 bytes")
	}

	// Convert to *btcec.PrivateKey
	btcecPrivateKey, _ := btcec.PrivKeyFromBytes(privateKey.D.Bytes())

	// Use SignCompact which handles everything for us
	compactSig := btcec_ecdsa.SignCompact(btcecPrivateKey, digest, false)

	// SignCompact returns: [recovery_id + 27][r 32 bytes][s 32 bytes] (65 bytes total)
	// go-ethereum expects: [r 32 bytes][s 32 bytes][recovery_id] (65 bytes total)

	// Create result in go-ethereum format
	result := make([]byte, 65)

	// Copy r (bytes 1-32 from compact sig to bytes 0-31 in result)
	copy(result[0:32], compactSig[1:33])

	// Copy s (bytes 33-64 from compact sig to bytes 32-63 in result)
	copy(result[32:64], compactSig[33:65])

	// Set recovery ID (remove Bitcoin's +27 offset)
	result[64] = compactSig[0] - 27

	if !verifySignature(digest, result, privateKey) {
		return nil, errors.New("local signature verification failed")
	}

	return result, nil
}

func verifySignature(digest []byte, signature []byte, privateKey *ecdsa.PrivateKey) bool {
	if len(signature) != 65 {
		return false
	}

	// Convert back to compact format for verification
	compactSig := make([]byte, 65)
	compactSig[0] = signature[64] + 27        // Add Bitcoin offset back
	copy(compactSig[1:33], signature[0:32])   // r
	copy(compactSig[33:65], signature[32:64]) // s

	// Recover public key using btcec
	recoveredPubKey, _, err := btcec_ecdsa.RecoverCompact(compactSig, digest)
	if err != nil {
		return false
	}

	// Compare with expected public key
	btcecPrivateKey, _ := btcec.PrivKeyFromBytes(privateKey.D.Bytes())
	expectedPubKey := btcecPrivateKey.PubKey()

	return recoveredPubKey.IsEqual(expectedPubKey)
}
