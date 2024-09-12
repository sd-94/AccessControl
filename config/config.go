package config

import (
	"log"
	"prac/databases"

	"github.com/spf13/viper"
)

type Auth struct {
	Secret string
	Header string
	Ignore []string
}

type Config struct {
	Auth     Auth
	Port     string
	DBConfig databases.PostgresConfig
}

func LoadConfig() (*Config, error){
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return &config, nil
}
