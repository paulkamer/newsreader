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

	storeArticles(dbconn, &articles)

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
			createdAt, _ := time.Parse(time.RFC1123, item.PubDate)

			article := models.Article{
				ID:        uuid.New(),
				Source:    newssource.ID,
				Title:     item.Title,
				Url:       item.Link,
				Body:      item.Description,
				CreatedAt: createdAt,
			}
			articles = append(articles, article)
		}
	case models.ATOM:
		log.Print("Parsing Atom feed\n")
		atom, err := feedparsers.ParseAtomFeed(body)
		log.Debugf("Atom: %v\n", atom)

		if err != nil {
			log.Errorf("Error parsing Atom feed: %v\n", err)
			return nil, err
		}

		for _, entry := range atom.Entries {
			createdAt, _ := time.Parse(time.RFC3339, entry.Updated)

			article := models.Article{
				ID:        uuid.New(),
				Source:    newssource.ID,
				Title:     entry.Title,
				Body:      entry.Summary,
				CreatedAt: createdAt,
			}

			articles = append(articles, article)
		}

	default:
		log.Debugf("Unknown feed type %s", newssource.FeedType)
		return nil, errors.New("unknown feed type")
	}

	return articles, nil
}

func storeArticles(dbconn *sql.DB, articles *[]models.Article) error {
	for _, article := range *articles {
		err := repositories.InsertArticle(dbconn, &article)

		if err != nil {
			log.Errorf("Failed to insert article: %v\n", err)
		}
	}

	if len(*articles) == 0 {
		log.Debug("No articles to store\n")
		return nil
	}

	newssource, err := repositories.FetchNewssource(dbconn, (*articles)[0].Source)
	if err != nil {
		log.Errorf("Failed to fetch newssource: %v\n", err)
		return err
	}
	repositories.UpdateNewssource(dbconn, &newssource)

	return nil
}
