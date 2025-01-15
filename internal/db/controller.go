package db

import (
	"database/sql"
	"fmt"

	"github.com/Corentin-cott/ServeurSentinel/config"
	"github.com/Corentin-cott/ServeurSentinel/internal/services"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// ConnectToDatabase initialises the connection to the MySQL database
func ConnectToDatabase() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.AppConfig.DB.User,
		config.AppConfig.DB.Password,
		config.AppConfig.DB.Host,
		config.AppConfig.DB.Port,
		config.AppConfig.DB.Name,
	)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %v", err)
	}

	// Verifier la connexion
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	fmt.Println("Successfully connected to the database.")
	return nil
}

// CheckAndInsertPlayer vérifie si un joueur existe, sinon il l'ajoute
func CheckAndInsertPlayer(playerName string, jeu string) (int, error) {
	// Utilisation du service pour récupérer l'UUID du joueur
	playerAcountID, err := services.GetPlayerUUID(playerName)
	if err != nil {
		return -1, fmt.Errorf("FAILED TO GET PLAYER ACCOUNT ID: %v", err)
	}

	// Vérifier si le joueur existe déjà en utilisant GetPlayerIdByAccoundId(accountId string)
	fmt.Println("Checking if player exists...")
	playerID, _ := GetPlayerIdByAccountId(playerAcountID)
	if playerID != -1 {
		fmt.Printf("PLAYER WITH ACCOUNT ID %s ALREADY EXISTS\n", playerAcountID)
		return playerID, nil
	}

	// Si le joueur n'existe pas, ajouter un nouveau joueur dans la table joueurs avec serveur_id = 1 [DEBUG]
	fmt.Println("Player does not exist. Inserting new player...")
	insertQuery := "INSERT INTO joueurs (utilisateur_id, jeu, compte_id, premiere_co, derniere_co) VALUES (NULL, ?, ?, NOW(), NOW())"
	_, err = db.Exec(insertQuery, "Minecraft", playerAcountID)
	if err != nil {
		return -1, fmt.Errorf("FAILED TO INSERT PLAYER: %v", err)
	}
	fmt.Println("Player successfully inserted !")

	// Récupérer l'id du joueur nouvellement inséré
	playerID, err = GetPlayerIdByAccountId(playerAcountID)
	if err != nil {
		return -1, fmt.Errorf("FAILED TO GET PLAYER ID: %v", err)
	} else if playerID == -1 {
		return -1, fmt.Errorf("PLAYER ID NOT FOUND")
	}

	return playerID, nil
}

// SaveConnectionLog sauvegarde un log de connexion dans la base de données
func SaveConnectionLog(playerName string, serverID int) error {
	game, _ := GetServerGameById(serverID)
	userID, err := CheckAndInsertPlayer(playerName, game)
	if err != nil {
		return fmt.Errorf("FAILED TO CHECK OR INSERT PLAYER: %v", err)
	}

	// Récupérer l'id du joueur
	playerID, err := GetPlayerIdByAccountId(userID)
	if err != nil {
		return fmt.Errorf("FAILED TO GET PLAYER ID: %v", err)
	} else if playerID == -1 {
		return fmt.Errorf("PLAYER ID NOT FOUND")
	}

	// Enregistrer le log de connexion
	insertQuery := `INSERT INTO joueurs_connections_log (serveur_id, joueur_id, date) VALUES (?, ?, NOW())`
	fmt.Println("Inserting connection log for player", playerID)
	_, err = db.Exec(insertQuery, serverID, playerID)
	if err != nil {
		return fmt.Errorf("FAILED TO INSERT CONNECTION LOG: %v", err)
	}

	fmt.Println("Connection log successfully saved.")
	return nil
}

// Getter pour récupéré l'id d'un joueur à partir de son id de compte. UUID pour les joueurs Minecraft, par exemple.
func GetPlayerIdByAccountId(accountId any) (int, error) {
	query := "SELECT id FROM joueurs WHERE compte_id = ?"
	var playerID int

	// Exécuter la requête et scanner le résultat dans la variable userId
	err := db.QueryRow(query, accountId).Scan(&playerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
		// Toute autre erreur
		fmt.Println("FAILED TO GET PLAYER ID:", err)
		return -1, fmt.Errorf("FAILED TO GET PLAYER ID: %v", err)
	}

	strPlayerID := fmt.Sprintf("%d", playerID)
	fmt.Println("Player ID retrieved successfully : "+strPlayerID+" for account ID : ", accountId)
	return playerID, nil
}

func GetPlayerAccountIdByPlayerName(playerName string, jeu string) (string, error) {
	switch jeu {
	case "Minecraft":
		return services.GetPlayerUUID(playerName)
	default:
		return "", fmt.Errorf("UNKNOWN GAME: %s", jeu)
	}
}

func GetServerGameById(serverID int) (string, error) {
	query := "SELECT jeu FROM serveurs WHERE id = ?"
	var jeu string

	err := db.QueryRow(query, serverID).Scan(&jeu)
	if err != nil {
		return "", fmt.Errorf("FAILED TO GET SERVER GAME: %v", err)
	}

	return jeu, nil
}
