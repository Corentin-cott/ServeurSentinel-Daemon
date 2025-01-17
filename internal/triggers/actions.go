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

// Envoi un message à un serveur Discord
func SendToDiscord(message string) {
	if !config.AppConfig.EnableBot {
		fmt.Println("Bot messages are disabled. Skipping Discord message.")
		return
	}

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

func PlayerJoinAction(line string) {
	if !config.AppConfig.EnableDatabase {
		fmt.Println("Database insertions are disabled. Skipping player join log.")
		return
	}

	// Extraire les informations du joueur
	re := regexp.MustCompile(`\[(\d{2}:\d{2}:\d{2})\] \[Server thread/INFO\]: (\w+) joined the game`)
	matches := re.FindStringSubmatch(line)
	if len(matches) < 3 {
		fmt.Println("Erreur lors de l'extraction des informations de connexion")
		return
	}

	playerName := matches[2] // Le UUID du joueur est le deuxième groupe

	// Utilisez la nouvelle fonction pour vérifier et enregistrer le joueur
	serverID := 1 // Vous pouvez utiliser l'ID de votre serveur ici

	// Enregistrer la connexion
	err := db.SaveConnectionLog(playerName, serverID)
	if err != nil {
		fmt.Printf("Erreur lors de l'enregistrement du log de connexion : %v\n", err)
	}
}

// SendToServer sends a message to a server
func SendToServer(serverID int, serverGame string, message string) error {
	// Not implemented yet
	return fmt.Errorf("ERROR: NOT IMPLEMENTED YET")
}
