package crypto

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	btcec_ecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

const (
	ethereumSignatureLength = 65 // 32 bytes r + 32 bytes s + 1 byte recovery ID
)

// SignFn accepts digest hash bytes and returns a signature and
// the public key bytes of the signer or an error.
//
// Signature and public key bytes are following the Ethereum format.
type SignFn func(digestHash []byte) (sig []byte, pubKeyBytes []byte, err error)

func NewSigner(privateKey *ecdsa.PrivateKey) SignFn {
	return func(digestHash []byte) ([]byte, []byte, error) {
		signature, err := sign(digestHash, privateKey)
		if err != nil {
			return nil, nil, fmt.Errorf("signing digest hash: %w", err)
		}

		return signature, GetPublicKeyBytes(&privateKey.PublicKey), nil
	}
}

func NewSignerFromHex(hexedPrivateKey string) (SignFn, error) {
	privateKey, err := HexToECDSA(hexedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("parsing private key: %w", err)
	}

	return NewSigner(privateKey), nil
}

func sign(digest []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	if len(digest) != 32 {
		return nil, errors.New("digest must be 32 bytes")
	}

	// Convert to btcec private key
	btcecPriv, _ := btcec.PrivKeyFromBytes(privateKey.D.Bytes())

	// Use SignCompact which returns recovery ID in the first byte
	compactSig := btcec_ecdsa.SignCompact(btcecPriv, digest, false)

	// compactSig is [recovery_id + 27][r][s] (65 bytes)
	// We need to adjust the recovery ID format for Ethereum
	signature := make([]byte, 65)
	copy(signature[:32], compactSig[1:33])  // R
	copy(signature[32:64], compactSig[33:]) // S
	signature[64] = compactSig[0] - 27      // V (recovery ID, subtract 27)

	return signature, nil
}
