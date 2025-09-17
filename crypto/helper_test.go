package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/t-0-network/provider-sdk-go/crypto"
)

func Test_PrivateKeyHelpers(t *testing.T) {
	// Generated using Ethereum crypto package
	privateKeyHex := "0x691db48202ca70d83cc7f5f3aa219536f9bb2dfe12ebb78a7bb634544858ee92"

	pk, err := crypto.GetPrivateKeyFromHex(privateKeyHex)
	require.NoError(t, err, "failed to get private key from hex")

	hexedPK := crypto.HexPrivateKey(pk)
	require.NoError(t, err, "failed to hex private key")
	require.Equal(t, privateKeyHex, hexedPK)
}

func Test_PublicKeyHelpers(t *testing.T) {
	// Generated using Ethereum crypto package
	publicKeyHex := "0x049bb924680bfba3f64d924bf9040c45dcc215b124b5b9ee73ca8e32c050d042c0bbd8dbb98e3929ed5bc2967f28c3a3b72dd5e24312404598bbf6c6cc47708dc7"

	pk, err := crypto.GetPublicKeyFromHex(publicKeyHex)
	require.NoError(t, err, "failed to get public key from hex")

	hexedPK := crypto.HexPublicKey(pk)
	require.NoError(t, err, "failed to hex public key")
	require.Equal(t, publicKeyHex, hexedPK)

	pkBytes := crypto.GetPublicKeyBytes(pk)
	pkFromBytes, err := crypto.GetPublicKeyFromBytes(pkBytes)
	require.NoError(t, err, "failed to get public key from bytes")

	require.True(t, pk.IsEqual(pkFromBytes))
}
