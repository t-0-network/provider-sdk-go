package crypto

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

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

		return signature, GetPublicKeyBytes(privateKey), nil
	}
}

func NewSignerFromHex(hexedPrivateKey string) (SignFn, error) {
	privateKey, err := hexToECDSA(hexedPrivateKey)
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

// hexToECDSA converts a 32-byte hex-encoded private key string (optionally prefixed with "0x")
// into a valid ECDSA private key on the secp256k1 curve using btcec.
func hexToECDSA(hexedPrivatekey string) (*ecdsa.PrivateKey, error) {
	hexedPrivatekey = strings.TrimPrefix(hexedPrivatekey, "0x")
	keyBytes, err := hex.DecodeString(hexedPrivatekey)
	if err != nil {
		return nil, fmt.Errorf("decoding private key: %w", err)
	}

	if keyLen := 32; len(keyBytes) != keyLen {
		return nil, fmt.Errorf("invalid private key length: expected %d bytes, got %d", keyLen, len(keyBytes))
	}

	// Convert to btcec private key
	privateKey, publicKey := btcec.PrivKeyFromBytes(keyBytes)

	// Convert to std library ecdsa.PrivateKey
	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: btcec.S256(),
			X:     publicKey.X(),
			Y:     publicKey.Y(),
		},
		D: new(big.Int).SetBytes(privateKey.Serialize()),
	}

	return priv, nil
}

func GetPublicKeyBytes(privateKey *ecdsa.PrivateKey) []byte {
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

func GetPublicKeyFromBytes(pubKeyBytes []byte) (*ecdsa.PublicKey, error) {
	if len(pubKeyBytes) != 65 || pubKeyBytes[0] != 0x04 {
		return nil, errors.New("invalid public key format")
	}

	// Extract x and y coordinates (32 bytes each)
	x := new(big.Int).SetBytes(pubKeyBytes[1:33])
	y := new(big.Int).SetBytes(pubKeyBytes[33:65])

	// Create public key
	pubKey := &ecdsa.PublicKey{
		Curve: btcec.S256(),
		X:     x,
		Y:     y,
	}

	return pubKey, nil
}
