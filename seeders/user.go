package seeders

import (
	"database/sql"
	"fmt"

	"github.com/edwinjordan/wmsTest_Golang/utils"
	"github.com/google/uuid"
)

// SeedUsers populates the users table with sample data
func SeedUsers(db *sql.DB) error {
	// Hash passwords
	password, err := utils.HashPassword("password123")
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Define user data with API keys
	users := []struct {
		username string
		email    string
		apiKey   string
	}{
		{"admin", "admin@wms.local", "wms-admin-" + uuid.New().String()[:8]},
		{"manager", "manager@wms.local", "wms-manager-" + uuid.New().String()[:8]},
		{"operator", "operator@wms.local", "wms-operator-" + uuid.New().String()[:8]},
		{"viewer", "viewer@wms.local", "wms-viewer-" + uuid.New().String()[:8]},
		{"alice", "alice@example.com", "wms-alice-" + uuid.New().String()[:8]},
		{"bob", "bob@example.com", "wms-bob-" + uuid.New().String()[:8]},
	}

	// Insert users with password hashes and API keys
	for _, user := range users {
		result, err := db.Exec(`
			INSERT INTO users (username, email, password, api_key, is_active) 
			VALUES ($1, $2, $3, $4, true)
			ON CONFLICT (email) DO NOTHING;
		`, user.username, user.email, password, user.apiKey)

		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", user.username, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected for user %s: %w", user.username, err)
		}

		if rowsAffected > 0 {
			fmt.Printf("Inserted user: %s (API Key: %s)\n", user.username, user.apiKey)
		} else {
			fmt.Printf("User %s already exists, skipping...\n", user.username)
		}
	}

	fmt.Println("User seeding completed!")
	return nil
}

// ClearUsers removes all users from the users table (useful for testing)
func ClearUsers(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		return fmt.Errorf("failed to clear users: %w", err)
	}

	fmt.Println("Cleared all users from the database")
	return nil
}
