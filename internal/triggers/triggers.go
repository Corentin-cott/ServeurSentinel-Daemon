package triggers

// This file contains the TRIGGERS functions

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Corentin-cott/ServeurSentinel/internal/console"
)

var (
	playerJoinedRegex       = regexp.MustCompile(`\[(\d{2}:\d{2}:\d{2})\] \[Server thread/INFO\]: (.+) joined the game`)
	playerDisconnectedRegex = regexp.MustCompile(`\[(\d{2}:\d{2}:\d{2})\] \[Server thread/INFO\]: (.+) lost connection: Disconnected`)
)

func GetTriggers() []console.Trigger {
	return []console.Trigger{
		{
			Condition: func(line string) bool {
				return strings.Contains(line, "joined the game")
			},
			Action: func(line string) {
				// On utilise une expression régulière pour récupérer le nom du joueur
				matches := playerJoinedRegex.FindStringSubmatch(line)
				if len(matches) < 3 {
					fmt.Println("Erreur lors de la récupération du nom du joueur")
					return
				}
				// fmt.Println(matches[2] + " à rejoint le serveur")
				SendToDiscord(matches[2] + " à rejoint le serveur")
				PlayerJoinAction(line)
				WriteToLogFile("/var/log/serversentinel/playerjoined.log", matches[2])
			},
		},
		{
			Condition: func(line string) bool {
				return strings.Contains(line, "lost connection: Disconnected")
			},
			Action: func(line string) {
				// On utilise une expression régulière pour récupérer le nom du joueur
				matches := playerDisconnectedRegex.FindStringSubmatch(line)
				if len(matches) < 3 {
					fmt.Println("Erreur lors de la récupération du nom du joueur")
					return
				}
				// fmt.Println(matches[2] + " à quitté le serveur")
				SendToDiscord(matches[2] + " à quitté le serveur")
				WriteToLogFile("/var/log/serversentinel/playerdisconnected.log", matches[2])
			},
		},
		{
			Condition: func(line string) bool {
				return strings.HasPrefix(line, "[ERROR]")
			},
			Action: func(line string) {
				fmt.Println("Error detected: ", line)
				fmt.Println("Sending error to the error log...")
				WriteToLogFile("/var/log/serversentinel/errors.log", line)
			},
		},
	}
}
