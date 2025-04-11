package config

import (
	"errors"
	"os"
	"reflect"

	_ "github.com/joho/godotenv/autoload" // Load .env file automatically
)

type Environment struct {
	DISCORD_BOT_TOKEN    string
	COMMANDS_CHANNEL_ID  string
	TWITCH_CLIENT_ID     string
	TWITCH_CLIENT_SECRET string
	ENCRYPTION_KEY       string
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

		// Check if the environment variable is set
		if value == "" {
			// No special handling for ENCRYPTION_KEY, require it in .env
			return errors.New("required environment variable " + envVar + " not set. Please add it to your .env file")
		}

		envValue.FieldByName(envVar).SetString(value)
	}

	return nil
}

func GetEnv() *Environment {
	return &environment
}
