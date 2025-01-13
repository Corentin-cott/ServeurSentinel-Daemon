package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Config struct {
	BotToken         string         `json:"botToken"`
	LogPath          string         `json:"logPath"`
	DiscordChannelID string         `json:"discordChannelID"`
	DB               DatabaseConfig `json:"db"`
}

var AppConfig Config

// LoadConfig charge la configuration d'un fichier JSON en paramètre
func LoadConfig(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("error opening configuration file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&AppConfig); err != nil {
		return fmt.Errorf("error decoding configuration: %v", err)
	}

	fmt.Printf("Configuration loaded successfully\n")
	return nil
}
