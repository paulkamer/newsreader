package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"newsreader/models"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

func InsertArticle(dbconn *sql.DB, article models.Article) error {
	query := `
		INSERT INTO articles (id, source_id, title, url, body) 
		            VALUES ($1, $2, $3, $4, $5)
		`

	_, err := dbconn.Exec(query, article.ID, article.Source, article.Title, article.Url, article.Body)

	if err != nil {
		return fmt.Errorf("failed to insert article: %s", err)
	}

	return err
}

func FetchArticle(dbconn *sql.DB, id uuid.UUID) (models.Article, error) {
	query := `SELECT id, source_id, title, url, body, created_at FROM articles WHERE id = ?`

	rows := dbconn.QueryRow(query, id)

	var article models.Article
	err := rows.Scan(&article.ID, &article.Source, &article.Title, &article.Url, &article.Body, &article.CreatedAt)
	if err != nil {
		return article, err
	}

	return article, nil
}

func UpdateArticle(dbconn *sql.DB, article models.Article) error {
	query := `
		UPDATE articles SET
			title = ?,
			url = ?,
			body = ?,
			update_at = CURRENT_TIMESTAMP
		WHERE id = ?
		`
	result, err := dbconn.Exec(query, article.Title, article.Url, article.Body, article.ID)
	if err != nil {
		return err
	}

	// Check the number of affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		log.Errorf("Failed to update article: %v", err)
		return errors.New("failed to update article")
	}
	log.Debugf("Successfully updated %d row(s)\n", rowsAffected)

	return err
}

func ListArticles(dbconn *sql.DB, source_guid uuid.UUID) ([]models.Article, error) {
	query := `SELECT id, source_id, title, url, body, created_at FROM articles WHERE source_id = ?`

	rows, err := dbconn.Query(query, source_guid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.Source, &article.Title, &article.Url, &article.Body, &article.CreatedAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}
