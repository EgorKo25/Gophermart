package config

import (
	"encoding/hex"
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address         string `env:"ADDRESS"`
	DB              string `env:"DB" default:"postgresql://localhost:5432/postgres"`
	SecretCookieKey []byte `env:"KEY"`
}

func NewConfig() *Config {
	var cfg Config

	_ = env.Parse(cfg)

	var secret string

	flag.StringVar(&cfg.Address,
		"a", "127.0.0.1:8080",
		"Адрес, на котором располагается сервер",
	)
	flag.StringVar(&cfg.DB,
		"d", "postgresql://localhost:5432/postgres",
		"Адрес базы данных с которой работает сервер",
	)
	flag.StringVar(&secret,
		"k", " 7de5a1a5ae85e3aef5376333c3410ca984ef56f0c8082f9d6703414c01affect3",
		"Ключ для шифрования куки",
	)
	flag.Parse()

	cfg.SecretCookieKey, _ = hex.DecodeString(secret)

	return &cfg
}
