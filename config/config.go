package config

import (
	"log"

	"github.com/spf13/viper"
)

func Start() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("'config.yaml' file does not exist.")
		} else {
			log.Fatalf("error reading config file: %v", err)
		}
	}
}
