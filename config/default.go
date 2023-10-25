package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	PostgreDriver          string        `mapstructure:"POSTGRES_DRIVER"`
	PostgreURI             string        `mapstructure:"POSTGRES_URI"`
	PORT                   string        `mapstructure:"PORT"`
	Origin                 string        `mapstructure:"ORIGIN"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAX_AGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAX_AGE"`
	AccessTokenExpiredIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiredIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
}

func LoadConfig() (c *Config, err error) {
	rootdir, _ := os.Getwd()
	viper.AddConfigPath(rootdir)
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		log.Fatal("env file tidak ditemukan pada root folder")
	}
	err = viper.Unmarshal(&c)
	return
}
