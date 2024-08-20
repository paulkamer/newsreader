package repositories

import (
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

	article := models.Article{
		ID:        uuid.New().String(),
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
