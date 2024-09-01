package feedparsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRssFeed(t *testing.T) {
	feedData := `<?xml version="1.0" encoding="UTF-8" ?>
		<rss version="2.0">
		<channel>
			<title>Example Feed</title>
			<link>http://www.example.com/</link>
			<description>This is an example feed</description>
			<item>
				<title>Test entry</title>
				<link>http://www.example.com/test-entry</link>
				<description>This is a test entry</description>
			</item>
		</channel>
		</rss>`

	feed, err := ParseRssFeed([]byte(feedData))

	assert.Nil(t, err)
	assert.Equal(t, "Example Feed", feed.Channel.Title)
	assert.Equal(t, "Test entry", feed.Channel.Items[0].Title)
}

func TestParseRssFeed_Fail(t *testing.T) {
	feedData := `<?xml version="1.0" encoding="UTF-8" ?>
		<rss version="2.0">
		<channel>
			<title>Example Feed</title>
			<link>http://www.example.com/</link>
			<description>This is an example feed</description>
			<item_________________________INVALID_________>
				<title>Test entry</title>
				<link>http://www.example.com/test-entry</link>
				<description>This is a test entry</description>
			</item>
		</channel>
		</rss>`

	_, err := ParseRssFeed([]byte(feedData))

	assert.Error(t, err)
}
