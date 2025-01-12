package main

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/Corentin-cott/ServeurSentinel/internal/console"
	"github.com/Corentin-cott/ServeurSentinel/internal/triggers"
)

func main() {
	logDirPath := "/opt/serversentinel/serverslog/" // Dossier contenant les fichiers log des serveurs

	fmt.Println("Starting the Server Sentinel daemon...")

	// Récupére tous les fichiers .log du dossier
	logFiles, err := filepath.Glob(filepath.Join(logDirPath, "*.log"))
	if err != nil {
		log.Fatalf("Error searching for log files: %v", err)
	}

	if len(logFiles) == 0 {
		log.Println("No log files found in the directory.")
		return
	}

	// Crée la liste des triggers
	triggersList := triggers.GetTriggers()

	var wg sync.WaitGroup

	// Commence un écouteur pour chaque fichier log
	for _, logFile := range logFiles {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			err := console.StartFileLogListener(file, triggersList)
			if err != nil {
				log.Printf("Error with file %s: %v\n", file, err)
			}
		}(logFile)
	}

	// Attend que tous les écouteurs soient terminés
	wg.Wait()
	fmt.Println("Server Sentinel daemon stopped.")
}
