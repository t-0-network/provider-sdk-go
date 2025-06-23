package provider

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/t-0-network/provider-sdk-go/pkg/internal/crypto"
)

var (
	ErrEmptyBaseURL    = errors.New("base URL is not set")
	ErrInvalidBaseURL  = errors.New("base URL is not valid")
	ErrEmptyPrivateKey = errors.New("network private key is not set")
	ErrInvalidTimeOut  = errors.New("timeout must be greater than zero")
)

type clientOptions struct {
	baseURL                string
	networkPrivateKeyHexed string
	signFn                 crypto.SignFn
	timeout                time.Duration
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
	baseURL:                "",
	networkPrivateKeyHexed: "",
	signFn:                 nil,
	timeout:                defaultTimeout,
}

type Option func(*clientOptions)

func WithBaseURL(url string) Option {
	return func(c *clientOptions) {
		c.baseURL = url
	}
}

func WithNetworkPrivateKeyHexed(privateKey string) Option {
	return func(c *clientOptions) {
		c.networkPrivateKeyHexed = privateKey
	}
}

func WithSignatureFunction(fn crypto.SignFn) Option {
	return func(c *clientOptions) {
		c.signFn = fn
	}
}

func WithTimeout(t time.Duration) Option {
	return func(c *clientOptions) {
		c.timeout = t
	}
}
