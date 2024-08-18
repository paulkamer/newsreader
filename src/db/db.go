package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase(dbType, dataSourceName string) (*sql.DB, error) {
	var err error
	dbConn, err := OpenDatabase(dbType, dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to %s database: %v", dbType, err)
	}
	fmt.Printf("Connected to %s database!\n", dbType)

	CreateDbStructure(dbConn)
	Seed(dbConn)

	return dbConn, nil
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

func CreateDbStructure(dbConn *sql.DB) {
	createNewssourceTableSQL := `CREATE TABLE IF NOT EXISTS newssources (
        id GUID PRIMARY KEY,
        title TEXT NOT NULL,
        url TEXT NOT NULL,
		update_priority TEXT CHECK( update_priority IN ('URGENT','HIGH','MED', 'LOW') ) NOT NULL DEFAULT 'MED',
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME
    );`

	_, err := dbConn.Exec(createNewssourceTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	createArticleTableQuery := `CREATE TABLE IF NOT EXISTS articles (
		id GUID PRIMARY KEY,
		source_id GUID NOT NULL,
		title TEXT NOT NULL,
		url TEXT NOT NULL,
		body TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME
	);`

	res, err := dbConn.Exec(createArticleTableQuery)
	if err != nil || res == nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	fmt.Printf("Created DB structure")
}
