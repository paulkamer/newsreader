package db

import (
	"database/sql"
	"fmt"
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
		{
			ID:             uuid.New(),
			Title:          "NASA",
			Url:            "https://earthobservatory.nasa.gov/feeds/earth-observatory.rss",
			FeedType:       models.RSS,
			UpdatePriority: models.URGENT,
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	for _, newssource := range newssources {
		err := repositories.InsertNewssource(dbconn, newssource)
		if err != nil {
			_ = fmt.Errorf("failed to insert newssource %s: %v", newssource.Title, err)
			return
		}
	}
}
