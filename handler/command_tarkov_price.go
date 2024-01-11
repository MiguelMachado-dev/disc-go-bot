package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type PriceHandler struct{}

func (h *PriceHandler) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "price",
		Description: "Check Flea Market price for an item",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item",
				Description: "Name of the item to check",
				Required:    true,
			},
		},
	}
}

func relativeTime(timeStr string) string {
	// Parse the time string
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return ""
	}

	// Calculate the difference
	now := time.Now().UTC()
	diff := now.Sub(parsedTime)

	// Format the difference
	if diff.Hours() < 24 {
		return fmt.Sprintf("%dh ago", int(diff.Hours()))
	} else if diff.Hours() < 48 {
		return "Yesterday"
	} else {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	}
}

func (h *PriceHandler) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Extract the item name
	in := i.ApplicationCommandData().Options[0].StringValue()

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	body := strings.NewReader(`{"query": "{ items(name: \"` + in + `\") {name avg24hPrice baseImageLink link basePrice wikiLink updated sellFor { vendor { name } price }} }"}`)
	req, err := http.NewRequest("POST", "https://api.tarkov.dev/graphql", body)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := netClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Unmarshal the response into a struct
	var response struct {
		Data struct {
			Items []struct {
				Name          string `json:"name"`
				Avg24hPrice   int    `json:"avg24hPrice"`
				BaseImageLink string `json:"baseImageLink"`
				Link          string `json:"link"`
				BasePrice     int    `json:"basePrice"`
				WikiLink      string `json:"wikiLink"`
				Updated       string `json:"updated"`
				SellFor       []struct {
					Vendor struct {
						Name string `json:"name"`
					} `json:"vendor"`
					Price int `json:"price"`
				} `json:"sellFor"`
			} `json:"items"`
		} `json:"data"`
	}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		log.Fatalln(err)
	}

	// Check if the item exists
	if len(response.Data.Items) == 0 {
		errMessage := fmt.Sprintf("Nenhum item encontrado sobre '%s'", in)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: errMessage,
			},
		})
		return
	}

	itemName := response.Data.Items[0].Name

	itemPrice := response.Data.Items[0].Avg24hPrice
	p := message.NewPrinter(language.English)
	itemPriceString := p.Sprintf("%d ₽", itemPrice)

	itemImage := response.Data.Items[0].BaseImageLink
	itemLink := response.Data.Items[0].Link

	// Add dynamic sellFor fields
	sellForFields := make([]*discordgo.MessageEmbedField, len(response.Data.Items[0].SellFor))
	for i, sellFor := range response.Data.Items[0].SellFor {
		sellForFields[i] = &discordgo.MessageEmbedField{
			Name:   sellFor.Vendor.Name,
			Value:  p.Sprintf("%d ₽", sellFor.Price),
			Inline: true,
		}
	}

	defer resp.Body.Close()

	avg24hPriceField := &discordgo.MessageEmbedField{
		Name:   "Average 24h price",
		Value:  itemPriceString,
		Inline: true,
	}

	// convert string Updated to date
	itemUpdated := response.Data.Items[0].Updated
	relativeDate := relativeTime(itemUpdated)

	lastUpdatedField := &discordgo.MessageEmbedField{
		Name:   "Last updated",
		Value:  relativeDate,
		Inline: true,
	}

	wikiLinkField := &discordgo.MessageEmbedField{
		Name:  "Wiki link",
		Value: response.Data.Items[0].WikiLink,
	}

	basePriceField := &discordgo.MessageEmbedField{
		Name:   "Base price",
		Value:  p.Sprintf("%d ₽", response.Data.Items[0].BasePrice),
		Inline: true,
	}

	// create an embed message and respond to the interaction
	embed := &discordgo.MessageEmbed{
		Title:       itemName,
		Description: "Price check for " + itemName,
		Color:       0x00ff00,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: itemImage,
		},
		URL: itemLink,
		Fields: []*discordgo.MessageEmbedField{
			basePriceField,
			avg24hPriceField,
			lastUpdatedField,
			wikiLinkField,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Data provided by tarkov.dev and compiled by Migtito",
		},
	}

	// Add dynamic sellFor fields
	embed.Fields = append(embed.Fields, sellForFields...)

	// Respond to the interaction with the embed
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
