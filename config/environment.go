package config

import (
	"errors"
	"os"
	"reflect"

	_ "github.com/joho/godotenv/autoload" // Load .env file automatically
)

type Environment struct {
	DISCORD_BOT_TOKEN string
	HUNT_CHANNEL_ID   string
}

var environment Environment

func initializeEnvironment() error {
	envType := reflect.TypeOf(environment)
	envValue := reflect.ValueOf(&environment).Elem()

	for i := 0; i < envType.NumField(); i++ {
		field := envType.Field(i)
		envVar := field.Name

		// Get environment variable value
		value := os.Getenv(envVar)

		// Check if the environment variable is set, otherwise return an error
		if value == "" {
			return errors.New("required environment variable " + envVar + " not set")
		}

		envValue.FieldByName(envVar).SetString(value)
	}

	return nil
}

func GetEnv() *Environment {
	return &environment
}
