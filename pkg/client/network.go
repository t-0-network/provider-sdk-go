// Package network provides a client for interacting with the TZero service.
package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/t-0-network/provider-sdk-go/pkg/gen/proto/network/networkconnect"
	"github.com/t-0-network/provider-sdk-go/pkg/internal/crypto"
)

const defaultClientTimeout = time.Duration(15) * time.Second

func NewNetworkServiceClient(cfg Config) (networkconnect.NetworkServiceClient, error) {
	signFn, err := crypto.NewSignerFromHex(cfg.HexedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("creating signer from hexed private key: %w", err)
	}

	client := http.Client{
		Timeout: defaultClientTimeout,
		Transport: &signingTransport{
			transport: http.DefaultTransport,
			sign:      signFn,
		},
	}

	return networkconnect.NewNetworkServiceClient(&client, cfg.BaseURL), nil
}
