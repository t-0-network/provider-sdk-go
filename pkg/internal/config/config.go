package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/t-0-network/provider-sdk-go/pkg/client"
)

type Config struct {
	Server ServerConfig `envPrefix:"SERVER_"`

	NetworkClient client.Config `envPrefix:"NETWORK_CLIENT_"`
}

type ServerConfig struct {
	Address      string        `env:"PORT" envDefault:":8080"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
}

func LoadConfig() (Config, error) {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Trace().Err(err).Msg("local.env is not loaded")
	}

	_ = godotenv.Load("/mnt/secrets/local.env")

	cfg := Config{}
	return cfg, env.Parse(&cfg)
}
