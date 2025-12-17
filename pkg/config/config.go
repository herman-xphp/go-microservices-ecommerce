package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port  string `mapstructure:"PORT"`
	DBUrl string `mapstructure:"DB_URL"`
	Env   string `mapstructure:"ENV"`
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Config file not found, relying on env vars: %v", err)
	}

	err := viper.Unmarshal(&config)
	return config, err
}
