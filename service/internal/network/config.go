package network

type Config struct {
	BaseURL         string `envDefault:"https://api.t-0.network/"`
	HexedPrivateKey string `env:"HEXED_PRIVATE_KEY"`
}
