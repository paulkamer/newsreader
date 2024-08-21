package jobs

import (
	"database/sql"
	"errors"
	"io"
	"net/http"
	"newsreader/db"
	"newsreader/feedparsers"
	"newsreader/models"
	"newsreader/repositories"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

func FetchNews(newssource_guid uuid.UUID) ([]models.Article, error) {
	log.Debugf("Fetching news for source: %s\n", newssource_guid)

	dbconn, _ := db.Connect(db.SQLiteType, db.SQLiteDataSource)
	defer dbconn.Close()

	newssource, err := repositories.FetchNewssource(dbconn, newssource_guid)
	if err != nil {
		log.Fatalf("Failed to fetch newssource %s: %v", newssource_guid, err)
	}

	body, err := FetchFeed(newssource.Url)
	if err != nil {
		log.Errorf("Failed to fetch feed: %v\n", err)
		return nil, err
	}

	articles, err := ParseFeed(body, newssource)
	if err != nil {
		log.Errorf("Failed to parse feed: %v\n", err)
		return nil, err
	}

	storeArticles(dbconn, articles)

	return articles, nil
}

func FetchFeed(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Error making HTTP request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body: %v\n", err)
		return nil, err
	}

	return body, nil
}

func ParseFeed(body []byte, newssource models.Newssource) ([]models.Article, error) {
	var articles []models.Article

	// TODO determine feed type & update newssource if necessary

	switch newssource.FeedType {
	case models.RSS:
		log.Print("Parsing RSS feed\n")
		rss, err := feedparsers.ParseRssFeed(body)

		if err != nil {
			log.Errorf("Error parsing RSS feed: %v\n", err)
			return nil, err
		}

		for _, item := range rss.Channel.Items {
			CreatedAt, _ := time.Parse(time.RFC1123Z, item.PubDate)

			article := models.Article{
				ID:        item.Guid,
				Source:    newssource.ID,
				Title:     item.Title,
				Url:       item.Link,
				CreatedAt: CreatedAt,
			}
			articles = append(articles, article)
		}
	case models.ATOM:
		log.Print("Parsing Atom feed\n")
		atom, _ := feedparsers.ParseAtomFeed(body)
		log.Debugf("Atom: %v\n", atom)
	default:
		log.Debugf("Unknown feed type %s", newssource.FeedType)
		return nil, errors.New("unknown feed type")
	}

	return articles, nil
}

func storeArticles(dbconn *sql.DB, articles []models.Article) error {

	for _, article := range articles {
		err := repositories.InsertArticle(dbconn, article)

		if err != nil {
			log.Errorf("Failed to insert article: %v\n", err)
		}
	}

	return nil
}
