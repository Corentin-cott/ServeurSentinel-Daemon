package triggers

import (
	"fmt"
	"strings"

	"github.com/Corentin-cott/ServeurSentinel/internal/console"
)

// GetTriggers renvoie la liste de tous les triggers
func GetTriggers() []console.Trigger {
	return []console.Trigger{
		{
			Condition: func(line string) bool {
				return strings.Contains(line, "joined the game")
			},
			Action: func(line string) {
				fmt.Println("Player joined: ", line)
				console.WriteToLogFile("/var/log/serversentinel/playerjoined.log", line)
			},
		},
		{
			Condition: func(line string) bool {
				return strings.Contains(line, "ALERT")
			},
			Action: func(line string) {
				fmt.Println("[ALERT]: ", line)
				console.WriteToLogFile("/var/log/serversentinel/alerts.log", line)
			},
		},
		{
			Condition: func(line string) bool {
				return strings.Contains(line, "disconnected")
			},
			Action: func(line string) {
				fmt.Println("Player disconnected: ", line)
				console.WriteToLogFile("/var/log/serversentinel/disconnected.log", line)
			},
		},
		{
			Condition: func(line string) bool {
				return strings.HasPrefix(line, "[ERROR]")
			},
			Action: func(line string) {
				fmt.Println("Error detected: ", line)
				console.WriteToLogFile("/var/log/serversentinel/errors.log", line)
			},
		},
	}
}
