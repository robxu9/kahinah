package kahinah

import (
	"time"

	rt "gopkg.in/dancannon/gorethink.v1"
)

const (
	TableUpdate = "updates"
)

var (
	TimeMin = time.Unix(0, 0)
	TimeMax = time.Unix(1<<63-1, 999999999)
)

// Update represents an update submitted.
type Update struct {
	// rethinkdb unique id
	ID string `gorethink:"id,omitempty"`

	// common update attributes (indexes)
	Platform  string    `gorethink:"platform"`
	Name      string    `gorethink:"name"`
	EVR       string    `gorethink:"evr"`
	Submitter string    `gorethink:"submitter"`
	Date      time.Time `gorethink:"created_at"`
	Type      string    `gorethink:"type"`
	Connector string    `gorethink:"connector"` // which connector maintains this update

	// non-indexable update attributes usually
	Diff      []byte `gorethink:"diff"`
	Changelog []byte `gorethink:"changelog"`

	// connector tampering updates
	Advisory       string            `gorethink:"advisory_id"`
	Deprecated     bool              `gorethink:"deprecated"`
	Pushed         bool              `gorethink:"pushed"`
	ConnectorStore map[string]string `gorethink:"connector_store"` // connector specific info

	// list of new packages - arch -> packages
	Packages map[string][]string `gorethink:"packages"`
}

// InsertUpdate inserts an update into the database. This is mutable.
func (k *K) InsertUpdate(u *Update) string {
	rw, err := k.RunWriteErr(k.DB().Table(TableUpdate).Insert(u, rt.InsertOpts{}))
	if err != nil {
		panic(err)
	}

	if len(rw.GeneratedKeys) != 1 {
		panic("# generated keys should only be 1")
	}

	u.ID = rw.GeneratedKeys[0]
	return rw.GeneratedKeys[0]
}

// GetUpdate gets a specific update with the specific uuid, or nil if it doesn't
// exist.
func (k *K) GetUpdate(uuid string) *Update {
	cursor, err := k.Run(k.DB().Table(TableUpdate).Get(uuid))
	if err != nil {
		panic(err)
	}

	var u *Update

	defer cursor.Close()
	if err = cursor.One(u); err != nil {
		panic(err)
	}

	return u
}

// SearchUpdate is used to search for updates in the system.
type SearchUpdate struct {
	Platform   string `gorethink:"platform,omitempty"`
	Name       string `gorethink:"name,omitempty"`
	EVR        string `gorethink:"evr,omitempty"`
	Submitter  string `gorethink:"submitter,omitempty"`
	Type       string `gorethink:"type,omitempty"`
	Connector  string `gorethink:"connector,omitempty"` // which connector maintains this update
	Deprecated bool   `gorethink:"deprecated,omitempty"`
	Pushed     bool   `gorethink:"pushed,omitempty"`

	DateBefore time.Time `gorethink:"-"`
	DateAfter  time.Time `gorethink:"-"`
}

// SearchUpdate searches for an update with the specified parameters.
// Unfortunately, for more powerful functionality, direct manipulation of the
// DB is needed.
func (k *K) SearchUpdate(search SearchUpdate) []*Update {
	table := k.DB().Table(TableUpdate)
	if !search.DateBefore.Equal(time.Time{}) || !search.DateAfter.Equal(time.Time{}) {
		// 'created_at' - see update struct
		var lower interface{} = search.DateBefore
		if search.DateBefore.Equal(TimeMin) {
			lower = rt.MinVal
		}

		var higher interface{} = search.DateAfter
		if search.DateAfter.Equal(TimeMax) {
			higher = rt.MaxVal
		}

		table = table.Between(lower, higher, rt.BetweenOpts{
			Index: "created_at",
		})
	}

	cursor, err := k.Run(table.Filter(search))
	if err != nil {
		panic(err)
	}

	defer cursor.Close()

	var updates []*Update
	err = cursor.All(&updates)
	if err != nil {
		panic(err)
	}

	return updates
}
