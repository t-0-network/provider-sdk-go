package crypto

import (
	"crypto/ecdsa"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	dcrececdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
)

type VerifySignatureFn func(digest []byte, signature []byte, privateKey *ecdsa.PrivateKey) bool

func VerifySignature(pubKey *secp256k1.PublicKey, digest []byte, signature []byte) bool {
	if len(digest) != 32 {
		return false
	}

	// Support both 64-byte (r+s) and 65-byte (r+s+v) signatures
	if len(signature) != 64 && len(signature) != 65 {
		return false
	}

	// Extract R and S from the signature (ignore recovery ID at position 64)
	rBytes := signature[:32]
	sBytes := signature[32:64]

	// Convert to ModNScalar
	var r, s secp256k1.ModNScalar
	r.SetByteSlice(rBytes)
	s.SetByteSlice(sBytes)

	// Create the signature
	sig := dcrececdsa.NewSignature(&r, &s)

	// Verify
	return sig.Verify(digest, pubKey)
}
