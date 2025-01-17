package commands

import (
	"fmt"
	"os/exec"

	"github.com/Corentin-cott/ServeurSentinel/internal/db"
)

// StartMinecraftServer starts a Minecraft server in a tmux session
func StartMinecraftServer(serverID int) error {
	tmuxSession := fmt.Sprintf("minecraft_%d", serverID)

	// Check if the server exists in the database
	serverName, err := db.GetServerNameByID(serverID)
	if err != nil {
		return fmt.Errorf("FAILED to get server name by ID: %v", err)
	}

	// Command to start a Minecraft server : tmux new-session -s ${id}_${nom_serv} 'java -Xmx1024M -Xms1024M -jar server.jar nogui 2>&1 || tee /opt/serversentinel/serverslog/${id}.log'
	cmd := exec.Command("tmux", "new-session", "-d", "-s", tmuxSession, "java -Xms1G -Xmx1G -jar server.jar nogui")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("FAILED to start Minecraft server '%s' in tmux session: '%s': %v", serverName, tmuxSession, err)
	}

	fmt.Printf("Minecraft server '%s' started in tmux session: '%s'\n", serverName, tmuxSession)
	return nil
}
