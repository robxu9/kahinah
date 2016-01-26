package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/naoina/toml"
)

type Config struct {
}

func parseConfig(file string) *Config {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal("failed to open configuration:", err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("couldn't read all configuration:", err)
	}
	var config Config
	if err := toml.Unmarshal(buf, &config); err != nil {
		log.Fatal("couldn't unmarshal configuration:", err)
	}

	return &config
}
