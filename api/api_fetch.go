package api

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
)

// RSS feeds are always XML spec

/*
	You can test with these ones:

	https://blog.boot.dev/index.xml
	https://wagslane.dev/index.xml
	And any other blogs you enjoy that have RSS feeds.
*/

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Guid        string `xml:"guid"`
	Description string `xml:"description"`
}

type rssChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Generator     string    `xml:"generator"`
	Language      string    `xml:"language"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []rssItem `xml:"item"`
	//AtomLink      string    `xml:"xmlns atom,attr"`
	// atom:link?
}

type rss struct {
	Channels []rssChannel `xml:"channel"`
}

func (config *ApiConfig) FetchFeed(url string) (*rss, error) {
	var rawData []byte
	var feed *rss

	log.Printf("Reading from %v", url)
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("Error getting feed: %v", err)
		return feed, err
	}

	//log.Printf("Raw response: %v", resp.Body)
	log.Printf("Response code: %v", resp.StatusCode)
	rawData, err = io.ReadAll(resp.Body)
	log.Printf("Bytes read: %v", len(rawData))
	defer resp.Body.Close()

	if err != nil {
		log.Printf("Error reading feed response: %v", err)
		return feed, err
	}

	err = xml.Unmarshal(rawData, &feed)
	if err != nil {
		log.Printf("Error unmarshalling XML: %v", err)
		return feed, err
	}

	// test := string(rawData[:])
	// log.Printf(test)
	return feed, nil
}
