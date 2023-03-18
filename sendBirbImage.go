package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

func SendBirbImage(s *discordgo.Session, channelID string) {
	// Send a "loading" message to the channel
	loadingMessage, err := s.ChannelMessageSend(channelID, "Loading birb image...")
	if err != nil {
		fmt.Println("Error sending loading message:", err)
		return
	}

	url := "https://random.birb.pw/tweet.json"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making API request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading API response:", err)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return
	}

	imagePath := result["file"].(string)
	imageURL := "https://random.birb.pw/img/" + imagePath

	// Download the image data using another HTTP GET request
	imageResp, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("Error downloading image:", err)
		return
	}
	defer imageResp.Body.Close()

	// Use the image data as the Reader for discordgo.File
	attachment := discordgo.File{
		Name:   "birb_image.jpg",
		Reader: imageResp.Body,
	}

	// Send the image and delete the "loading" message
	_, err = s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Files: []*discordgo.File{&attachment},
	})

	if err != nil {
		fmt.Println("Error sending dog image:", err)
		return
	}

	err = s.ChannelMessageDelete(channelID, loadingMessage.ID)
	if err != nil {
		fmt.Println("Error deleting loading message:", err)
	}
}
