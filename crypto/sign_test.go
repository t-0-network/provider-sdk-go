package crypto_test

import (
	"encoding/hex"
	"testing"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/stretchr/testify/require"
	"github.com/t-0-network/provider-sdk-go/crypto"
)

func TestSignAndVerify(t *testing.T) {
	expectedPubKeyHex := "044fa1465c087aaf42e5ff707050b8f77d2ce92129c5f300686bdd3adfffe44567713bb7931632837c5268a832512e75599b6964f4484c9531c02e96d90384d9f0"
	privateKeyHex := "6b30303de7b26bfb1222b317a52113357f8bb06de00160b4261a2fef9c8b9bd8"

	sign, err := crypto.NewSignerFromHex(privateKeyHex)
	require.NoError(t, err, "failed to create signer from hex private key")

	msg := []byte("please sign me!")
	digest := crypto.LegacyKeccak256(msg)

	signature, pubKeyBytes, err := sign(digest)
	require.NoError(t, err)
	require.Len(t, signature, 65, "signature length should be 65 bytes")
	require.Len(t, pubKeyBytes, 65, "public key bytes length should be 65 bytes")
	require.Equal(t, expectedPubKeyHex, hex.EncodeToString(pubKeyBytes))

	publicKey, err := crypto.GetPublicKeyFromBytes(pubKeyBytes)
	require.NoError(t, err, "failed to get public key from bytes")

	valid := crypto.VerifySignature(publicKey, digest, signature)
	require.True(t, valid, "signature verification should succeed")

	t.Run("Should fail to verify due to key miss match", func(t *testing.T) {
		pk, err := secp256k1.GeneratePrivateKey()
		require.NoError(t, err, "failed to generate private key key")

		sign = crypto.NewSigner(pk)
		signature, pubKeyBytes, err := sign(digest)
		require.NoError(t, err)
		require.Len(t, pubKeyBytes, 65, "public key bytes length should be 65 bytes")

		valid := crypto.VerifySignature(publicKey, digest, signature)
		require.False(t, valid, "signature verification should fail with different key")
	})
}
