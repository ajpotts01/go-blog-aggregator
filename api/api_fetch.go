package api

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ajpotts01/go-blog-aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

func createPostParams(post rssItem, feedId uuid.UUID) (database.CreatePostParams, error) {
	newId, err := uuid.NewUUID()

	if err != nil {
		return database.CreatePostParams{}, err
	}

	createdAt := time.Now()

	// Times of all sampled feeds appear to conform to RFC1123Z standard
	// this is despite the RSS spec saying RFC822...
	// ...supported by examples that aren't RFC822. Nice one.
	publishedAt, err := time.Parse(time.RFC1123Z, post.PubDate)

	if err != nil {
		return database.CreatePostParams{}, err
	}

	params := database.CreatePostParams{
		ID:        newId,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		Title:     post.Title,
		PublishedAt: sql.NullTime{
			Time: publishedAt,
		},
		Description: sql.NullString{
			String: post.Description,
		},
		FeedID: feedId,
		Url:    post.Link,
	}

	return params, nil
}

func (config *ApiConfig) fetchFeed(url string) (*rss, error) {
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

func (config *ApiConfig) processFeed(feed *rss, feedId uuid.UUID) error {
	for _, c := range feed.Channels {
		for _, item := range c.Items {
			params, err := createPostParams(item, feedId)

			if err != nil {
				return err
			}

			post, err := config.DbConn.CreatePost(context.TODO(), params)
			if err != nil {
				if postgresErr, ok := err.(*pq.Error); ok {
					if postgresErr.Code == "23505" {
						log.Printf("Post already exists: %v", item.Title)
					} else {
						return err
					}
				}
			} else {
				log.Printf("Created post: %v", post.Title)
			}
		}

		log.Printf("Processed feed: %v", c.Title)
	}

	return nil
}

func (config *ApiConfig) FetchLoop() {
	loopTimer := 60 * time.Second
	ticker := time.NewTicker(loopTimer)

	log.Printf("Init fetch loop")

	for {
		var urlPool sync.WaitGroup
		// Block until a signal is received from the ticker
		<-ticker.C

		log.Printf("Fetch loop running...")
		// Want to grab up to X feeds at once.
		// This is configured in config.MaxNumProcessed
		feeds, err := config.DbConn.GetNextFeedsToFetch(context.TODO(), int32(config.MaxFeedsProcessed))

		if err != nil {
			log.Printf("Error: failed to retrieve feeds to fetch: %v", err)
			continue
		}

		log.Printf("Got %v feeds", len(feeds))

		for _, feed := range feeds {
			urlPool.Add(1)
			go func(url string, id uuid.UUID) {
				defer urlPool.Done()
				log.Printf("Fetching from %s", url)
				rss, err := config.fetchFeed(url)
				config.DbConn.MarkFeedFetched(context.TODO(), id)
				if err != nil {
					log.Printf("Error: failed to retrieve items from feed %s: %v", url, err)
					return
				}
				config.processFeed(rss, id)
			}(feed.Url, feed.ID)
		}
		log.Printf("Waiting for fetching to end...")
		urlPool.Wait()
	}
}
