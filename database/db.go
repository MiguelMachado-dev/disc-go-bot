package database

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/MiguelMachado-dev/disc-go-bot/utils"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Init initializes the database connection
func Init() error {
	// Ensure data directory exists
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		return err
	}

	// Open SQLite database
	var err error
	db, err = sql.Open("sqlite3", filepath.Join("data", "disc-go-bot.db")+"?_busy_timeout=5000&_journal=WAL&_sync=NORMAL")
	if err != nil {
		return err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Create tables if not exist
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS gemini_api_keys (
		user_id TEXT PRIMARY KEY,
		api_key TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// StoreGeminiAPIKey saves or updates a user's Gemini API key
func StoreGeminiAPIKey(userID, apiKey string) error {
	encryptedKey, err := utils.Encrypt(apiKey)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	INSERT INTO gemini_api_keys (user_id, api_key, updated_at)
	VALUES (?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(user_id)
	DO UPDATE SET api_key = ?, updated_at = CURRENT_TIMESTAMP
	`, userID, encryptedKey, encryptedKey)
	return err
}

// GetGeminiAPIKey retrieves a user's Gemini API key
func GetGeminiAPIKey(userID string) (string, error) {
	var encryptedKey string
	err := db.QueryRow("SELECT api_key FROM gemini_api_keys WHERE user_id = ?", userID).Scan(&encryptedKey)
	if err != nil {
		return "", err
	}

	return utils.Decrypt(encryptedKey)
}
