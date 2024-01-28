package natss

import (
	"github.com/nats-io/nats.go"
)

type NatsServer struct {
	Config *Config
	Natc   *nats.Conn
}

// Иницализация  NATS
func New(config *Config) *NatsServer {
	return &NatsServer{
		Config: config,
	}
}

// Созданние соединения с NATS сервером
func (n *NatsServer) InitConnect() error {
	nc, err := nats.Connect(n.Config.URL)

	if err != nil {
		return err
	}
	n.Natc = nc
	return nil
}

// Закрытие соединения с NATS сервером
func (n *NatsServer) Close() {
	n.Natc.Close()
}
