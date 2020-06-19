package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type endpoint struct {
	URL   string `yaml:"url"`
	Query string `yaml:"query"`
}

type dbConfig struct {
	Connection map[string]interface{} `yaml:"connection"`
	Endpoints  []endpoint             `yaml:"endpoints"`
}

func readConfig(file string) (c map[string]dbConfig, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &c)

	return
}
