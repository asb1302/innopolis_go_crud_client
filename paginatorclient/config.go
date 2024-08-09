package paginatorclient

import (
	"github.com/spf13/viper"
	"log"
)

type PaginatorClientConfig struct {
	CrudClientURL string
	AuthToken     string
}

var config PaginatorClientConfig

func InitConfig() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		// Если файл не найден, продолжаем, т.к. переменные могут быть заданы в окружении
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Error reading config file, %s", err)
		}
	}

	viper.AutomaticEnv()

	viper.SetEnvPrefix("PAGINATOR_CLIENT_")
	viper.BindEnv("CRUD_CLIENT_URL")
	viper.BindEnv("AUTH_TOKEN")

	config.CrudClientURL = viper.GetString("CRUD_CLIENT_URL")
	if config.CrudClientURL == "" {
		log.Fatal("CRUD_CLIENT_URL is not set in the .env file or environment variables")
	}

	config.AuthToken = viper.GetString("AUTH_TOKEN")
	if config.AuthToken == "" {
		log.Fatal("AUTH_TOKEN is not set in the .env file or environment variables")
	}
}

func GetConfig() *PaginatorClientConfig {
	return &config
}
