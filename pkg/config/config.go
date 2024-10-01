package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	PostGresAddr     string
	PostGresUser     string
	PostGresPassword string
	PostGresDB       string
}

func NewConfig() *Config {
	viper.AutomaticEnv()
	return &Config{
		PostGresAddr:     viper.GetString("SENDMIND_POSTGRES_ADDR"),
		PostGresUser:     viper.GetString("SENDMIND_POSTGRES_USER"),
		PostGresPassword: viper.GetString("SENDMIND_POSTGRES_PASSWORD"),
		PostGresDB:       viper.GetString("SENDMIND_POSTGRES_DB"),
	}
}
