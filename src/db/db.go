package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// TODO dont use global var :x
var DB *sql.DB

func InitDatabase(dbType, dataSourceName string) {
	var err error
	DB, err = OpenDatabase(dbType, dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to %s database: %v", dbType, err)
	}
	fmt.Printf("Connected to %s database!\n", dbType)

	CreateDbStructure()
	Seed()
}

func OpenDatabase(dbType, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(dbType, dataSourceName)
	if err != nil {
		return nil, err
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateDbStructure() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS newssources (
        id GUID PRIMARY KEY,
        title TEXT NOT NULL,
        url TEXT NOT NULL,
		update_priority TEXT CHECK( update_priority IN ('URGENT','HIGH','MED', 'LOW') ) NOT NULL DEFAULT 'MED',
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME
    );`

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	fmt.Printf("Created DB structure")
}

func InsertNewssource(newssource Newssource) error {
	query := `
		INSERT INTO newssources (id, title, url, update_priority, is_active) 
		            VALUES ($1, $2, $3, $4, $5)
		`

	_, err := DB.Exec(query, newssource.ID, newssource.Title, newssource.Url, newssource.UpdatePriority, newssource.IsActive)

	if err != nil {
		return fmt.Errorf("failed to insert newssource: %s", err)
	}

	return err
}

func FetchNewssource(guid uuid.UUID) (Newssource, error) {
	query := `SELECT * FROM newssources WHERE id = $1`

	rows := DB.QueryRow(query, guid)

	var newssource Newssource
	err := rows.Scan(&newssource.ID, &newssource.Title, &newssource.Url, &newssource.UpdatePriority, &newssource.IsActive, &newssource.CreatedAt, &newssource.UpdatedAt)
	if err != nil {
		return newssource, err
	}

	return newssource, nil
}

func UpdateNewssource(newssource Newssource) error {
	log.Printf("UpdateNewssource %v", newssource)

	query := `
		UPDATE newssources SET
			title = ?,
			url = ?,
			update_priority = ?,
			is_active = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		`

	result, err := DB.Exec(query, newssource.Title, newssource.Url, newssource.UpdatePriority, newssource.IsActive, newssource.ID)

	if err != nil {
		return fmt.Errorf("failed to update newssource: %s", err)
	}

	// Check the number of affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal("Failed to retrieve rows affected:", err)
	}
	fmt.Printf("Successfully updated %d row(s)\n", rowsAffected)

	return err
}

func DeleteNewssource(guid uuid.UUID) error {
	query := `DELETE FROM newssources WHERE id = $1`

	_, err := DB.Exec(query, guid)
	if err != nil {
		log.Printf("failed to delete newssource: %s", err)
	}

	return err
}

func HasNewssources() (bool, error) {
	query := `SELECT id FROM newssources LIMIT 1`

	rows := DB.QueryRow(query)

	var newssource Newssource
	err := rows.Scan(&newssource.ID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ListNewssources() ([]Newssource, error) {
	query := `SELECT id, title, url, created_at FROM newssources`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newssources []Newssource
	for rows.Next() {
		var newssource Newssource
		err := rows.Scan(&newssource.ID, &newssource.Title, &newssource.Url, &newssource.CreatedAt)
		if err != nil {
			return nil, err
		}
		newssources = append(newssources, newssource)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return newssources, nil
}
