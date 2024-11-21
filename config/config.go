package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Start all the initial configuration required by this program
func Start(provider string) error {
	// Enable config.yaml file to be used
	viper.AddConfigPath(".")
	viper.AddConfigPath("/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// Enable environment variables
	viper.AutomaticEnv()

	// Verify the cloud provider being used for config initialization
	switch provider {
	case "aws":
		err := loadAWSConfig()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("provider %s is not supported", provider)
	}

	// Attempt to read the config.yaml file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.New("'config.yaml' file does not exist")
		} else {
			return fmt.Errorf("error reading config file: %v", err)
		}
	}
	return nil
}

func loadAWSConfig() error {
	viper.SetEnvPrefix("AWS")

	// AWS_REGION env variable
	if err := viper.BindEnv("region"); err != nil {
		return fmt.Errorf("error binding region variable: %w", err)
	}

	// AWS_PROFILE_NAME env variable
	if err := viper.BindEnv("profile.name", "AWS_PROFILE_NAME"); err != nil {
		return fmt.Errorf("error binding profile_name variable: %w", err)
	}

	// AWS_ACCESS_KEY env variable
	if err := viper.BindEnv("credentials.access_key", "AWS_ACCESS_KEY"); err != nil {
		return fmt.Errorf("error binding access_key variable: %w", err)
	}

	// AWS_SECRET_KEY env variable
	if err := viper.BindEnv("credentials.secret_key", "AWS_SECRET_KEY"); err != nil {
		return fmt.Errorf("error binding secret_key variable: %w", err)
	}

	return nil
}
