package db

import (
	"context"
	"testing"
)

type DBMock struct{}

func (db *DBMock) Init(config map[string]interface{}) (err error) {
	return
}

func (db *DBMock) Run(ctx context.Context, query string) (result []byte, err error) {
	return
}

func TestRegisterDB(t *testing.T) {
	registerDb("testDB", &DBMock{})

	if _, exists := registry["testDB"]; !exists {
		t.Errorf("WTF? registry: %s", registry)
	}
}

func TestGetDB(t *testing.T) {
	registerDb("testDB", &DBMock{})

	t.Run("test missing DB", func(t *testing.T) {
		_, err := GetDb("missingDB")
		if err == nil && err.Error() != "database not found missingDB" {
			t.Error("this DB should missing")
		}
	})

	t.Run("test existing DB", func(t *testing.T) {
		_, err := GetDb("testDB")
		if err != nil {
			t.Error("this DB should not be missing")
		}
	})
}
