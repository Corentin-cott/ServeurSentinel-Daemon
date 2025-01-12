package console

import (
	"bufio"
	"fmt"
	"io"
	"os"

	// "regexp"
	"strings"
	"time"
)

type Trigger struct {
	Condition func(string) bool // Condition du trigger
	Action    func(string)      // Function à exécuter si la condition est vraie
}

// StartFileLogListener démarre un écouteur en temps réel pour un fichier log
func StartFileLogListener(logFilePath string, triggers []Trigger) error {
	file, err := os.Open(logFilePath)
	if err != nil {
		return fmt.Errorf("impossible d'ouvrir le fichier log : %v", err)
	}
	defer file.Close()

	// Position le curseur à la fin du fichier
	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("impossible de se positionner à la fin du fichier : %v", err)
	}

	fmt.Printf("Écoute des logs en temps réel démarrée : %s\n", logFilePath)

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return fmt.Errorf("erreur lors de la lecture du fichier log : %v", err)
		}

		line = strings.TrimSpace(line)
		if line != "" {
			for _, trigger := range triggers {
				if trigger.Condition(line) {
					trigger.Action(line)
				}
			}
		}
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
