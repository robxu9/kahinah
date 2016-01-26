package main

import (
	"log"

	"github.com/robxu9/kahinah/backend"
	"github.com/unrolled/render"
)

type Server1 struct {
	R *render.Render       // renderer
	K *kahinah.K           // kahinah
	C *Config              // config
	M map[string]Connector // connectors
}

type Connector interface {
	Name() string
	Init(*kahinah.K) error
	Pass(*kahinah.Update)
	Fail(*kahinah.Update)
	Close()
}

var (
	Close = make(chan struct{})
)

func main() {
	// starting up server1
	log.Printf("starting up server1...")

	// server1 should read configuration
	log.Printf("reading configuration...")
	c := parseConfig("config.toml")

	// server1 should call kahinah
	log.Printf("initialising kahinah...")
	k, err := kahinah.Open(&kahinah.KOpts{})
	if err != nil {
		log.Fatalf("unable to initialise kahinah: %v", err)
	}

	// server1 should manage abf, including initialising rules

	// server1 should start the rss reader

	// server1 should initialise the renderer

	// and finally, server1 should start its frontend stack
	server1 := &Server1{
		K: k,
		C: c,
		M: map[string]Connector{},
	}

	<-Close
	log.Printf("stopping server1...")
	if err = server1.K.Close(); err != nil {
		log.Printf("warning: failed to close kahinah connection: %v", err)
	}
	for _, v := range server1.M {
		v.Close()
	}
}
