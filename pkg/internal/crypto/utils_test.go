package crypto_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/t-0-network/provider-sdk-go/pkg/internal/crypto"
)

func Test_HexToECDSA(t *testing.T) {
	hexedPubKey := "0x049bb924680bfba3f64d924bf9040c45dcc215b124b5b9ee73ca8e32c050d042c0bbd8dbb98e3929ed5bc2967f28c3a3b72dd5e24312404598bbf6c6cc47708dc7"
	hexedPrivateKey := "691db48202ca70d83cc7f5f3aa219536f9bb2dfe12ebb78a7bb634544858ee92"
	privateKey, err := crypto.HexToECDSA(hexedPrivateKey)
	require.NoError(t, err, "failed to convert hex to ECDSA")

	pubKeyBytes := crypto.GetPublicKeyBytes(&privateKey.PublicKey)
	require.Equal(t, hexedPubKey, fmt.Sprintf("0x%x", pubKeyBytes))
}

func Test_HexToECDSAPublicKey(t *testing.T) {
	hexedPubKey := "0x049bb924680bfba3f64d924bf9040c45dcc215b124b5b9ee73ca8e32c050d042c0bbd8dbb98e3929ed5bc2967f28c3a3b72dd5e24312404598bbf6c6cc47708dc7"
	pubKey, err := crypto.HexToECDSAPublicKey(hexedPubKey)
	require.NoError(t, err, "failed to convert hex to ECDSA public key")

	pubKeyBytes := crypto.GetPublicKeyBytes(pubKey)
	require.Len(t, pubKeyBytes, 65, "public key bytes length should be 65 bytes")
	require.Equal(t, hexedPubKey, fmt.Sprintf("0x%x", pubKeyBytes))

	pubKeyFromBytes, err := crypto.GetPublicKeyFromBytes(pubKeyBytes)
	require.NoError(t, err, "failed to get public key from bytes")
	require.True(t, pubKey.Equal(pubKeyFromBytes))
}
