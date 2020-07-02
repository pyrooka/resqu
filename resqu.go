package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"

	ccache "github.com/pyrooka/resqu/cache"
	"github.com/pyrooka/resqu/db"
	"github.com/robfig/cron"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

var cronJobs = []string{}

var mutex = &sync.Mutex{}

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
	c, err := readConfig("config.yaml")
	if err != nil {
		log.Fatal("Error while reading the config:", err)
	}

	port, exists := os.LookupEnv("SERVER_PORT")
	if exists != true {
		port = "8888"
	}

	cacheType, exists := os.LookupEnv("CACHE_TYPE")
	if exists != true {
		cacheType = "local"
	}
	cacheType = strings.ToLower(cacheType)

	debug, exists := os.LookupEnv("DEBUG")
	if exists != true {
		debug = "true"
	}
	debug = strings.ToLower(debug)

	if debug == "true" || debug == "1" {
		log.SetLevel(log.DebugLevel)
	}

	// Create a cache.
	cache, err := ccache.NewCache(cacheType)
	if err != nil {
		log.WithFields(log.Fields{
			"cacheType": cacheType,
		}).Fatalf("Error while initializing cache: %s", err)
	}

	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	// Init the database backends and set up the endpoints.
	for n, cc := range c {
		// Copy the vars.
		name := n
		config := cc

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
			cacheConf := e.Cache

			var watcher *cron.Cron
			if cacheConf.Enabled && cacheConf.ClearTime != "" {
				watcher = cron.New()
				watcher.Start()
			}

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

				if cacheConf.Enabled {
					// Read from cache is exists.
					result, err := cache.Get(query)
					if err != nil {
						l := log.WithField("query", query)
						if err == ccache.ErrNotFound {
							l.Infof("[%s] Missing key.", name)
						} else {
							l.Errorf("[%s] Cannot get cache: %s", name, err)
						}

					} else {
						log.WithField("query", query).Infof("[%s] Using result from cache.", name)
						resp := response{
							Data: result,
						}

						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(resp)

						return
					}
				}

				log.WithFields(log.Fields{
					"vars": vars,
				}).Debugf("[%s] Executing a query: %s", name, query)

				// BOOOM.
				result, err := db.Run(r.Context(), query)
				if err != nil {
					log.WithFields(log.Fields{
						"query": query,
					}).Errorf("[%s] Error while executing the query: %s", name, err)

					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if cacheConf.Enabled {
					// Cache the result.
					err = cache.Set(query, result)
					if err != nil {
						log.WithFields(log.Fields{
							"query":      query,
							"resultSize": len(result),
						}).Errorf("[%s] Cannot set to cache: %s", name, err)
					} else {
						// Check if we already have a cron job for this query.
						if cacheConf.ClearTime != "" && !stringInSlice(cronJobs, query) {
							err = watcher.AddFunc(cacheConf.ClearTime, func() {
								err = cache.Remove(query)
								if err != nil {
									l := log.WithField("query", query)
									if err == ccache.ErrNotFound {
										l.Infof("[%s] Missing key.", name)
									} else {
										l.Errorf("[%s] Cannot remove key from cache: %s", name, err)
									}
								} else {
									log.Infof("[%s] Removed from cache.", name)
								}
							})
							if err != nil {
								log.WithFields(log.Fields{
									"clearTime": cacheConf.ClearTime,
									"query":     query,
								}).Errorf("[%s] Error while creating a cron function: %s", name, err)
							} else {
								mutex.Lock()
								cronJobs = append(cronJobs, query)
								mutex.Unlock()
							}
						}
					}
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
