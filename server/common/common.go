package common

import (
	"log"
	"time"

	"github.com/pmylund/go-cache"
	"github.com/robxu9/kahinah/connectors/abf"
	"github.com/robxu9/kahinah/kahinah"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Common struct {
	K *kahinah.Kahinah
	C *Config

	cache *cache.Cache
}

// Open a common object with the given configuration.
func Open(c *Config) *Common {
	openFunc := kahinah.Open
	if c.DebugMode {
		log.Print("[warn] debug mode enabled")
		openFunc = kahinah.OpenDebug
	}

	k, err := openFunc(c.Database.Dialect, c.Database.Params)
	if err != nil {
		log.Printf("[err] failed to connect to database: %v", err)
		return nil
	}

	// common
	ret := &Common{
		K:     k,
		C:     c,
		cache: cache.New(5*time.Minute, 30*time.Second),
	}

	// set options
	log.Print("setting options...")
	k.AdvisoryProcessFunc = ret.AdvisoryProcessFunc

	if c.Connectors.ABF.Enabled {
		connector := &abf.Connector{
			Platforms:  c.Connectors.ABF.PlatformIds,
			User:       c.Connectors.ABF.User,
			APIKey:     c.Connectors.ABF.APIKey,
			CheckEvery: time.Duration(c.Connectors.ABF.CheckEveryMin) * time.Minute,
		}
		if err := k.Attach(connector); err != nil {
			panic(err)
		}
	}

	// setup additional tables
	k.DB().AutoMigrate(&UserToken{})

	return ret
}

func (c *Common) AdvisoryProcessFunc(a *kahinah.Advisory) kahinah.AdvisoryStatus {
	total := 0

	list := make(map[int64]int)

	for _, v := range a.Comments {
		switch v.Verdict {
		case kahinah.NEUTRAL:
			list[v.UserId] = 0
		case kahinah.NO:
			list[v.UserId] = c.C.Karma.AddFailKarma
		case kahinah.YES:
			list[v.UserId] = c.C.Karma.AddPassKarma
		case kahinah.BLOCK:
			list[v.UserId] = c.C.Karma.AddBlockKarma
		case kahinah.OVERRIDE:
			list[v.UserId] = c.C.Karma.AddOverrideKarma
		}
	}

	for _, v := range list {
		total += v
	}

	if total >= c.C.Karma.PassLimit {
		return kahinah.PASS
	} else if total <= c.C.Karma.FailLimit {
		return kahinah.FAIL
	}

	return kahinah.OPEN
}

func (c *Common) Close() {
	c.K.Close()
}
