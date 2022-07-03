package config

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v3"
)

func Parse(filePath string) (*Config, error) {
	configData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	cfg := NewConfig()
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, err
}
