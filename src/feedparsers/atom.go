package feedparsers

import (
	"encoding/xml"
	"newsreader/feedtypes"
)

func ParseAtomFeed(body []byte) (*feedtypes.AtomFeed, error) {
	var feed feedtypes.AtomFeed
	err := xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}
