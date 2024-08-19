package feedtypes

import "encoding/xml"

// AtomFeed represents the top-level Atom feed element
type AtomFeed struct {
	XMLName xml.Name    `xml:"feed"`
	Title   string      `xml:"title"`
	Link    []AtomLink  `xml:"link"`
	Updated string      `xml:"updated"`
	Author  AtomAuthor  `xml:"author"`
	ID      string      `xml:"id"`
	Entries []AtomEntry `xml:"entry"`
}

// AtomLink represents the link element in the feed and entry
type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
}

// AtomAuthor represents the author element in the feed and entry
type AtomAuthor struct {
	Name  string `xml:"name"`
	URI   string `xml:"uri,omitempty"`
	Email string `xml:"email,omitempty"`
}

// AtomEntry represents an individual entry in the Atom feed
type AtomEntry struct {
	Title   string      `xml:"title"`
	Link    []AtomLink  `xml:"link"`
	ID      string      `xml:"id"`
	Updated string      `xml:"updated"`
	Summary string      `xml:"summary,omitempty"`
	Content string      `xml:"content,omitempty"`
	Author  *AtomAuthor `xml:"author,omitempty"`
}
