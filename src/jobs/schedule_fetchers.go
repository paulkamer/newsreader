package jobs

import (
	"newsreader/db"
	"newsreader/repositories"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func DetermineOutdatedNewssources(listChan chan uuid.UUID) {
	dbconn, err := db.Connect(db.SQLiteType, db.SQLiteDataSource)
	if err != nil {
		log.Error(err)
		return
	}

	// TODO Determine which feeds are outdated instead of all feeds
	newssources, err := repositories.ListNewssources(dbconn)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("Outdated newssources: %d", len(newssources))

	for _, newssource := range newssources {
		log.Debugf("Outdated newssource: %s", newssource.ID)

		listChan <- newssource.ID
	}
}
