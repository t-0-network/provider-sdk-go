package crypto

import "golang.org/x/crypto/sha3"

func LegacyKeccak256(b []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(b)

	return hash.Sum(make([]byte, 0, 32))
}
