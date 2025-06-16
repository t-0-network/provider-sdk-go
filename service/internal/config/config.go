package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/t-0-network/provider-sdk-go/service/internal/network"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`

	NetworkClient network.Config `envPrefix:"NETWORK_CLIENT_"`
}

func NewConfig() (Config, error) {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Trace().Err(err).Msg("local.env is not loaded")
	}

	_ = godotenv.Load("/mnt/secrets/local.env")

	cfg := Config{}
	return cfg, env.Parse(&cfg)
}
