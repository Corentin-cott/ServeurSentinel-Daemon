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

// Écrit une ligne dans un fichier log
func WriteToLogFile(logPath, line string) error {
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("erreur d'ouverture du fichier log : %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(line + "\n")
	if err != nil {
		return fmt.Errorf("erreur d'écriture dans le fichier log : %v", err)
	}
	return nil
}

// Envoi un message à un serveur Discord
func SendToDiscord(message string) {
	botToken := config.AppConfig.Bot.BotToken
	channelID := config.AppConfig.Bot.DiscordChannelID

	switch {
	case botToken == "" && channelID == "":
		fmt.Println("Bot token and channel ID not set. Skipping Discord message.")
		return
	case botToken == "":
		fmt.Println("Bot token not set. Skipping Discord message.")
		return
	case channelID == "":
		fmt.Println("Channel ID not set. Skipping Discord message.")
		return
	}

	apiURL := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", channelID)

	type DiscordBotMessage struct {
		Content string `json:"content"`
	}

	payload := DiscordBotMessage{Content: message}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Erreur lors de la sérialisation du message Discord : %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Printf("Erreur lors de la création de la requête : %v\n", err)
		return
	}

	req.Header.Set("Authorization", "Bot "+botToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Erreur lors de l'envoi à Discord : %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		fmt.Printf("Erreur lors de l'envoi à Discord : Status %d\n", resp.StatusCode)
	} else {
		fmt.Println("Message envoyé à Discord avec succès.")
	}
}

func PlayerJoinAction(line string) {
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
