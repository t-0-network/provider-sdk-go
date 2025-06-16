// Package network provides a client for interacting with the TZero service.
package network

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1"
	"github.com/t-0-network/provider-sdk-go/service/gen/proto/network/networkconnect"
)

const defaultClientTimeout = time.Duration(15) * time.Second

func NewNetworkServiceClient(cfg Config) (networkconnect.NetworkServiceClient, error) {
	privateKey, err := hexToECDSA(string(cfg.HexedPrivateKey))
	if err != nil {
		return nil, fmt.Errorf("parsing private key: %w", err)
	}

	client := http.Client{
		Timeout: defaultClientTimeout,
		Transport: &signingTransport{
			transport:  http.DefaultTransport,
			privateKey: privateKey,
		},
	}

	return networkconnect.NewNetworkServiceClient(&client, string(cfg.BaseURL)), nil
}

// hexToECDSA converts a 32-byte hex-encoded private key string (optionally prefixed with "0x")
// into a valid ECDSA private key on the secp256k1 curve.
//
// This function is intended as a replacement for Ethereum's crypto.HexToECDSA and produces
// keys fully compatible with Ethereum's cryptography. It uses the secp256k1 curve implementation
// from the github.com/decred/dcrd/dcrec/secp256k1/v4
func hexToECDSA(hexedPrivatekey string) (*ecdsa.PrivateKey, error) {
	hexedPrivatekey = strings.TrimPrefix(hexedPrivatekey, "0x")
	keyBytes, err := hex.DecodeString(hexedPrivatekey)
	if err != nil {
		return nil, fmt.Errorf("decoding private key: %w", err)
	}

	if keyLen := 32; len(keyBytes) != keyLen {
		return nil, fmt.Errorf("invalid private key length: expected %d bytes, got %d", keyLen, len(keyBytes))
	}

	privateKey, _ := secp256k1.PrivKeyFromBytes(keyBytes)

	// Convert to std library ecdsa.PrivateKey
	pubKey := privateKey.PubKey()
	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: secp256k1.S256(),
			X:     pubKey.X,
			Y:     pubKey.Y,
		},
		D: new(big.Int).SetBytes(privateKey.Serialize()),
	}

	return priv, nil
}
