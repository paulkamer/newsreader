package feedparsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAtomFeed(t *testing.T) {
	feedData := `<?xml version="1.0" encoding="utf-8"?>
	<feed xmlns="http://www.w3.org/2005/Atom">
		<title>Example Feed</title>
		<subtitle>A subtitle.</subtitle>
		<link href="http://example.org/feed/" rel="self" />
		<link href="http://example.org/" />
		<id>urn:uuid:60a76c80-d399-11d9-b91C-0003939e0af6</id>
		<updated>2003-12-13T18:30:02Z</updated>
		<entry>
			<title>Test entry</title>
			<link href="http://example.org/2003/12/13/atom03" />
			<id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
			<updated>2003-12-13T18:30:02Z</updated>
			<summary>Some text</summary>
		</entry>
	</feed>`

	feed, err := ParseAtomFeed([]byte(feedData))

	assert.Nil(t, err)
	assert.Equal(t, "Example Feed", feed.Title)
	assert.Equal(t, "Test entry", feed.Entries[0].Title)
}

func TestParseAtomFeed_Fail(t *testing.T) {
	feedData := `<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom">
			<title>Example Feed</title>
			<entry_________________________INVALID_________>
				<title>Test entry</title>
			</entry>
		</feed>`

	_, err := ParseAtomFeed([]byte(feedData))

	assert.Error(t, err)
}
