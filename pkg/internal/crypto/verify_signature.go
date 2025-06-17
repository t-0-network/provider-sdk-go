package crypto

import (
	"crypto/ecdsa"
	"math/big"
)

type VerifySignatureFn func(digest []byte, signature []byte, privateKey *ecdsa.PrivateKey) bool

// VerifySignature verifies an Ethereum signature by performing ECDSA verification.
func VerifySignature(pubKey *ecdsa.PublicKey, digest []byte, signature []byte) bool {
	if len(digest) != 32 || len(signature) != 65 {
		return false
	}

	// Extract R and S from signature (ignore recovery ID for verification)
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])

	// Perform ECDSA verification
	return ecdsa.Verify(pubKey, digest, r, s)
}
