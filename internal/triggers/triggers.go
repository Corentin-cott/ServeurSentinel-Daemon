package triggers

import (
	"fmt"
	"os"
	"regexp"
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
				// On utilise une expression régulière pour récupérer le nom du joueur
				re := regexp.MustCompile(`\[(\d{2}:\d{2}:\d{2})\] \[Server thread/INFO\]: (.+) joined the game`)
				matches := re.FindStringSubmatch(line)
				if len(matches) < 3 {
					fmt.Println("Erreur lors de la récupération du nom du joueur")
					return
				}
				fmt.Println("Player joined: ", matches[2])
				WriteToLogFile("/var/log/serversentinel/playerjoined.log", matches[2])
			},
		},
		{
			Condition: func(line string) bool {
				return strings.Contains(line, "lost connection: Disconnected")
			},
			Action: func(line string) {
				// On utilise une expression régulière pour récupérer le nom du joueur
				re := regexp.MustCompile(`\[(\d{2}:\d{2}:\d{2})\] \[Server thread/INFO\]: (.+) lost connection: Disconnected`)
				matches := re.FindStringSubmatch(line)
				if len(matches) < 3 {
					fmt.Println("Erreur lors de la récupération du nom du joueur")
					return
				}
				fmt.Println("Player disconnected: ", matches[2])
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
