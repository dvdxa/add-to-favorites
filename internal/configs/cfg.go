package configs

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Dbname   string
	Password string
	SSLMode  string
}

func InitConfig() (config Config, err error) {
	viper.AddConfigPath("./internal/configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
	}
	return config, nil
}
