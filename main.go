package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/pyrooka/resqu/db"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type response struct {
	Data json.RawMessage `json:"data"`
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"user-agent": r.UserAgent(),
		}).Infof("%s %s", r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Read and parse the config.
	c, err := readConfig()
	if err != nil {
		log.Fatal("Error while reading the config:", err)
	}

	port, exists := os.LookupEnv("SERVER_PORT")
	if exists != true {
		port = "8888"
	}

	debug, exists := os.LookupEnv("DEBUG")
	if exists != true {
		debug = "true"
	}
	debug = strings.ToLower(debug)

	if debug == "true" || debug == "1" {
		log.SetLevel(log.DebugLevel)
	}

	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	// Init the database backends and set up the endpoints.
	for name, config := range c {
		// Get the DB from the name of the map in the yaml.
		db, err := db.GetDb(name)
		if err != nil {
			log.WithFields(log.Fields{
				"config": config,
			}).Fatalf("[%s] Error while loading the DB: %s", name, err)
		}

		// Initialize this DB backend.
		err = db.Init(config.Connection)
		if err != nil {
			log.WithFields(log.Fields{
				"config": config,
			}).Fatalf("[%s] Error while initializing the DB: %s", name, err)
		}

		log.Infof("[%s] Initialized.", name)

		// Create the endpoints.
		for i, e := range config.Endpoints {
			// Copy the values.
			URL := e.URL
			rawQuery := e.Query

			// Build the template.
			t, err := template.New(fmt.Sprintf("%s_%d", name, i)).Parse(rawQuery)
			if err != nil {
				log.WithFields(log.Fields{
					"rawQuery": rawQuery,
				}).Fatalf("[%s] Error building the template from the query: %s", name, err)
			}

			router.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				params := r.URL.Query()

				// Check the parameters and add them to the vars.
				for param, values := range params {
					if len(values) > 1 {
						log.WithFields(log.Fields{
							"vars":   vars,
							"params": params,
						}).Errorf(`[%s] Ooops. There are more than 1 parameter for key "%s".`, name, param)

						http.Error(w, fmt.Sprintf(`there are more than 1 parameter for key "%s"`, param), http.StatusBadRequest)
						return
					}

					vars[param] = values[0]
				}

				q := new(bytes.Buffer)

				err = t.Execute(q, vars)
				if err != nil {
					log.WithFields(log.Fields{
						"vars":     vars,
						"rawQuery": rawQuery,
					}).Errorf(`[%s]. Error while executing the template "%s".`, name, err)

					http.Error(w, "Error while executing the template. Check the server logs.", http.StatusInternalServerError)
					return
				}

				query := q.String()

				// TODO: implement proper context handling.
				ctx := context.Background()

				log.WithFields(log.Fields{
					"vars": vars,
				}).Debugf("[%s] Executing a query: %s", name, query)

				// BOOOM.
				result, err := db.Run(ctx, query)
				if err != nil {
					log.WithFields(log.Fields{
						"query": query,
					}).Error("[%s] Error while executing the query:", name, err)

					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				resp := response{
					Data: result,
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			})
		}
	}

	addr := fmt.Sprintf("0.0.0.0:%s", port)

	log.WithFields(log.Fields{
		"Address": addr,
	}).Info("HTTP server is listening")

	log.Fatal(http.ListenAndServe(addr, router))
}
