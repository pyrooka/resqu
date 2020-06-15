package db

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v4/stdlib" // PostgreSQL driver
)

// PostgreSQL struct
type PostgreSQL struct {
	client *sql.DB
}

func init() {
	registerDb("postgresql", &PostgreSQL{})
}

// Init DB interface implementation.
func (psql *PostgreSQL) Init(config map[string]interface{}) (err error) {
	var connectionURL string

	// Check the necessary variables.
	if path, ok := config["connectionURL"]; ok {
		connectionURL = path.(string)
	} else {
		return errors.New("connectionURL is missing from the connection config")
	}

	psql.client, err = sql.Open("pgx", connectionURL)

	return
}

// Run DB interface implementation.
func (psql *PostgreSQL) Run(ctx context.Context, query string) (result []byte, err error) {
	rows, err := psql.client.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()

	result, err = processSQLRows(rows)

	return
}
