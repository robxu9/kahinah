package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/robxu9/kahinah/server/apiv1"
	"github.com/robxu9/kahinah/server/common"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/negroni"
	"github.com/stretchr/graceful"
)

const (
	// VERSION contains the server version. Bump for configuration changes.
	VERSION = 0
)

func defaultAndExit() {
	configFile, err := os.Create("config.toml.new")
	if err != nil {
		panic(err)
	}

	encoder := toml.NewEncoder(configFile)
	encoder.Indent = "\t"

	if err := encoder.Encode(common.DefaultConfig(VERSION)); err != nil {
		panic(err)
	}

	configFile.Close()
	os.Exit(2)
}

func main() {
	logFile, err := os.OpenFile("output.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Printf("[warn] failed to open output.log: %v", err)
	} else {
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
		defer logFile.Close()
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("starting kahinah, version %v...", VERSION)

	// read configuration
	log.Print("reading configuration...")

	configFile, err := os.Open("config.toml")
	if err != nil {
		log.Print("[err] failed to open config.toml")
		log.Print("[err] writing default config to config.toml.new and exiting")
		defaultAndExit()
	}
	defer configFile.Close()

	config := common.DefaultConfig(VERSION)

	if _, err := toml.DecodeReader(configFile, config); err != nil {
		log.Print("[err] failed to decode config.toml")
		log.Print("[err] writing default config to config.toml.new and exiting")
		defaultAndExit()
	}

	if config.Version != VERSION {
		log.Print("[err] config version is not the same!")
		log.Print("[err] writing default config to config.toml.new and exiting")
		log.Print("[err] modify your existing config to add new variables and bump version")
		defaultAndExit()
	}

	// create common object
	log.Print("opening database...")

	common := common.Open(config)
	if common == nil {
		log.Fatal("[err] failed to open database")
		os.Exit(1)
	}
	defer common.Close()

	// start http server
	log.Print("creating routes...")
	mux := http.NewServeMux()

	// API handlers
	// ~> /api/v1: json api endpoint
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiv1.New(common)))
	// ~> /: ember.js client
	mux.Handle("/", &ClientEndpoint{})

	// start middleware
	log.Print("starting middleware...")
	n := negroni.Classic()
	n.UseHandler(mux)

	log.Printf("running on %v", config.HTTP)
	graceful.Run(config.HTTP, 10*time.Second, n)

	log.Printf("*** gracefully terminating")
}
