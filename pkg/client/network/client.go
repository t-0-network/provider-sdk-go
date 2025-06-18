// Package network provides a client for interacting with the TZero service.
package network

import (
	"fmt"
	"net/http"

	"github.com/t-0-network/provider-sdk-go/pkg/client"
	"github.com/t-0-network/provider-sdk-go/pkg/gen/proto/network/networkconnect"
	"github.com/t-0-network/provider-sdk-go/pkg/internal/crypto"
)

func NewServiceClient(opts ...ClientOption) (networkconnect.NetworkServiceClient, error) {
	options := defaultClientOptions
	for _, opt := range opts {
		opt(options)
	}

	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("validating client options: %w", err)
	}

	if options.signFn == nil {
		if options.hexedProviderPrivateKey == "" {
			return nil, ErrEmptyPrivateKey
		}

		defaultSignFn, err := crypto.NewSignerFromHex(options.hexedProviderPrivateKey)
		if err != nil {
			return nil, fmt.Errorf("creating signer from hexed private key: %w", err)
		}

		options.signFn = defaultSignFn
	}

	client := http.Client{
		Timeout:   options.timeout,
		Transport: client.NewEthereumSigningTransport(options.signFn),
	}

	return networkconnect.NewNetworkServiceClient(&client, options.baseURL), nil
}
