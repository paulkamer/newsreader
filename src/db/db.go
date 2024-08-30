package db

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase(dbType, dataSourceName string, migrationsDir string) (*sql.DB, error) {
	var err error
	dbConn, err := Connect(dbType, dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to %s database: %v", dbType, err)
	}

	runMigrations(dbConn, migrationsDir)

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

func runMigrations(dbConn *sql.DB, migrationsDir string) {
	driver, err := sqlite3.WithInstance(dbConn, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Failed to set up DB migrations: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsDir, "sqlite3", driver)
	if err != nil {
		log.Fatalf("Failed to start DB migrations: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run DB migrations: %v", err)
	}
}
