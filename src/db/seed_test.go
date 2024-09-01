package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSeed(t *testing.T) {
	dbconn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbconn.Close()

	t.Run("HasNewssources returns empty result - InsertNewssource succeeds", func(t *testing.T) {
		rows := mock.NewRows([]string{"id"})
		mock.ExpectQuery("SELECT id FROM newssources LIMIT 1").WillReturnRows(rows)

		mock.ExpectExec("INSERT INTO newssources").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		Seed(dbconn)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("HasNewssources returns err", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM newssources LIMIT 1").WillReturnError(nil)

		Seed(dbconn)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("InsertNewssource returns err", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"})
		mock.ExpectQuery("SELECT id FROM newssources LIMIT 1").WillReturnRows(rows)

		mock.ExpectExec("INSERT INTO newssources").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(nil)

		Seed(dbconn)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
