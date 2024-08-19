package feedparsers

import (
	"encoding/xml"
	"newsreader/feedtypes"
)

func ParseRssFeed(body []byte) (*feedtypes.RSS, error) {
	var rss feedtypes.RSS
	err := xml.Unmarshal(body, &rss)
	if err != nil {
		return nil, err
	}

	return &rss, nil
}
