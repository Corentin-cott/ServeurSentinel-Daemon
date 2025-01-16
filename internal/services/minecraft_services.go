package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetMinecraftPlayerUUID gets the UUID of a Minecraft player by their username
func GetMinecraftPlayerUUID(playerName string) (string, error) {
	fmt.Println("\nGetting Minecraft player UUID for player " + playerName + "...")

	// Send a request to the Mojang API to get the player UUID by their username
	APIUrl := "https://api.mojang.com/users/profiles/minecraft/" + playerName
	resp, err := http.Get(APIUrl)
	if err != nil {
		return "", fmt.Errorf("FAILED TO SEND REQUEST TO MOJANG API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // API returns an error
		return "", fmt.Errorf("FAILED TO GET PLAYER UUID, STATUS CODE: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil { // Failed to read response body
		return "", fmt.Errorf("FAILED TO READ RESPONSE BODY: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil { // Failed to parse JSON response
		return "", fmt.Errorf("FAILED TO READ JSON RESPONSE: %v", err)
	}

	playerUUID, ok := result["id"].(string)
	if !ok || playerUUID == "" { // Failed to find player UUID in response
		return "", fmt.Errorf("FAILED TO GET PLAYER UUID: %v", result)
	}

	fmt.Println("Player UUID retrieved successfully : " + playerUUID + " for player name : " + playerName + "\n")
	return playerUUID, nil
}
