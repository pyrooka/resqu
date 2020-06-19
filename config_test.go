package main

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	t.Run("Test missing config", func(t *testing.T) {
		_, err := readConfig("config.mp3")
		if err == nil {
			t.Error(`error shouldn't be nil because "config.mp3" is probably missing`)
		}

	})

	t.Run("Test missing config", func(t *testing.T) {
		config, err := readConfig("examples/config.yaml")
		if err != nil {
			t.Error("where is the config from the examples")
		}

		tDbConfig := map[string]dbConfig{
			"sqlite": {
				Connection: map[string]interface{}{"path": "/db/employees.sqlite3"},
				Endpoints: []endpoint{
					{
						URL:   "/employees",
						Query: "SELECT * FROM emp {{if .limit}} LIMIT {{.limit}} {{end}} {{if .offset}} OFFSET {{.offset}} {{end}}",
					},
					{
						URL:   "/employees/{empNo}",
						Query: "SELECT * FROM emp WHERE empno = {{.empNo}}",
					},
				},
			},
		}

		t.Run("Test config", func(t *testing.T) {
			/*
				if reflect.DeepEqual(config, tDbConfig) != true {
					t.Error("mismatching configs")
				}
			*/
			if config["sqlite"].Connection["path"] != tDbConfig["sqlite"].Connection["path"] {
				t.Error("wrong path")
			}

			for i, e := range config["sqlite"].Endpoints {
				if e.URL != tDbConfig["sqlite"].Endpoints[i].URL {
					t.Errorf("mismatching URL: %d", i)
				}
				if e.Query != tDbConfig["sqlite"].Endpoints[i].Query {
					t.Errorf("mismatching query: %d", i)
				}
			}
		})

	})
}
