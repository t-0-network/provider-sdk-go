// Package network provides a client for interacting with the TZero service.
package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/t-0-network/provider-sdk-go/pkg/gen/proto/network/networkconnect"
	"github.com/t-0-network/provider-sdk-go/pkg/internal/crypto"
)

const (
	defaultBaseURL = "https://api.t-0.network"
	defaultTimeout = 15 * time.Second
)

func NewNetworkServiceClient(opts ...Option) (networkconnect.NetworkServiceClient, error) {
	c := &client{
		baseURL: defaultBaseURL,
		timeOut: defaultTimeout,
	}
	for _, opt := range opts {
		opt(c)
	}

	if err := c.validate(); err != nil {
		return nil, fmt.Errorf("validating client options: %w", err)
	}

	if c.signFn == nil {
		if c.hexedPrivateKey == "" {
			return nil, ErrEmptyPrivateKey
		}

		defaultSignFn, err := crypto.NewSignerFromHex(c.hexedPrivateKey)
		if err != nil {
			return nil, fmt.Errorf("creating signer from hexed private key: %w", err)
		}

		c.signFn = defaultSignFn
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
