// Package network provides a client for interacting with the TZero service.
package network

import (
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/t-0-network/provider-sdk-go/crypto"
)

type PrivateKeyHexed string

type ClientFactory[T any] func(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) T

func NewServiceClient[T any](
	privateKey PrivateKeyHexed, clientFactory ClientFactory[T], opts ...ClientOption,
) (T, error) {
	options := defaultClientOptions
	for _, opt := range opts {
		opt(&options)
	}

	var t T

	if err := options.validate(); err != nil {
		return t, fmt.Errorf("validating client options: %w", err)
	}

	if options.signFn == nil {
		if privateKey == "" {
			return t, ErrEmptyPrivateKey
		}

		defaultSignFn, err := crypto.NewSignerFromHex(string(privateKey))
		if err != nil {
			return t, fmt.Errorf("creating signer from hexed private key: %w", err)
		}

		options.signFn = defaultSignFn
	}

	client := http.Client{
		Timeout: options.timeout,
		Transport: NewSigningTransport(
			options.signFn, time.Now,
		),
	}

	return clientFactory(&client, options.baseURL, options.connectOptions...), nil
}
