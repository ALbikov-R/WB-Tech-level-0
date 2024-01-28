package natss

type Config struct {
	URL     string `yaml:"natsURL"`
	Subject string `yaml:"subject"`
}

// Иницализация конфига NATS
func NewConfig() *Config {
	return &Config{}
}
