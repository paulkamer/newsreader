package feedparsers

import (
	"encoding/xml"
	"fmt"
	"newsreader/feedtypes"
)

func ParseRssFeed(body []byte) (*feedtypes.RSS, error) {
	var rss feedtypes.RSS
	err := xml.Unmarshal(body, &rss)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Parsed RSS feed: %v", rss)

	debugParsedRssFeed(&rss)

	return &rss, nil
}

func debugParsedRssFeed(feed *feedtypes.RSS) {
	// Print out the channel information
	fmt.Printf("Channel Title: %s\n", feed.Channel.Title)
	fmt.Printf("Channel Link: %s\n", feed.Channel.Link)
	fmt.Printf("Channel Description: %s\n", feed.Channel.Description)

	// Iterate over the items and print them out
	for _, item := range feed.Channel.Items {
		fmt.Printf("Item Title: %s\n", item.Title)
		fmt.Printf("Item Link: %s\n", item.Link)
		fmt.Printf("Item Description: %s\n", item.Description)
		fmt.Printf("Item PubDate: %s\n", item.PubDate)
		fmt.Printf("Item GUID: %s\n\n", item.Guid)
	}
}
