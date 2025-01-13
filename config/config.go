package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	BotToken         string `json:"botToken"`
	LogPath          string `json:"logPath"`
	DiscordChannelID string `json:"discordChannelID"`
}

var AppConfig Config

// LoadConfig charge la configuration d'un fichier JSON en paramètre
func LoadConfig(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier de configuration : %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&AppConfig); err != nil {
		return fmt.Errorf("erreur lors du décodage de la configuration : %v", err)
	}

	return nil
}
