package app

import (
	"WB_Tech_level_0/internal/store"
	natss "WB_Tech_level_0/internal/transport/nats"
)

type Config struct {
	Nats     *natss.Config `yaml:"nats"`
	Store    *store.Config `yaml:"postgresql"`
	BindAddr string        `yaml:"bind_addr"`
}

// Инициализация конфиг-сервиса
func NewConfig() *Config {
	return &Config{
		BindAddr: "8080",
		Nats:     natss.NewConfig(),
		Store:    store.NewConfig(),
	}
}
