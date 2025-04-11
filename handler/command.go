package handler

import "github.com/bwmarrin/discordgo"

type CommandHandler interface {
	Command() *discordgo.ApplicationCommand
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var CommandHandlers = []CommandHandler{
	&DeleteHandler{},
	&LiveHandler{},
	&PriceHandler{},
	&AskHandler{},
	&SetGeminiKeyHandler{},
}
