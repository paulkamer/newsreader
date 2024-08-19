package feedtypes

import "encoding/xml"

// RSS represents the entire RSS feed
type RSS struct {
	XMLName xml.Name   `xml:"rss"`
	Channel RssChannel `xml:"channel"`
}

// RssChannel represents the channel information of the RSS feed
type RssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []RssItem `xml:"item"`
}

// RssItem represents each individual item (e.g., a blog post) in the RSS feed
type RssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Guid        string `xml:"guid"`
}
