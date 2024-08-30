package feedparsers

import (
	"encoding/xml"
	"newsreader/feedtypes"

	log "github.com/sirupsen/logrus"
)

func ParseRssFeed(body []byte) (*feedtypes.RSS, error) {
	var rss feedtypes.RSS
	err := xml.Unmarshal(body, &rss)
	if err != nil {
		return nil, err
	}

	if log.GetLevel() == log.DebugLevel {
		log.Debugf("Parsed RSS feed: %v", rss)
		debugParsedRssFeed(&rss)
	}

	return &rss, nil
}

func debugParsedRssFeed(feed *feedtypes.RSS) {
	// Print out the channel information
	log.Debugf("Channel Title: %s\n", feed.Channel.Title)
	log.Debugf("Channel Link: %s\n", feed.Channel.Link)
	log.Debugf("Channel Description: %s\n", feed.Channel.Description)

	// Iterate over the items and print them out
	for _, item := range feed.Channel.Items {
		log.Debugf("Item Title: %s\n", item.Title)
		log.Debugf("Item Link: %s\n", item.Link)
		log.Debugf("Item Description: %s\n", item.Description)
		log.Debugf("Item PubDate: %s\n", item.PubDate)
		log.Debugf("Item GUID: %s\n\n", item.Guid)
	}
}
