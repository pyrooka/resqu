package db

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestProcessSQLRows(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	t.Run("invalid columns", func(t *testing.T) {
		rows := sql.Rows{}
		_, err := processSQLRows(&rows)
		if err == nil {
			t.Error("wut? error expected")
		}
	})

	t.Run("empty result", func(t *testing.T) {
		testRows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT t FROM test").WillReturnRows(testRows)

		rows, err := db.Query("SELECT t FROM test")
		if err != nil {
			t.Errorf("error was not expected while querying, but got: %v", err)
		}
		defer rows.Close()

		result, err := processSQLRows(rows)
		if err != nil {
			t.Errorf("error was not expected while processing rows, but got: %v", err)
		}

		if string(result) != "[]" {
			t.Errorf("expected an empty JSON array, but got %s", result)
		}
	})

	t.Run("results types", func(t *testing.T) {
		testRows := sqlmock.NewRows([]string{"id", "name", "height", "birth", "active"})
		testRows.AddRow(0, "John Wick", 179.8, "1985-01-01", true)
		testRows.AddRow(1, "Jane Doe", 165.3, "2000-04-02", false)

		mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(testRows)

		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			t.Errorf("error was not expected while querying, but got: %v", err)
		}
		defer rows.Close()

		result, err := processSQLRows(rows)
		if err != nil {
			t.Errorf("error was not expected while processing rows, but got: %v", err)
		}

		expectedResult := `[{"active":true,"birth":"1985-01-01","height":179.8,"id":0,"name":"John Wick"},{"active":false,"birth":"2000-04-02","height":165.3,"id":1,"name":"Jane Doe"}]`
		if string(result) != expectedResult {
			t.Errorf("results dont't match. Exptected: %s Got: %s", expectedResult, string(result))
		}
	})
}
