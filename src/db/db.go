package db

import (
	"database/sql"

	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase(dbType, dataSourceName string) (*sql.DB, error) {
	var err error
	dbConn, err := Connect(dbType, dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to %s database: %v", dbType, err)
	}

	err = createDbStructure(dbConn)
	if err != nil {
		log.Fatalf("Failed to create database structure: %v", err)
	}

	Seed(dbConn)

	return dbConn, nil
}

func Connect(dbType, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(dbType, dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func createDbStructure(dbConn *sql.DB) error {
	createNewssourceTableSQL := `CREATE TABLE IF NOT EXISTS newssources (
        id GUID PRIMARY KEY,
        title TEXT NOT NULL,
        url TEXT NOT NULL,
		update_priority TEXT CHECK( update_priority IN ('URGENT','HIGH','MED', 'LOW') ) NOT NULL DEFAULT 'MED',
		feed_type TEXT CHECK( feed_type IN ('rss','atom') ) NOT NULL,
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME
    );`

	_, err := dbConn.Exec(createNewssourceTableSQL)
	if err != nil {
		return err
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

	_, err = dbConn.Exec(createArticleTableQuery)

	return err
}
