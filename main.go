package main

import (
	"github.com/MiguelMachado-dev/disc-go-bot/config"
	"github.com/MiguelMachado-dev/disc-go-bot/database"
	"github.com/MiguelMachado-dev/disc-go-bot/discord"
)

func main() {
	config.Init()

	log := config.NewLogger("main")

	// Initialize the database
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize the discord Bot
	discord.Init()

	log.Fatalf("No service to run. Exiting...")
}
