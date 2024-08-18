package repositories

import (
	"newsreader/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInsertNewsources(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	newssource := models.Newssource{
		ID:             uuid.New(),
		Title:          "title",
		Url:            "url",
		UpdatePriority: "URGENT",
		IsActive:       true,
	}

	mock.ExpectExec("INSERT INTO newssources \\(id, title, url, update_priority, is_active\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\)").WillReturnResult(sqlmock.NewResult(1, 1))

	err = InsertNewssource(db, newssource)

	assert.NoError(t, err)
}

func TestFetchNewssource(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	uuid := uuid.New()
	created_at := time.Now()

	rows := sqlmock.NewRows([]string{"id", "title", "url", "update_priority", "is_active", "created_at", "updated_at"}).AddRow(uuid, "title", "url", "URGENT", true, created_at, created_at)
	mock.ExpectQuery("SELECT \\* FROM newssources WHERE id = ?").WithArgs(uuid).WillReturnRows(rows)

	newssource, err := FetchNewssource(db, uuid)

	assert.NoError(t, err)
	assert.Equal(t, newssource.ID, uuid)
	assert.Equal(t, newssource.CreatedAt, created_at)
}

func TestListNewssources(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	uuid := uuid.New()
	created_at := time.Now()

	rows := sqlmock.NewRows([]string{"id", "title", "url", "created_at"}).AddRow(uuid, "title", "url", created_at)
	mock.ExpectQuery("SELECT id, title, url, created_at FROM newssources").WillReturnRows(rows)

	newssources, err := ListNewssources(db)

	assert.NoError(t, err)
	assert.Equal(t, newssources[0].ID, uuid)
	assert.Equal(t, newssources[0].CreatedAt, created_at)
}
