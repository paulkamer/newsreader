package feedparsers

import (
	"newsreader/feedtypes"

	log "github.com/sirupsen/logrus"
)

func ParseAtomFeed(body []byte) (*feedtypes.AtomFeed, error) {
	log.Fatal("ParseAtomFeed not implemented")
	return nil, nil
}

func debugParsedAtomFeed(feed *feedtypes.AtomFeed) {
	// Print out the feed information
	log.Debugf("Feed Title: %s\n", feed.Title)
	log.Debugf("Feed Updated: %s\n", feed.Updated)
	log.Debugf("Feed ID: %s\n", feed.ID)

	// Iterate over the entries and print them out
	for _, entry := range feed.Entries {
		log.Debugf("Entry Title: %s\n", entry.Title)
		log.Debugf("Entry Updated: %s\n", entry.Updated)
		log.Debugf("Entry ID: %s\n", entry.ID)
		log.Debugf("Entry Summary: %s\n", entry.Summary)
		log.Debugf("Entry Content: %s\n", entry.Content)
	}
}
