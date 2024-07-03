package util

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	DBDriver      string        `mapstructure:"DB_DRIVER"`
	DBSource      string        `mapstructure:"DB_SOURCE"`
	ServerAddress string        `mapstructure:"SERVER_ADDRESS"`
	EncryptionKey string        `mapstructure:"PAESTO_TOKEN_CREATOR"`
	TokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

func GetConfig(path string) (*Config, error) {
	var once sync.Once
	var err error
	// if found an instance of config return it, else create a new instance only once
	if config != nil {
		fmt.Println("The config found", config)
		return config, nil
	}
	once.Do(func() {
		config_obj, err := LoadConfig(path)
		fmt.Println("The config was not found it was created", config_obj)
		if err != nil {
			return
		}
		config = &config_obj
	})

	return config, err
}
