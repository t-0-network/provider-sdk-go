// Package network provides a client for interacting with the TZero service.
package provider

import (
	"fmt"
	"net/http"
	"time"

	"github.com/t-0-network/provider-sdk-go/pkg/gen/proto/network/networkconnect"
	"github.com/t-0-network/provider-sdk-go/pkg/internal/crypto"
	"github.com/t-0-network/provider-sdk-go/pkg/internal/transport"
)

const (
	defaultTimeout = 15 * time.Second
)

func NewServiceClient(opts ...Option) (networkconnect.ProviderServiceClient, error) {
	options := defaultClientOptions
	for _, opt := range opts {
		opt(&options)
	}

	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("validating client options: %w", err)
	}

	if options.signFn == nil {
		if options.networkPrivateKeyHexed == "" {
			return nil, ErrEmptyPrivateKey
		}

		defaultSignFn, err := crypto.NewSignerFromHex(options.networkPrivateKeyHexed)
		if err != nil {
			return nil, fmt.Errorf("creating signer from hexed private key: %w", err)
		}

		options.signFn = defaultSignFn
	}

	client := http.Client{
		Timeout:   options.timeout,
		Transport: transport.NewEthereumSigningTransport(options.signFn),
	}

	return networkconnect.NewProviderServiceClient(&client, options.baseURL), nil
}
