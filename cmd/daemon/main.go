package main

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/Corentin-cott/ServeurSentinel/internal/console"
)

func main() {
	logDirPath := "/opt/serversentinel/serverslog/" // Répertoire contenant les fichiers .log

	fmt.Println("Démarrage du daemon Server Sentinel...")

	// Récupérer tous les fichiers .log dans le répertoire
	logFiles, err := filepath.Glob(filepath.Join(logDirPath, "*.log"))
	if err != nil {
		log.Fatalf("Erreur lors de la recherche des fichiers log : %v", err)
	}

	if len(logFiles) == 0 {
		log.Println("Aucun fichier log trouvé dans le répertoire.")
		return
	}

	var wg sync.WaitGroup

	// Lancer un écouteur pour chaque fichier log
	for _, logFile := range logFiles {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			err := console.StartFileLogListener(file, console.ExampleAction)
			if err != nil {
				log.Printf("Erreur pour le fichier %s : %v\n", file, err)
			}
		}(logFile)
	}

	// Attendre que tous les écouteurs soient terminés
	wg.Wait()
	fmt.Println("Arrêt du daemon Server Sentinel.")
}
