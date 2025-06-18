package client

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
	ErrEmptyPrivateKey = errors.New("private key is not set")
	ErrInvalidTimeOut  = errors.New("timeout must be greater than zero")
)

type client struct {
	baseURL         string
	hexedPrivateKey string
	signFn          crypto.SignFn
	timeOut         time.Duration
}

func (c *client) validate() error {
	if c.baseURL == "" {
		return ErrEmptyBaseURL
	}

	if _, err := url.Parse(c.baseURL); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidBaseURL, err)
	}

	if c.timeOut <= 0 {
		return ErrInvalidTimeOut
	}

	return nil
}

type Option func(*client)

func WithBaseURL(url string) Option {
	return func(c *client) {
		c.baseURL = url
	}
}

func WithHexedPrivateKey(privateKey string) Option {
	return func(c *client) {
		c.hexedPrivateKey = privateKey
	}
}

func WithSignFn(fn crypto.SignFn) Option {
	return func(c *client) {
		c.signFn = fn
	}
}

func WithTimeOut(t time.Duration) Option {
	return func(c *client) {
		c.timeOut = t
	}
}
