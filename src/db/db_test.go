package db

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInitDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	_, err = InitDatabase("sqlite3", "")
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Test returned error is handled

	// mock.ExpectExec("CREATE TABLE IF NOT EXISTS newssources").
	// 	WillReturnError(fmt.Errorf("error"))

	// _, err = InitDatabase("sqlite3", "")

	// assert.Error(t, err)

	// if err := mock.ExpectationsWereMet(); err != nil {
	// 	t.Errorf("there were unfulfilled expectations: %s", err)
	// }
}

func TestConnect(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	conn, err := Connect("sqlite3", "")

	assert.NoError(t, err)
	assert.NotNil(t, conn)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateDbStructure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS newssources").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS articles").WillReturnResult(sqlmock.NewResult(1, 1))

	err = createDbStructure(db)
	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Test returned error is handled

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS newssources").
		WillReturnError(fmt.Errorf("error"))

	err = createDbStructure(db)

	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
