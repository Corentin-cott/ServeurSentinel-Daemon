package db

import (
	"database/sql"
	"fmt"

	"github.com/Corentin-cott/ServeurSentinel/config"
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
func CheckAndInsertPlayer(playerName string) (int, error) {
	// Pour récupéré l'UUID a partir du nom du joueur, ont utilise une api
	playerUUID := playerName

	// Vérifier si le joueur existe déjà
	var userID int
	query := "SELECT utilisateur_id FROM joueurs_mc WHERE uuid = ?"
	err := db.QueryRow(query, playerUUID).Scan(&userID)

	// Si le joueur n'existe pas, l'ajouter
	if err == sql.ErrNoRows {
		// Ajouter un nouveau joueur dans la table joueurs_mc avec serveur_id = 1 [DEBUG]
		insertQuery := "INSERT INTO joueurs_mc (uuid, utilisateur_id, premiere_co, derniere_co) VALUES (?, ?, NOW(), NOW())"
		result, err := db.Exec(insertQuery, playerUUID, 1) // 1 est l'ID de l'utilisateur par défaut, ici à ajuster selon les besoins.
		if err != nil {
			return 0, fmt.Errorf("failed to insert player: %v", err)
		}

		// Récupérer l'ID de l'utilisateur
		userID64, err := result.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("failed to get last insert id: %v", err)
		}

		userID := int(userID64) // Convertir int64 en int car en go Int64 n'est pas le même type que Int

		fmt.Printf("Player with UUID %s was added with ID %d\n", playerUUID, userID)
	} else if err != nil {
		return 0, fmt.Errorf("failed to query player: %v", err)
	} else {
		fmt.Printf("Player with UUID %s found with ID %d\n", playerUUID, userID)
	}

	return userID, nil // Retourner l'ID de l'utilisateur
}

// SaveConnectionLog sauvegarde un log de connexion dans la base de données
func SaveConnectionLog(playerUUID string, serverID int) error {
	userID, err := CheckAndInsertPlayer(playerUUID)
	if err != nil {
		return fmt.Errorf("failed to check or insert player: %v", err)
	}

	// Enregistrer le log de connexion
	query := `
		INSERT INTO serveurs_connections_log (joueurs_mc_uuid, serveur_id, date)
		VALUES (?, ?, NOW())
	`

	_, err = db.Exec(query, userID, serverID)
	if err != nil {
		return fmt.Errorf("failed to insert connection log: %v", err)
	}

	fmt.Println("Connection log successfully saved.")
	return nil
}
