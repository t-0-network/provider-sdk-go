package provider

import "errors"

var (
	ErrMissingRequiredHeader       = errors.New("missing required header")
	ErrInvalidHeaderEncoding       = errors.New("invalid header encoding")
	ErrUnknownPublicKey            = errors.New("request signed with unknown public key")
	ErrSignatureVerificationFailed = errors.New("signature verification failed")
	ErrInvalidSignature            = errors.New("invalid signature")
	ErrNoSignatureResult           = errors.New("no signature result in context")
	ErrNetworkPublicKeyIsRequired  = errors.New("network public key is not set")
)
