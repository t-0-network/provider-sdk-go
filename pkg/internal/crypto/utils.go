package crypto

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
)

// HexToECDSA converts a 32-byte hex-encoded private key string (optionally prefixed with "0x")
// into a valid ECDSA private key on the secp256k1 curve using btcec.
func HexToECDSA(hexedPrivatekey string) (*ecdsa.PrivateKey, error) {
	keyBytes, err := hex.DecodeString(strings.TrimPrefix(hexedPrivatekey, "0x"))
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

func HexToECDSAPublicKey(hexedPublicKey string) (*ecdsa.PublicKey, error) {
	keyBytes, err := hex.DecodeString(strings.TrimPrefix(hexedPublicKey, "0x"))
	if err != nil {
		return nil, fmt.Errorf("decoding public key: %w", err)
	}

	// Parse using btcec which handles both compressed and uncompressed formats
	btcecPubKey, err := btcec.ParsePubKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("parsing public key: %w", err)
	}

	// Convert to std library ecdsa.PublicKey
	pubKey := &ecdsa.PublicKey{
		Curve: btcec.S256(),
		X:     btcecPubKey.X(),
		Y:     btcecPubKey.Y(),
	}

	return pubKey, nil
}

func GetPublicKeyBytes(publicKey *ecdsa.PublicKey) []byte {
	// Ethereum uses uncompressed public key format: 0x04 + x + y
	pubKeyBytes := make([]byte, 65)
	pubKeyBytes[0] = 0x04

	// Pad coordinates to 32 bytes each
	xBytes := publicKey.X.Bytes()
	yBytes := publicKey.Y.Bytes()

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
