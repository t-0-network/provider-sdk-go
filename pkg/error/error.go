package error

import "errors"

var (
	ErrMissingRequiredHeader       error = errors.New("missing required header")
	ErrInvalidHeaderEncoding       error = errors.New("invalid header encoding")
	ErrUnknownPublicKey            error = errors.New("unknown public key")
	ErrSignatureVerificationFailed error = errors.New("signature verification failed")
	ErrInvalidSignature            error = errors.New("invalid signature")
	ErrNoSignatureResult           error = errors.New("no signature result in context")
)
