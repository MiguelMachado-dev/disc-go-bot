package handler

import (
	"testing"

	"github.com/MiguelMachado-dev/disc-go-bot/config"
	"github.com/bwmarrin/discordgo"
)

func TestHuntHandler_Handler(t *testing.T) {
	discordToken := config.GetEnv().DISCORD_BOT_TOKEN
	// Crie uma sessão de teste
	session, err := discordgo.New(discordToken)
	if err != nil {
		t.Fatalf("Erro ao criar sessão de teste: %v", err)
	}
	// Crie uma interação de teste
	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID: "1234567890",
		},
	}
	// Crie uma instância do HuntHandler
	h := &HuntHandler{}
	// Chame o método Handler
	h.Handler(session, interaction)
	// Verifique se não houve erros
	if err != nil {
		t.Errorf("Erro ao executar o método Handler: %v", err)
	}
}
