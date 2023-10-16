package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	PostgreDriver          string `mapstructure:"POSTGRES_DRIVER"`
	PostgreURI             string `mapstructure:"POSTGRES_URI"`
	PORT                   string `mapstructure:"PORT"`
	Origin                 string `mapstructure:"ORIGIN"`
	AccessTokenPublicKey   string `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	AccessTokenPrivateKey  string `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	RefreshTokenPrivateKey string `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenMaxAge      int    `mapstructure:"ACCESS_TOKEN_MAX_AGE"`
	RefreshTokenMaxAge     int    `mapstructure:"REFRESH_TOKEN_MAX_AGE"`
}

func LoadConfig() (c *Config, err error) {
	viper.AddConfigPath("../simple_api")
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		log.Fatal("env file tidak ditemukan pada root folder")
	}
	err = viper.Unmarshal(&c)
	return
}
