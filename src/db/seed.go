package db

import (
	"log"
	"time"

	"github.com/google/uuid"
)

func Seed() {
	hasrecords, _ := HasNewssources()

	if hasrecords {
		return
	}

	newssources := []Newssource{
		{uuid.New(), "CNN", "https://cnn.com", URGENT, true, time.Now(), time.Now()},
	}

	for _, newssource := range newssources {
		err := InsertNewssource(newssource)
		if err != nil {
			log.Printf("Failed to insert newssource %s: %v", newssource.Title, err)
		}
	}
}
