package network

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/t-0-network/provider-sdk-go/pkg/internal/crypto"
)

const (
	defaultBaseURL = "https://api.t-0.network"
	defaultTimeout = 15 * time.Second
)

var (
	ErrEmptyBaseURL    = errors.New("base URL is not set")
	ErrInvalidBaseURL  = errors.New("base URL is not valid")
	ErrEmptyPrivateKey = errors.New("provider private key is not set")
	ErrInvalidTimeOut  = errors.New("timeout must be greater than zero")
)

type clientOptions struct {
	baseURL                 string
	providerPrivateKeyHexed string
	signFn                  crypto.SignFn
	timeout                 time.Duration
}

func (c *clientOptions) validate() error {
	if c.baseURL == "" {
		return ErrEmptyBaseURL
	}

	if _, err := url.Parse(c.baseURL); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidBaseURL, err)
	}

	if c.timeout <= 0 {
		return ErrInvalidTimeOut
	}

	return nil
}

var defaultClientOptions = clientOptions{
	baseURL:                 defaultBaseURL,
	providerPrivateKeyHexed: "",
	signFn:                  nil,
	timeout:                 defaultTimeout,
}

type ClientOption func(*clientOptions)

func WithBaseURL(url string) ClientOption {
	return func(c *clientOptions) {
		c.baseURL = url
	}
}

func WithProviderPrivateKeyHexed(privateKey string) ClientOption {
	return func(c *clientOptions) {
		c.providerPrivateKeyHexed = privateKey
	}
}

func WithSignatureFunction(fn crypto.SignFn) ClientOption {
	return func(c *clientOptions) {
		c.signFn = fn
	}
}

func WithTimeout(t time.Duration) ClientOption {
	return func(c *clientOptions) {
		c.timeout = t
	}
}
