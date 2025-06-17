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

	"github.com/btcsuite/btcd/btcec/v2"
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

func hexToECDSA(hexedPrivatekey string) (*ecdsa.PrivateKey, error) {
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
