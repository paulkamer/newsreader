package feedparsers

import (
	"fmt"
	"log"
	"newsreader/feedtypes"
)

func ParseAtomFeed(body []byte) (*feedtypes.AtomFeed, error) {
	log.Fatal("ParseAtomFeed not implemented")
	return nil, nil
}

func debugParsedAtomFeed(feed *feedtypes.AtomFeed) {
	// Print out the feed information
	fmt.Printf("Feed Title: %s\n", feed.Title)
	fmt.Printf("Feed Updated: %s\n", feed.Updated)
	fmt.Printf("Feed ID: %s\n", feed.ID)

	// Iterate over the entries and print them out
	for _, entry := range feed.Entries {
		fmt.Printf("Entry Title: %s\n", entry.Title)
		fmt.Printf("Entry Updated: %s\n", entry.Updated)
		fmt.Printf("Entry ID: %s\n", entry.ID)
		fmt.Printf("Entry Summary: %s\n", entry.Summary)
		fmt.Printf("Entry Content: %s\n", entry.Content)
	}
}
