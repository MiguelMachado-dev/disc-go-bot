package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/MiguelMachado-dev/disc-go-bot/config"
)

func GetTwitchAccessToken() (string, error) {
	TwitchClientID := config.GetEnv().TWITCH_CLIENT_ID
	TwitchClientSecret := config.GetEnv().TWITCH_CLIENT_SECRET

	url := "https://id.twitch.tv/oauth2/token"
	payload := strings.NewReader(fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials", TwitchClientID, TwitchClientSecret))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if accessToken, ok := result["access_token"].(string); ok {
		return accessToken, nil
	}
	return "", fmt.Errorf("access token not found in response: %v", result)
}
