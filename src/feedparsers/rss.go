package feedparsers

import (
	"encoding/xml"
	"log"
	"newsreader/feedtypes"
)

func ParseRssFeed(body []byte) (*feedtypes.RSS, error) {
	var rss feedtypes.RSS
	err := xml.Unmarshal(body, &rss)
	if err != nil {
		return nil, err
	}

	log.Printf("Parsed RSS feed: %v", rss)

	debugParsedRssFeed(&rss)

	return &rss, nil
}

func debugParsedRssFeed(feed *feedtypes.RSS) {
	// Print out the channel information
	log.Printf("Channel Title: %s\n", feed.Channel.Title)
	log.Printf("Channel Link: %s\n", feed.Channel.Link)
	log.Printf("Channel Description: %s\n", feed.Channel.Description)

	// Iterate over the items and print them out
	for _, item := range feed.Channel.Items {
		log.Printf("Item Title: %s\n", item.Title)
		log.Printf("Item Link: %s\n", item.Link)
		log.Printf("Item Description: %s\n", item.Description)
		log.Printf("Item PubDate: %s\n", item.PubDate)
		log.Printf("Item GUID: %s\n\n", item.Guid)
	}
}
