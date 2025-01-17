package triggers

// This file contains the ACTIONS functions for the triggers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Corentin-cott/ServeurSentinel/config"
)

// WriteToLogFile writes a line to a log file
func WriteToLogFile(logPath string, line string) error {
	// Open the log file
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("ERROR WHILE OPENING LOG FILE: %v", err)
	}
	defer file.Close()

	// Write the line to the log file
	_, err = file.WriteString(line + "\n")
	if err != nil {
		return fmt.Errorf("ERROR WHILE WRITING TO LOG FILE: %v", err)
	}
	return nil
}

// SendToDiscord sends a message to a Discord channel
func SendToDiscord(message string) error {
	// Get parameters from the configuration
	botToken := config.AppConfig.Bot.BotToken
	channelID := config.AppConfig.Bot.DiscordChannelID

	// Checks if one of the parameters is missing
	switch {
	case botToken == "" && channelID == "":
		return fmt.Errorf("ERROR: BOT TOKEN AND CHANNEL ID NOT SET")
	case botToken == "":
		return fmt.Errorf("ERROR: BOT TOKEN NOT SET")
	case channelID == "":
		return fmt.Errorf("ERROR: CHANNEL ID NOT SET")
	}

	// Prepare the request to the Discord API
	apiURL := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", channelID)
	type DiscordBotMessage struct {
		Content string `json:"content"`
	}

	// Serialize the message to JSON
	payload := DiscordBotMessage{Content: message}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ERROR WHILE SERIALISING DISCORD MESSAGE: %v", err)
	}

	// Create the HTTP request to send the message
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("ERROR WHILE CREATING REQUEST TO DISCORD: %v", err)
	}

	// Set the headers for the request
	req.Header.Set("Authorization", "Bot "+botToken)
	req.Header.Set("Content-Type", "application/json")

	// Finally, send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ERROR WHILE SENDING MESSAGE TO DISCORD: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("ERROR WHILE SENDING MESSAGE TO DISCORD, RESPONSE STATUS: %v", resp.Status)
	} else {
		return nil
	}
}

// SendToServer sends a message to a server
func SendToServer(serverID int, serverGame string, message string) error {
	// Not implemented yet
	return fmt.Errorf("ERROR: NOT IMPLEMENTED YET")
}
