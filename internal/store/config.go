package store

type Config struct {
	Database_URL string `yaml:"database_URL"`
}

// Инициалзация конфига для БД
func NewConfig() *Config {
	return &Config{}
}
