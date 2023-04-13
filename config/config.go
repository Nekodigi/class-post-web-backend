package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	ClientId     string
	ProjectId    string
	RefreshToken string
	ClientSecret string
}

var config *Config

func Load() *Config {
	err := godotenv.Load("dev.env")
	if err == nil {
		log.Infoln("Load dev.env file for local dev")
	}

	if config == nil {
		if os.Getenv("CLIENT_ID") == "" { //other env value might not set as well
			log.Fatalln("CLIENT_ID is not set:")
		}

		config = &Config{
			ClientId:  os.Getenv("CLIENT_ID"),
			ProjectId: os.Getenv("PROJECT_ID"),

			RefreshToken: os.Getenv("REFRESH_TOKEN"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
		}
	}
	fmt.Printf("LOAD CONFIG:%v", config)
	return config
}
