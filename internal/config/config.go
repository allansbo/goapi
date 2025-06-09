package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// EnvConfig is the configuration for the application.
type EnvConfig struct {
	DBUser       string `mapstructure:"DB_USER"`
	DBPass       string `mapstructure:"DB_PASS"`
	DBName       string `mapstructure:"DB_NAME"`
	DBCollection string `mapstructure:"DB_COLLECTION"`
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       string `mapstructure:"DB_PORT"`
	AppPort      string `mapstructure:"APP_PORT"`
}

// isValidConfig is a function that checks if the configuration is valid.
func (e *EnvConfig) isValidConfig() error {
	requiredFields := map[string]string{
		"DB_USER":       e.DBUser,
		"DB_PASS":       e.DBPass,
		"DB_NAME":       e.DBName,
		"DB_COLLECTION": e.DBCollection,
		"DB_HOST":       e.DBHost,
		"DB_PORT":       e.DBPort,
		"APP_PORT":      e.AppPort,
	}

	for key, value := range requiredFields {
		if value == "" {
			return fmt.Errorf("%s is required", key)
		}
	}

	return nil
}

// LoadEnvConfig is a function that loads the configuration from the environment variables.
func LoadEnvConfig() (*EnvConfig, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config EnvConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := config.isValidConfig(); err != nil {
		return nil, err
	}

	return &config, nil
}
