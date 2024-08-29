package repositories

import (
	"fmt"
	"newsreader/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInsertArticle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	article := &models.Article{
		ID:        uuid.New(),
		Source:    uuid.New(),
		Title:     "title",
		Url:       "url",
		Body:      "body",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectExec("INSERT INTO articles \\(id, source_id, title, url, body\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\)").WillReturnResult(sqlmock.NewResult(1, 1))

	err = InsertArticle(db, article)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.ExpectExec("INSERT INTO newssources").WillReturnError(nil)
	err = InsertArticle(db, article)

	assert.Error(t, err)
}

func TestFetchArticle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := uuid.New()

	article := models.Article{
		ID:        id,
		Source:    uuid.New(),
		Title:     "title",
		Url:       "url",
		Body:      "body",
		CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "source_id", "title", "url", "body", "created_at"}).
		AddRow(article.ID, article.Source, article.Title, article.Url, article.Body, article.CreatedAt)

	mock.ExpectQuery("SELECT id, source_id, title, url, body, created_at FROM articles WHERE id = \\?").WillReturnRows(rows)

	_, err = FetchArticle(db, article.ID)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	rows = sqlmock.NewRows([]string{"id", "source_id", "title", "url", "body", "created_at"}).
		RowError(1, fmt.Errorf("row error"))

	mock.ExpectQuery("SELECT id, source_id, title, url, body, created_at FROM articles WHERE id = \\?").WillReturnRows(rows)

	_, err = FetchArticle(db, article.ID)
	assert.Error(t, err)
}

func TestUpdateArticle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := uuid.New()

	article := &models.Article{
		ID:        id,
		Source:    uuid.New(),
		Title:     "title",
		Url:       "url",
		Body:      "body",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectExec("UPDATE articles SET title = \\?, url = \\?, body = \\?, update_at = CURRENT_TIMESTAMP WHERE id = \\?").WillReturnResult(sqlmock.NewResult(1, 1))

	err = UpdateArticle(db, article)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListArticles(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := uuid.New()
	source_id := uuid.New()

	article := models.Article{
		ID:        id,
		Source:    source_id,
		Title:     "title",
		Url:       "url",
		Body:      "body",
		CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "source_id", "title", "url", "body", "created_at"}).
		AddRow(article.ID, article.Source, article.Title, article.Url, article.Body, article.CreatedAt)

	mock.ExpectQuery("SELECT id, source_id, title, url, body, created_at FROM articles").WillReturnRows(rows)

	_, err = ListArticles(db, source_id)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
