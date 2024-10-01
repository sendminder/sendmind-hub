package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	PostGresUser     string
	PostGresPassword string
	PostGresDB       string
}

func NewConfig() *Config {
	viper.AutomaticEnv()
	return &Config{
		PostGresUser:     viper.GetString("SENDMIND_POSTGRES_USER"),
		PostGresPassword: viper.GetString("SENDMIND_POSTGRES_PASSWORD"),
		PostGresDB:       viper.GetString("SENDMIND_POSTGRES_DB"),
	}
}
