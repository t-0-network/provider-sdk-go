// Package network provides a client for interacting with the TZero service.
package network

import (
	"fmt"
	"net/http"

	"github.com/t-0-network/provider-sdk-go/pkg/gen/proto/network/networkconnect"
	"github.com/t-0-network/provider-sdk-go/pkg/internal/crypto"
	"github.com/t-0-network/provider-sdk-go/pkg/internal/transport"
)

func NewServiceClient(opts ...ClientOption) (networkconnect.NetworkServiceClient, error) {
	options := defaultClientOptions
	for _, opt := range opts {
		opt(&options)
	}

	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("validating client options: %w", err)
	}

	if options.signFn == nil {
		if options.providerPrivateKeyHexed == "" {
			return nil, ErrEmptyPrivateKey
		}

		defaultSignFn, err := crypto.NewSignerFromHex(options.providerPrivateKeyHexed)
		if err != nil {
			return nil, fmt.Errorf("creating signer from hexed private key: %w", err)
		}

		options.signFn = defaultSignFn
	}

	client := http.Client{
		Timeout:   options.timeout,
		Transport: transport.NewEthereumSigningTransport(options.signFn),
	}

	return networkconnect.NewNetworkServiceClient(&client, options.baseURL), nil
}
