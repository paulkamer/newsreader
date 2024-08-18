package db

import (
	"database/sql"
	"log"
	"newsreader/models"
	"newsreader/repositories"
	"time"

	"github.com/google/uuid"
)

func Seed(dbconn *sql.DB) {
	hasrecords, _ := repositories.HasNewssources(dbconn)

	if hasrecords {
		return
	}

	newssources := []models.Newssource{
		{uuid.New(), "CNN", "https://cnn.com", models.URGENT, true, time.Now(), time.Now()},
	}

	for _, newssource := range newssources {
		err := repositories.InsertNewssource(dbconn, newssource)
		if err != nil {
			log.Printf("Failed to insert newssource %s: %v", newssource.Title, err)
		}
	}
}
