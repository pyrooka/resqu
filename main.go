package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pyrooka/resqu/db"

	"github.com/gorilla/mux"
)

type response struct {
	Data  json.RawMessage `json:"data,omitempty"`
	Error error           `json:"error,omitempty"`
}

func main() {
	// Read and parse the config.
	c, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	port, exists := os.LookupEnv("SERVER_PORT")
	if exists != true {
		port = "8888"
	}

	router := mux.NewRouter()

	// Init the database backends and set up the endpoints.
	for name, config := range c {
		// Get the DB from the name of the map in the yaml.
		db, err := db.GetDb(name)
		if err != nil {
			log.Fatal(err)
		}

		// Initialize this DB backend.
		err = db.Init(config.Connection)
		if err != nil {
			log.Fatalf("Error while initializing %s: %s", name, err)
		}

		// Create the endpoints.
		for _, e := range config.Endpoints {
			// Copy the values.
			URL := e.URL
			rawQuery := e.Query

			router.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				params := r.URL.Query()
				query := rawQuery

				// Replace the mux variables and the parameters in the SQL query with the values from the URL.
				for muxVar, value := range vars {
					query = strings.ReplaceAll(query, fmt.Sprintf("{%s}", muxVar), value)
				}
				for param, values := range params {
					if len(values) > 1 {
						log.Println(fmt.Sprintf(`Ooops. There are more than 1 parameter for key "%s" in the query: %s`, param, values))
						continue
					}
					// Only use the first value for each parameter in the HTTP query.
					query = strings.ReplaceAll(query, fmt.Sprintf("{%s}", param), values[0])
				}

				// TODO: implement proper context handling.
				ctx := context.Background()

				// BOOOM.
				result, err := db.Run(ctx, query)
				if err != nil {
					log.Println("Error while executing the query:", err)
				}

				resp := response{
					Data:  result,
					Error: err,
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			})
		}
	}

	addr := fmt.Sprintf("0.0.0.0:%s", port)

	log.Println("HTTP server listening on:", addr)

	log.Fatal(http.ListenAndServe(addr, router))
}
