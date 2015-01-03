package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/robxu9/kahinah/connectors/abf"
	"github.com/robxu9/kahinah/kahinah"
	"github.com/unrolled/render"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	VERSION = 0
)

var (
	global *kahinah.Kahinah
	config Config
	r      *render.Render
)

func defaultAndExit() {
	configFile, err := os.Create("config.toml.new")
	if err != nil {
		panic(err)
	}

	encoder := toml.NewEncoder(configFile)
	encoder.Indent = "\t"

	if err := encoder.Encode(configFile); err != nil {
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

	if _, err := toml.DecodeReader(configFile, &config); err != nil {
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

	// connect to kahinah database
	log.Print("opening database...")

	openFunc := kahinah.Open
	if config.DebugMode {
		log.Print("[warn] DEBUG MODE ENABLED")
		openFunc = kahinah.OpenDebug
	}

	global, err = openFunc(config.Database.Dialect, config.Database.Params)
	if err != nil {
		log.Fatalf("[err] failed to connect to database: %v", err)
	}

	// set options
	log.Print("setting options...")
	global.AdvisoryProcessFunc = AdvisoryProcessFunc
	c := cookiestore.New([]byte(config.SecretKey))

	if config.Connectors.ABF.Enabled {
		connector := &abf.Connector{
			Platforms:  config.Connectors.ABF.PlatformIds,
			User:       config.Connectors.ABF.User,
			APIKey:     config.Connectors.ABF.APIKey,
			CheckEvery: time.Duration(config.Connectors.ABF.CheckEveryMin) * time.Minute,
		}
		if err := global.Attach(connector); err != nil {
			panic(err)
		}
	}

	// start renderer
	log.Print("starting renderer...")
	r = render.New(render.Options{
		Layout:        "layout",
		Extensions:    []string{".tmpl", ".html"},
		IsDevelopment: config.DevMode,
	})

	// start http server
	log.Print("creating routes...")
	mux := http.NewServeMux()

	// start middleware
	log.Print("starting middleware...")
	n := negroni.Classic()
	n.Use(sessions.Sessions("kahinah", c))
	n.UseHandler(mux)

	log.Printf("running on %v", config.HTTP)
	n.Run(config.HTTP)
}
