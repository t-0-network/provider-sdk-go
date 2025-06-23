package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

const bip39EntropySize = 128

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Generate mnemonic
	entropy, err := bip39.NewEntropy(bip39EntropySize)
	if err != nil {
		log.Fatal(err)
	}

	defaultMnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Fatal(err)
	}

	mnemonic := flag.String("mnemonic", defaultMnemonic, "type your mnemonic here, or leave empty to generate a new one")
	flag.Parse()

	if !bip39.IsMnemonicValid(*mnemonic) {
		log.Fatalf("Invalid mnemonic: %s", *mnemonic)
	}

	// Generate seed and master key
	seed := bip39.NewSeed(*mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		log.Fatal(err)
	}

	// Derive path m/44'/60'/0'/0/0 (first address)
	key := masterKey
	derivationPath := []uint32{
		bip32.FirstHardenedChild + 44, // purpose
		bip32.FirstHardenedChild + 60, // coin type (ETH)
		bip32.FirstHardenedChild + 0,  // account
		0,                             // change
		0,                             // address index
	}

	for _, index := range derivationPath {
		key, err = key.NewChildKey(index)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Generate the key pair
	privateKey := secp256k1.PrivKeyFromBytes(key.Key)
	privateKeyHex := hex.EncodeToString(privateKey.Serialize())
	pubKeyHex := hex.EncodeToString(privateKey.PubKey().SerializeUncompressed())

	fmt.Println("Mnemonic:", *mnemonic)
	fmt.Printf("Private Key: 0x%s\n", privateKeyHex)
	fmt.Printf("Public Key: 0x%s\n", pubKeyHex)
}
