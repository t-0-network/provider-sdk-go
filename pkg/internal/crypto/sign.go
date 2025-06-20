package crypto

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

const (
	ethereumSignatureLength = 65 // 32 bytes r + 32 bytes s + 1 byte recovery ID
)

// SignFn accepts a raw message/payload hashes and signs it
type SignFn func(digest []byte) (sig []byte, pubKeyBytes []byte, err error)

func NewSigner(privateKey *secp256k1.PrivateKey) SignFn {
	return func(message []byte) ([]byte, []byte, error) {
		return sign(message, privateKey),
			GetPublicKeyBytes(privateKey.PubKey()), nil
	}
}

func NewSignerFromHex(hexedPrivateKey string) (SignFn, error) {
	privateKey, err := GetPrivateKeyFromHex(hexedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("creating signer from hexed private key: %w", err)
	}

	return NewSigner(privateKey), nil
}

func sign(digest []byte, privateKey *secp256k1.PrivateKey) []byte {
	// Use SignCompact which returns recovery ID in the first byte
	compactSig := ecdsa.SignCompact(privateKey, digest, false)

	// compactSig is [recovery_id + 27][r][s] (65 bytes)
	// We need to adjust the recovery ID format for Ethereum
	signature := make([]byte, ethereumSignatureLength)
	copy(signature[:32], compactSig[1:33])  // R
	copy(signature[32:64], compactSig[33:]) // S
	signature[64] = compactSig[0] - 27      // V (recovery ID, subtract 27)

	return signature
}
