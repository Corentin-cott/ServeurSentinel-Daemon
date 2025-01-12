package main

import (
	"fmt"
	"log"

	"github.com/Corentin-cott/ServeurSentinel/internal/console"
)

func main() {
	logFilePath := "/opt/serversentinel/serverslog/primaire.log"

	fmt.Println("Démarrage du daemon Server Sentinel...")

	// Démarrer l'écoute des logs avec le trigger ExampleAction
	err := console.StartFileLogListener(logFilePath, console.ExampleAction)
	if err != nil {
		log.Fatalf("Erreur lors du lancement de l'écoute des logs : %v", err)
	}
}
