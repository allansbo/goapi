package config

import (
	"github.com/spf13/viper"
)

type EnvConfig struct {
	DBUser  string `mapstructure:"DB_USER"`
	DBPass  string `mapstructure:"DB_PASSWORD"`
	DBName  string `mapstructure:"DB_NAME"`
	DBHost  string `mapstructure:"DB_HOST"`
	DBPort  string `mapstructure:"DB_PORT"`
	AppPort string `mapstructure:"APP_PORT"`
}

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

	return &config, nil
}
