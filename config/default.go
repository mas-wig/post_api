package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	PostgreDriver string `mapstructure:"POSTGRES_DRIVER"`
	PostgreURI    string `mapstructure:"POSTGRES_URI"`
	PORT          string `mapstructure:"PORT"`
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
