package conf

import (
	"github.com/robxu9/kahinah/log"

	"github.com/pelletier/go-toml"
)

var (
	Config *toml.TomlTree
)

func init() {
	var err error
	Config, err = toml.LoadFile("app.toml")
	if err != nil {
		log.Logger.Fatalf("unable to load configuration: %v", err)
	}
}
