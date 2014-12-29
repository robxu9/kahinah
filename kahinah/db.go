package kahinah

import (
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pmylund/go-cache"
)

// Kahinah represents the main entry point into all functions
type Kahinah struct {
	// AdvisoryID for AdvisoryFamily mutex
	advisoryFamIDmutex *sync.Mutex

	// A mutex for the db is not needed.
	// In gorm, the db is cloned before any operation, so variables
	// occur separately from each other, it seems.
	db    *gorm.DB
	cache *cache.Cache
}

// Open sets up the Kahinah database, creating and
// adjusting tables as necessary. It follows the standard
// sql.Open() syntax.
func Open(dialect, params string) (*Kahinah, error) {
	db, err := gorm.Open(dialect, params)

	if err != nil {
		return nil, err
	}

	dbptr := &db

	// Auto-Migration

	// migrate updates
	dbptr.AutoMigrate(&Update{}, &UpdateChange{}, &UpdateContent{}, &UpdatePackage{})
	// migrate connectors
	dbptr.AutoMigrate(&ConnectUpdateContentUpdateChange{}, &ConnectUpdateContentUpdatePackage{})

	// migrate advisories
	dbptr.AutoMigrate(&Advisory{}, &Comment{})

	// migrate users
	dbptr.AutoMigrate(&User{}, &UserIP{}, &UserToken{})

	// Create cache
	c := cache.New(5*time.Minute, 30*time.Second)

	return &Kahinah{db: dbptr, cache: c, advisoryFamIDmutex: &sync.Mutex{}}, nil
}

// OpenDebug is like open, except it enables debug logging. DO NOT
// ENABLE THIS IN PRODUCTION - it logs EVERYTHING!
func OpenDebug(dialect, params string) (*Kahinah, error) {
	k, err := Open(dialect, params)
	if err != nil {
		return nil, err
	}

	k.db.LogMode(true)
	return k, nil
}

// Close closes the database. Any operations afterwards on Kahinah
// WILL panic.
func (k *Kahinah) Close() error {
	// delete everything from the cache
	k.cache.Flush()
	// close the db
	return k.db.Close()
}
