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

	newssource := &models.Newssource{
		ID:             uuid.New(),
		Title:          "title",
		Url:            "url",
		FeedType:       "rss",
		UpdatePriority: "URGENT",
		IsActive:       true,
	}

	mock.ExpectExec("INSERT INTO newssources \\(id, title, url, feed_type, update_priority, is_active\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6\\)").WillReturnResult(sqlmock.NewResult(1, 1))

	err = InsertNewssource(db, newssource)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.ExpectExec("INSERT INTO newssources").WillReturnError(nil)
	err = InsertNewssource(db, newssource)

	assert.Error(t, err)
}

func TestFetchNewssource(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	uuid := uuid.New()
	created_at := time.Now()

	rows := sqlmock.NewRows([]string{"id", "title", "url", "feed_type", "update_priority", "is_active", "created_at"}).AddRow(uuid, "title", "url", "rss", "URGENT", true, created_at)
	mock.ExpectQuery("SELECT id, title, url, feed_type, update_priority, is_active, created_at FROM newssources WHERE id = ?").WithArgs(uuid).WillReturnRows(rows)

	newssource, err := FetchNewssource(db, uuid)

	assert.NoError(t, err)
	assert.Equal(t, newssource.ID, uuid)
	assert.Equal(t, newssource.CreatedAt, created_at)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.ExpectExec("SELECT \\* FROM newssources WHERE id = ?").WillReturnError(nil)
	_, err = FetchNewssource(db, uuid)

	assert.Error(t, err)
}

func TestUpdateNewssource(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	newssource := &models.Newssource{
		ID:             uuid.New(),
		Title:          "title",
		Url:            "url",
		FeedType:       "rss",
		UpdatePriority: "URGENT",
		IsActive:       true,
	}

	mock.ExpectExec("UPDATE newssources SET title = \\?, url = \\?, feed_type = \\?, update_priority = \\?, is_active = \\?, updated_at = CURRENT_TIMESTAMP WHERE id = \\?").WithArgs(newssource.Title, newssource.Url, newssource.FeedType, newssource.UpdatePriority, newssource.IsActive, newssource.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = UpdateNewssource(db, newssource)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.ExpectExec("UPDATE newssources SET title").WillReturnError(nil)
	err = UpdateNewssource(db, newssource)

	assert.Error(t, err)

	mockResult := sqlmock.NewResult(0, 0)
	mock.ExpectExec("UPDATE newssources SET title").WillReturnResult(mockResult)
	err = UpdateNewssource(db, newssource)

	assert.Error(t, err)

}

func TestDeleteNewssource(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	uuid := uuid.New()

	mock.ExpectExec("DELETE FROM newssources WHERE id = ?").WithArgs(uuid).WillReturnResult(sqlmock.NewResult(1, 1))

	err = DeleteNewssource(db, uuid)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.ExpectExec("DELETE FROM newssources WHERE id = ?").WillReturnError(nil)
	err = DeleteNewssource(db, uuid)
	assert.Error(t, err)
}

func TestHasNewssources(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	uuid := uuid.New()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(uuid)
	mock.ExpectQuery("SELECT id FROM newssources LIMIT 1").WillReturnRows(rows)

	hasNewssource, err := HasNewssources(db)

	assert.NoError(t, err)
	assert.Equal(t, hasNewssource, true)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.ExpectQuery("SELECT id FROM newssources LIMIT 1").WillReturnError(nil)

	hasNewssource, err = HasNewssources(db)
	assert.Error(t, err)
	assert.Equal(t, hasNewssource, false)
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
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.ExpectQuery("SELECT id, title, url, created_at FROM newssources").WillReturnError(nil)

	_, err = ListNewssources(db)
	assert.Error(t, err)
}
