package feedparsers

import (
	"log"
	"newsreader/feedtypes"
)

func ParseAtomFeed(body []byte) (*feedtypes.AtomFeed, error) {
	log.Fatal("ParseAtomFeed not implemented")
	return nil, nil
}

func debugParsedAtomFeed(feed *feedtypes.AtomFeed) {
	// Print out the feed information
	log.Printf("Feed Title: %s\n", feed.Title)
	log.Printf("Feed Updated: %s\n", feed.Updated)
	log.Printf("Feed ID: %s\n", feed.ID)

	// Iterate over the entries and print them out
	for _, entry := range feed.Entries {
		log.Printf("Entry Title: %s\n", entry.Title)
		log.Printf("Entry Updated: %s\n", entry.Updated)
		log.Printf("Entry ID: %s\n", entry.ID)
		log.Printf("Entry Summary: %s\n", entry.Summary)
		log.Printf("Entry Content: %s\n", entry.Content)
	}
}
