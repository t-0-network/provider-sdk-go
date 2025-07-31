// Package network provides a client for interacting with the TZero service.
package network

import (
	"fmt"
	"net/http"
	"time"

	"github.com/t-0-network/provider-sdk-go/api/gen/proto/tzero/v1/payment/paymentconnect"
	"github.com/t-0-network/provider-sdk-go/pkg/crypto"
)

type PrivateKeyHexed string

func NewServiceClient(
	privateKey PrivateKeyHexed, opts ...ClientOption,
) (paymentconnect.NetworkServiceClient, error) {
	options := defaultClientOptions
	for _, opt := range opts {
		opt(&options)
	}

	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("validating client options: %w", err)
	}

	if options.signFn == nil {
		if privateKey == "" {
			return nil, ErrEmptyPrivateKey
		}

		defaultSignFn, err := crypto.NewSignerFromHex(string(privateKey))
		if err != nil {
			return nil, fmt.Errorf("creating signer from hexed private key: %w", err)
		}

		options.signFn = defaultSignFn
	}

	client := http.Client{
		Timeout: options.timeout,
		Transport: newSigningTransport(
			options.signFn, time.Now,
		),
	}

	return paymentconnect.NewNetworkServiceClient(&client, options.baseURL), nil
}
