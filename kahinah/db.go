package kahinah

import (
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pmylund/go-cache"
)

// Kahinah represents the main entry point into all functions
type Kahinah struct {
	// Function to determine whether an advisory passes, fails, or remains the same
	AdvisoryProcessFunc func(*Advisory) AdvisoryStatus

	// Process changed Advisories
	advisoryProcessQueue    chan *Advisory
	advisoryProcessRoutines *sync.WaitGroup

	// AdvisoryID for AdvisoryFamily mutex
	advisoryFamIDmutex *sync.Mutex

	// User mutex (in case frontend is stupid enough to try creating two users w/ same email)
	userMutex *sync.Mutex

	// Connectors
	connectors map[string]Connector

	// A mutex for the db is not needed.
	// In gorm, the db is cloned before any operation, so variables
	// occur separately from each other, it seems.
	db    *gorm.DB
	cache *cache.Cache
}

// DefaultAdvisoryProcessFunc provides the default advisory process function,
// a simple karma check, with limits -3 and 3.
func DefaultAdvisoryProcessFunc(a *Advisory) AdvisoryStatus {
	total := 0
	for _, v := range a.Comments {
		switch v.Verdict {
		case NEUTRAL:
			break
		case NO:
			total--
		case YES:
			total++
		case BLOCK:
			total -= 9999
		case OVERRIDE:
			total += 9999
		}
	}

	if total >= 3 {
		return PASS
	} else if total <= -3 {
		return FAIL
	}

	return OPEN
}

// Open sets up the Kahinah database, creating and
// adjusting tables as necessary. It follows the standard
// sql.Open() syntax.
func Open(dialect, params string) (*Kahinah, error) {
	return open(dialect, params, false)
}

// OpenDebug is like open, except it enables debug logging. DO NOT
// ENABLE THIS IN PRODUCTION - it logs EVERYTHING!
func OpenDebug(dialect, params string) (*Kahinah, error) {
	return open(dialect, params, true)
}

func open(dialect, params string, debug bool) (*Kahinah, error) {
	db, err := gorm.Open(dialect, params)

	if err != nil {
		return nil, err
	}

	dbptr := &db

	if debug {
		dbptr.LogMode(true)
	}

	// Auto-Migration

	// migrate updates
	dbptr.AutoMigrate(&Update{}, &UpdateChange{}, &UpdateContent{}, &UpdatePackage{})
	// migrate connectors
	dbptr.AutoMigrate(&ConnectUpdateContentUpdateChange{}, &ConnectUpdateContentUpdatePackage{})

	// migrate advisories
	dbptr.AutoMigrate(&Advisory{}, &Comment{})

	// migrate users
	dbptr.AutoMigrate(&User{})

	// Create cache
	c := cache.New(5*time.Minute, 30*time.Second)

	// Create kahinah obj
	kahinah := &Kahinah{
		connectors:              make(map[string]Connector),
		db:                      dbptr,
		cache:                   c,
		advisoryFamIDmutex:      &sync.Mutex{},
		userMutex:               &sync.Mutex{},
		advisoryProcessQueue:    make(chan *Advisory, 100),
		advisoryProcessRoutines: &sync.WaitGroup{},
		AdvisoryProcessFunc:     DefaultAdvisoryProcessFunc,
	}

	// Start worker queue
	go kahinah.processAdvisory()

	return kahinah, nil
}

// DB returns a copy of the database object for client side use.
func (k *Kahinah) DB() *gorm.DB {
	return k.db.New()
}

// Close closes the database. Any operations afterwards on Kahinah
// WILL panic.
func (k *Kahinah) Close() error {
	// finish processing everything
	close(k.advisoryProcessQueue)
	k.advisoryProcessRoutines.Wait() // we must wait for processes to finish
	// signal to connectors that we're closing
	for _, v := range k.connectors {
		v.Close()
	}
	// delete everything from the cache
	k.cache.Flush()
	// close the db
	return k.db.Close()
}
