package crypto

import (
	"crypto/ecdsa"
	"math/big"
)

type VerifySignatureFn func(digest []byte, signature []byte, privateKey *ecdsa.PrivateKey) bool

// VerifySignature verifies an Ethereum signature by performing ECDSA verification.
func VerifySignature(pubKey *ecdsa.PublicKey, digest []byte, signature []byte) bool {
	if len(digest) != 32 {
		return false
	}

	// Support both 64-byte (r+s) and 65-byte (r+s+v) signatures
	if len(signature) != 64 && len(signature) != 65 {
		return false
	}

	// Extract R and S from signature (always first 64 bytes)
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])

	// Validate r and s are not zero
	if r.Sign() == 0 || s.Sign() == 0 {
		return false
	}

	// Perform ECDSA verification
	return ecdsa.Verify(pubKey, digest, r, s)
}
