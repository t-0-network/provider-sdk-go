package crypto

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

func GetPrivateKeyBytes(privateKey *secp256k1.PrivateKey) []byte {
	return privateKey.Serialize()
}

func GetPrivateKeyFromHex(privateKeyHexed string) (*secp256k1.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(strings.TrimPrefix(strings.ToLower(privateKeyHexed), "0x"))
	if err != nil {
		return nil, fmt.Errorf("decoding private key hex: %w", err)
	}

	privateKey := secp256k1.PrivKeyFromBytes(privateKeyBytes)
	if privateKey == nil {
		return nil, errors.New("invalid private key bytes")
	}

	return privateKey, nil
}

func HexPrivateKey(privateKey *secp256k1.PrivateKey) string {
	return "0x" + hex.EncodeToString(GetPrivateKeyBytes(privateKey))
}

func GetPublicKeyBytes(publicKey *secp256k1.PublicKey) []byte {
	return publicKey.SerializeUncompressed()
}

func GetPublicKeyFromBytes(pubKeyBytes []byte) (*secp256k1.PublicKey, error) {
	publicKey, err := secp256k1.ParsePubKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("parsing public key bytes: %w", err)
	}

	return publicKey, nil
}

func GetPublicKeyFromHex(publicKeyHexed string) (*secp256k1.PublicKey, error) {
	pubKeyBytes, err := hex.DecodeString(strings.TrimPrefix(strings.ToLower(publicKeyHexed), "0x"))
	if err != nil {
		return nil, fmt.Errorf("decoding public key hex: %w", err)
	}

	return GetPublicKeyFromBytes(pubKeyBytes)
}

func HexPublicKey(publicKey *secp256k1.PublicKey) string {
	return "0x" + hex.EncodeToString(GetPublicKeyBytes(publicKey))
}
