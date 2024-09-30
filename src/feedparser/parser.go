package feedparser

import (
	"encoding/xml"
	"newsreader/feedtypes"
)

func ParseFeed[T feedtypes.RssFeed | feedtypes.AtomFeed](body []byte) (*T, error) {
	var feed T
	err := xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}
