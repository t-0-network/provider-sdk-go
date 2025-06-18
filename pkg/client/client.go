// Package network provides a client for interacting with the TZero service.
package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/t-0-network/provider-sdk-go/pkg/gen/proto/network/networkconnect"
	"github.com/t-0-network/provider-sdk-go/pkg/internal/crypto"
)

func NewNetworkServiceClient(opts ...Option) (networkconnect.NetworkServiceClient, error) {
	c := &client{
		baseURL: "https://api.t-0.network",
		timeOut: time.Duration(15) * time.Second,
	}
	for _, opt := range opts {
		opt(c)
	}

	if c.signFn == nil {
		defaultSignFn, err := crypto.NewSignerFromHex(c.hexedPrivateKey)
		if err != nil {
			return nil, fmt.Errorf("creating signer from hexed private key: %w", err)
		}
		c.signFn = defaultSignFn
	}

	if err := c.validate(); err != nil {
		return nil, fmt.Errorf("validating client options: %w", err)
	}

	client := http.Client{
		Timeout: c.timeOut,
		Transport: &signingTransport{
			transport: http.DefaultTransport,
			sign:      c.signFn,
		},
	}

	return networkconnect.NewNetworkServiceClient(&client, c.baseURL), nil
}
