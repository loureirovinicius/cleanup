package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

func Start() error {
	viper.AddConfigPath(".")
	viper.AddConfigPath("/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.New("'config.yaml' file does not exist")
		} else {
			return fmt.Errorf("error reading config file: %v", err)
		}
	}
	return nil
}
