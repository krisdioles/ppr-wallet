package config

import (
	"log"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	Bank1  Bank1Config
}

type ServerConfig struct {
	Port int64
}

type Bank1Config struct {
	Hostname             string
	APIKey               string
	DisbursementEndpoint string
}

var (
	config *Config
	once   sync.Once
)

func All() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("${PWD}/.")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AllowEmptyEnv(true)

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}

		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("Unable to decode into struct, %v", err)
		}
	})

	return config
}

func Get() *Config {
	return config
}
