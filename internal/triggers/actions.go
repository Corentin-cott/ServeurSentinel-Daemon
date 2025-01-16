package triggers

// This file contains the ACTIONS functions for the triggers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/Corentin-cott/ServeurSentinel/config"
	"github.com/Corentin-cott/ServeurSentinel/internal/db"
)

// WriteToLogFile writes a line to a log file
func WriteToLogFile(logPath, line string) error {
	// Open the log file
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("erreur d'ouverture du fichier log : %v", err)
	}
	defer file.Close()

	// Write the line to the log file
	_, err = file.WriteString(line + "\n")
	if err != nil {
		return fmt.Errorf("erreur d'écriture dans le fichier log : %v", err)
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
		fmt.Println("Bot token and channel ID not set. Skipping Discord message.")
		return fmt.Errorf("ERROR: BOT TOKEN AND CHANNEL ID NOT SET")
	case botToken == "":
		fmt.Println("Bot token not set. Skipping Discord message.")
		return fmt.Errorf("ERROR: BOT TOKEN NOT SET")
	case channelID == "":
		fmt.Println("Channel ID not set. Skipping Discord message.")
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

// PlayerJoinAction logs the player's connection
func PlayerJoinAction(line string) error {
	// Extract the player's name from the line
	re := regexp.MustCompile(`\[(\d{2}:\d{2}:\d{2})\] \[Server thread/INFO\]: (\w+) joined the game`)
	matches := re.FindStringSubmatch(line)
	if len(matches) < 3 {
		return fmt.Errorf("ERRROR WHILE EXTRACTING CONNECTION INFO")
	}

	playerName := matches[2] // Get the player's name from the regex match

	serverID := 1 // Temporary server ID, will be replaced by the real server ID later

	// Save the connection log to the database
	err := db.SaveConnectionLog(playerName, serverID)
	if err != nil {
		return fmt.Errorf("ERROR WHILE SAVING CONNECTION LOG: %v", err)
	}
	return nil
}
