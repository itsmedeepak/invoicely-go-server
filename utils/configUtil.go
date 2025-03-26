package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/joho/godotenv"
)

var sesClient *sesv2.Client
var envLoaded bool

func loadEnv() {
	if !envLoaded {
		if envErr := godotenv.Load(); envErr != nil {
			fmt.Println("Unable to load .env file")
		}
		envLoaded = true
	}
}

func GetEnv(key string) string {
	loadEnv()
	token := os.Getenv(key)
	if token == "" {
		fmt.Printf("Empty `%s` in .env file\n", key)
	}
	return token
}

func GetSesClient() (*sesv2.Client, error) {
	if sesClient == nil {
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
		if err != nil {
			fmt.Printf("Unable to load SDK config: %v\n", err)
			return nil, err
		}
		sesClient = sesv2.NewFromConfig(cfg)
	}
	return sesClient, nil
}
