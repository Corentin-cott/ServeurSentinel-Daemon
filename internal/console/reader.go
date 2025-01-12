package console

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

// StartFileLogListener démarre un écouteur en temps réel pour un fichier log
func StartFileLogListener(logFilePath string, handleLine func(string)) error {
	file, err := os.Open(logFilePath)
	if err != nil {
		return fmt.Errorf("impossible d'ouvrir le fichier log : %v", err)
	}
	defer file.Close()

	// Positionner le curseur à la fin du fichier pour ne pas lire les lignes précédentes
	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("impossible de se positionner à la fin du fichier : %v", err)
	}

	fmt.Printf("Écoute des logs en temps réel démarrée : %s\n", logFilePath)

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// Si une erreur arrive, on attend avant de réessayer
			if err.Error() == "EOF" {
				time.Sleep(100 * time.Millisecond) // Petite pause pour attendre de nouvelles données
				continue
			}
			return fmt.Errorf("erreur lors de la lecture du fichier log : %v", err)
		}

		// Supprimer les espaces autour de la ligne avant de la traiter
		line = strings.TrimSpace(line)
		if line != "" {
			handleLine(line)
		}
	}
}

// ExampleAction est un exemple d'action déclenchée par le log
func ExampleAction(line string) {
	if strings.Contains(line, "joined the game") {
		// Dans cet exemple, notre trigger sera la ligne "joined the game", donc dès qu'un joueur se connecte à un serveur Minecraft
		// On va utiliser un regex pour récupérer le nom du joueur
		re := regexp.MustCompile(`(?P<player>\w+) joined the game`)
		matches := re.FindStringSubmatch(line)
		if len(matches) < 2 {
			fmt.Println("Impossible de récupérer le nom du joueur")
			return
		}
		player := matches[1]

		// On peut maintenant faire ce que l'on veut avec le nom du joueur, par exemple l'afficher dans la console
		fmt.Printf("Le joueur %s a rejoint le serveur\n", player)

		// Mais aussi écrire le nom du joueur dans un fichier log pour surveiller les connexions comme un mani- euh, un professionnel. Un professionnel.
		writeToLogFile("/var/log/serversentinel/playerjoined.log", line)
	}
}

// Écrit une ligne dans un fichier log
func writeToLogFile(logPath, line string) error {
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
