package triggers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Corentin-cott/ServeurSentinel/internal/console"
)

// Variables to store the regex patterns
var (
	playerJoinedRegex       = regexp.MustCompile(`\[(\d{2}:\d{2}:\d{2})\] \[Server thread/INFO\]: (.+) joined the game`)
	playerDisconnectedRegex = regexp.MustCompile(`\[(\d{2}:\d{2}:\d{2})\] \[Server thread/INFO\]: (.+) lost connection: Disconnected`)
)

// GetTriggers returns the list of triggers filtered by names
func GetTriggers(selectedTriggers []string) []console.Trigger {
	// All available triggers
	allTriggers := []console.Trigger{
		{
			// This is an example trigger, use it as a template to create new triggers
			Name: "ExampleTrigger",
			Condition: func(line string) bool {
				// Here you can define the condition that will trigger the action, you're most probably looking for a specific string in the server log
				return strings.Contains(line, "whatever line you're looking for here")
			},
			Action: func(line string) {
				// Here you can define the action that will be executed
				fmt.Println("Example trigger action executed")
			},
		},
		{
			// This trigger is used to detect when a player joins a Minecraft server
			Name: "PlayerJoinedMinecraftServer",
			Condition: func(line string) bool {
				return strings.Contains(line, "joined the game")
			},
			Action: func(line string) {
				matches := playerJoinedRegex.FindStringSubmatch(line)
				if len(matches) < 3 {
					fmt.Println("ERROR WHILE EXTRACTING JOINED PLAYER NAME")
					return
				}
				SendToDiscord(matches[2] + " à rejoint le serveur")
				PlayerJoinAction(line)
				WriteToLogFile("/var/log/serversentinel/playerjoined.log", matches[2])
			},
		},
		{
			// This trigger is used to detect when a player disconnects from a Minecraft server
			Name: "PlayerDisconnectedMinecraftServer",
			Condition: func(line string) bool {
				return strings.Contains(line, "lost connection: Disconnected")
			},
			Action: func(line string) {
				matches := playerDisconnectedRegex.FindStringSubmatch(line)
				if len(matches) < 3 {
					fmt.Println("ERROR WHILE EXTRACTING DISCONNECTED PLAYER NAME")
					return
				}
				SendToDiscord(matches[2] + " à quitté le serveur")
				WriteToLogFile("/var/log/serversentinel/playerdisconnected.log", matches[2])
			},
		},
	}

	// If no specific triggers are requested, return all triggers
	if len(selectedTriggers) == 0 {
		return allTriggers
	}

	// Filter triggers based on the selected names
	var filteredTriggers []console.Trigger
	for _, trigger := range allTriggers {
		for _, name := range selectedTriggers {
			if trigger.Name == name {
				filteredTriggers = append(filteredTriggers, trigger)
				break
			}
		}
	}

	return filteredTriggers
}
