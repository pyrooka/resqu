package db

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3" // SQLite 3 driver
)

// SQLite struct
type SQLite struct {
	client *sql.DB
}

func init() {
	registerDb("sqlite", &SQLite{})
}

// Init DB interface implementation.
func (sqlt *SQLite) Init(config map[string]interface{}) (err error) {
	var dbPath string

	// Check the necessary variables.
	if path, ok := config["path"]; ok {
		dbPath = path.(string)
	} else {
		return errors.New("path is missing from the connection config")
	}

	sqlt.client, err = sql.Open("sqlite3", dbPath)

	return
}

// Run DB interface implementation.
func (sqlt *SQLite) Run(ctx context.Context, query string) (result []byte, err error) {
	rows, err := sqlt.client.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()

	result, err = processSQLRows(rows)

	return
}
