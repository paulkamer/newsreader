package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"newsreader/models"

	"github.com/google/uuid"
)

func InsertNewssource(dbconn *sql.DB, newssource models.Newssource) error {
	query := `
		INSERT INTO newssources (id, title, url, update_priority, is_active) 
		            VALUES ($1, $2, $3, $4, $5)
		`

	_, err := dbconn.Exec(query, newssource.ID, newssource.Title, newssource.Url, newssource.UpdatePriority, newssource.IsActive)

	if err != nil {
		return fmt.Errorf("failed to insert newssource: %s", err)
	}

	return err
}

func FetchNewssource(dbconn *sql.DB, guid uuid.UUID) (models.Newssource, error) {
	query := `SELECT * FROM newssources WHERE id = ?`

	rows := dbconn.QueryRow(query, guid)

	var newssource models.Newssource
	err := rows.Scan(&newssource.ID, &newssource.Title, &newssource.Url, &newssource.UpdatePriority, &newssource.IsActive, &newssource.CreatedAt, &newssource.UpdatedAt)
	if err != nil {
		return newssource, err
	}

	return newssource, nil
}

func UpdateNewssource(dbconn *sql.DB, newssource models.Newssource) error {
	query := `
		UPDATE newssources SET
			title = ?,
			url = ?,
			update_priority = ?,
			is_active = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		`

	result, err := dbconn.Exec(query, newssource.Title, newssource.Url, newssource.UpdatePriority, newssource.IsActive, newssource.ID)

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

func DeleteNewssource(dbconn *sql.DB, guid uuid.UUID) error {
	query := `DELETE FROM newssources WHERE id = $1`

	_, err := dbconn.Exec(query, guid)
	if err != nil {
		log.Printf("failed to delete newssource: %s", err)
	}

	return err
}

func HasNewssources(dbconn *sql.DB) (bool, error) {
	query := `SELECT id FROM newssources LIMIT 1`

	rows := dbconn.QueryRow(query)

	var newssource models.Newssource
	err := rows.Scan(&newssource.ID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ListNewssources(dbconn *sql.DB) ([]models.Newssource, error) {
	query := `SELECT id, title, url, created_at FROM newssources`

	rows, err := dbconn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newssources []models.Newssource
	for rows.Next() {
		var newssource models.Newssource
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
