package db

import (
	"context"
	"fmt"
)

// registry stores all the loaded DB backends.
var registry = map[string]DB{}

// DB interface defines that a database backend have to implement.
type DB interface {
	// Init makes the DB backend ready to query by create the connection, etc.
	Init(config map[string]interface{}) error
	// Run executes the query in the DB then returns with a JSON (byte slice) for the HTTP response or an error.
	Run(ctx context.Context, query string) (result []byte, err error)
}

// registerDb adds a new database backend to the registry.
func registerDb(name string, db DB) {
	registry[name] = db
}

// GetDb returns the DB with the given name or an error if not found.
func GetDb(name string) (DB, error) {
	if db, ok := registry[name]; ok {
		return db, nil
	}

	return nil, fmt.Errorf("database not found: %s", name)
}
