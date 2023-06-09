package main

import (
	"github.com/MiguelMachado-dev/disc-go-bot/config"
	"github.com/MiguelMachado-dev/disc-go-bot/discord"
)

func main() {
	config.Init()

	log := config.NewLogger("main")

	// Initialize the discord Bot
	discord.Init()

	log.Fatalf("No service to run. Exiting...")
}
