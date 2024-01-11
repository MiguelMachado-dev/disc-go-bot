package handler

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

type DeleteHandler struct{}

func (h *DeleteHandler) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "delete-messages",
		Description: "Delete last 100 messages from current channel",
	}
}

func (h *DeleteHandler) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if the user has admin permissions
	perms, _ := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	if perms&discordgo.PermissionAdministrator == 0 {
		// User does not have admin permissions
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You do not have permission to use this command.",
			},
		})
		return
	}

	// Get the number of messages to delete
	numMessages := 100 // default value

	// Respond that the deletion process is starting
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Starting to delete messages...",
		},
	})

	// Delete the messages
	messages, _ := s.ChannelMessages(i.ChannelID, numMessages, "", "", "")
	for _, message := range messages {
		s.ChannelMessageDelete(i.ChannelID, message.ID)
	}

	// Send a message (not a response to the interaction) indicating that the deletion process has completed
	s.ChannelMessageSend(i.ChannelID, "Deleted "+strconv.Itoa(len(messages))+" messages.")
}

func DeleteMessagesTicker(s *discordgo.Session, channelID string, hours int) {
	log.Infof("Starting to delete messages every %d hours", hours)
	ticker := time.NewTicker(time.Duration(hours) * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		deleteMessages(s, channelID, 100) // 100 is the maximum number of messages that can be fetched at once
	}
}

func deleteMessages(s *discordgo.Session, channelID string, numMessages int) {
	log.Infof("Deleting %d messages from channel %s", numMessages, channelID)
	// Fetch the messages
	messages, err := s.ChannelMessages(channelID, numMessages, "", "", "")
	if err != nil {
		log.Errorf("error fetching messages: %v", err)
		return
	}

	// Delete the messages
	for _, message := range messages {
		err := s.ChannelMessageDelete(channelID, message.ID)
		if err != nil {
			log.Errorf("error deleting message: %v", err)
			continue
		}
	}

	// If there are still messages left, delete them recursively
	if len(messages) == numMessages {
		deleteMessages(s, channelID, numMessages)
	}
}
