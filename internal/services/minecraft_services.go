package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetPlayerUUID(playerName string) (string, error) {
	fmt.Println("\nGetting player UUID for " + playerName + " from Mojang API...")
	mojangAPIUrl := "https://api.mojang.com/users/profiles/minecraft/" + playerName
	resp, err := http.Get(mojangAPIUrl)
	if err != nil {
		return "", fmt.Errorf("FAILED TO SEND REQUEST TO MOJANG API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("FAILED TO GET PLAYER UUID, STATUS CODE: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("FAILED TO READ RESPONSE BODY: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("FAILED TO READ JSON RESPONSE: %v", err)
	}

	playerUUID, ok := result["id"].(string)
	if !ok || playerUUID == "" {
		return "", fmt.Errorf("FAILED TO GET PLAYER UUID: %v", result)
	}

	fmt.Println("Player UUID retrieved successfully : " + playerUUID + " for player name : " + playerName + "\n")
	return playerUUID, nil
}
