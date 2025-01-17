package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/Corentin-cott/ServeurSentinel/config"
	"github.com/Corentin-cott/ServeurSentinel/internal/console"
	"github.com/Corentin-cott/ServeurSentinel/internal/db"
	"github.com/Corentin-cott/ServeurSentinel/internal/triggers"
)

func main() {

	// Commands handler, check if any arguments are passed
	var correctUsage = "correct usage: serversentinel <command> [arguments]"
	if len(os.Args) < 2 {
		fmt.Println("Starting the Server Sentinel daemon...")
	} else if len(os.Args) == 2 || len(os.Args) > 3 {
		fmt.Println("Invalid number of arguments, " + correctUsage)
		os.Exit(1)
	} else {
		fmt.Println("Valid number of arguments")

		// // Check if the second argument is an integer
		// var serverID int
		// if len(os.Args) > 1 {
		// 	var err error
		// 	serverID, err = strconv.Atoi(os.Args[2]) // Convert the third argument to an integer
		// 	if err != nil {
		// 		fmt.Println("Invalid server ID, " + correctUsage)
		// 		os.Exit(1)
		// 	}
		// }

		// // Command handler
		// command := os.Args[1]
		// switch command {
		// case "start-minecraft":
		// 	if len(os.Args) < 3 {
		// 		fmt.Println("Missing server ID, " + correctUsage)
		// 		os.Exit(1)
		// 	}

		// 	fmt.Println("Now starting the Minecraft server " + os.Args[2])
		// 	err := commands.StartMinecraftServer(serverID)
		// 	if err != nil {
		// 		fmt.Println("The Minecraft server could not be started: " + err.Error())
		// 	}
		// case "stop-minecraft":
		// 	if len(os.Args) < 3 {
		// 		fmt.Println("Missing server ID, " + correctUsage)
		// 		os.Exit(1)
		// 	}

		// 	fmt.Println("Not implemented yet")
		// default:
		// 	fmt.Println("Unknown command, " + correctUsage)
		// 	os.Exit(1)
		// }

		os.Exit(1)
	}

	// Load the configuration file
	err := config.LoadConfig("/opt/serversentinel/config.json")
	if err != nil {
		log.Fatalf("FATAL ERROR LOADING CONFIG JSON FILE: %v", err)
	}

	// Initialize the connection to the database
	err = db.ConnectToDatabase()
	if err != nil {
		log.Fatalf("FATAL ERROR TESTING DATABASE CONNECTION: %v", err)
	}

	// Log files handler, here get all log files in the directory
	logDirPath := "/opt/serversentinel/serverslog/" // Folder containing the log files
	logFiles, err := filepath.Glob(filepath.Join(logDirPath, "*.log"))
	if err != nil {
		log.Fatalf("FATAL ERROR WHEN GETTING LOG FILES: %v", err)
	}

	if len(logFiles) == 0 {
		log.Println("No log files found in the directory, did you forget to redirect the logs to the folder ?")
		return
	}

	// Create a list of triggers and create a wait group
	triggersList := triggers.GetTriggers()
	processLogFiles(logDirPath, triggersList)

	fmt.Println("Server Sentinel daemon stopped.")
}

// Function to process all log files in a directory
func processLogFiles(logDirPath string, triggersList []console.Trigger) {
	logFiles, err := filepath.Glob(filepath.Join(logDirPath, "*.log"))
	if err != nil {
		log.Fatalf("FATAL ERROR WHEN GETTING LOG FILES: %v", err)
	}

	if len(logFiles) == 0 {
		log.Println("No log files found in the directory, did you forget to redirect the logs to the folder?")
		return
	}

	// Create a wait group
	var wg sync.WaitGroup

	// Start a goroutine for each log file
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

	// Wait for all goroutines to finish
	wg.Wait()
}
