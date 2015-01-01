package kahinah

import (
	"errors"
	"strconv"
	"time"
)

// UpdateType defines the various update types
type UpdateType int

const (
	// NONE represents an undefined update type
	NONE UpdateType = iota
	// BUGFIX represents a bugfix update
	BUGFIX
	// SECURITY represents a security update
	SECURITY
	// ENHANCEMENT represents an update that provides improvements
	ENHANCEMENT
	// NEW represents a new package
	NEW
)

var (
	// ErrNoSuchUpdate - update does not exist.
	ErrNoSuchUpdate = errors.New("kahinah: update doesn't exist")

	// ErrCorruptUpdate - update is missing components, check db
	ErrCorruptUpdate = errors.New("kahinah: update is missing components - check db")

	// CacheUpdateExp - the default expiration for updates in the cache
	CacheUpdateExp = 12 * time.Hour
)

// Update represents an update that has been received from
// a build system (typically via connector). This update
// is immutable and is assigned a unique id for which it is
// referred by.
type Update struct {
	Id int64

	AdvisoryId int64 // attached to certain advisory

	Connector string     // connector that handled this update
	For       string     // update target
	Name      string     // update name
	Submitter string     // update submitter (their email)
	Type      UpdateType // update type

	Content   *UpdateContent `sql:"-"`       // contents of the update
	ContentId int64          `sql:"content"` // content ID

	ConnectorId   string // connector id
	ConnectorInfo string // connector information

	CreatedAt time.Time // time this was created
}

// UpdateContent represents the contents of an update (for
// example, changelog, subpackages, and a diff).
type UpdateContent struct {
	Id int64

	From string // from commit id
	To   string // to commit id

	Url     string    // url on the build system
	BuiltAt time.Time // time this requested for building

	Packages []*UpdatePackage `sql:"-"` // list of update packages

	Changes []*UpdateChange `sql:"-"` // list of changelog entries
}

// ConnectUpdateContentUpdatePackage is a sql relation connector
type ConnectUpdateContentUpdatePackage struct {
	Id              int64
	UpdateContentId int64
	UpdatePackageId int64
}

// UpdatePackage represents a package in the update
type UpdatePackage struct {
	Id int64

	// name [epoch]:version-release.arch.type
	Name    string
	Epoch   uint64
	Version string
	Release string
	Arch    string
	Type    string

	// direct url to the package (optional)
	Url string
}

// ConnectUpdateContentUpdateChange is a sql relation connector
type ConnectUpdateContentUpdateChange struct {
	Id              int64
	UpdateContentId int64
	UpdateChangeId  int64
}

// UpdateChange is a changelog entry in the update
type UpdateChange struct {
	Id int64

	ChangeAt time.Time
	For      string // epoch:version-release
	By       string // who made the entry
	Details  string
}

func updateCacheID(id int64) string {
	return "updates/" + strconv.FormatInt(id, 10)
}

func (k *Kahinah) setUpdateAdvisoryId(u *Update, id int64) error {
	u.AdvisoryId = id // we save ptrs in the cache so this should be okay
	return k.db.Save(u).Error
}

// NewUpdate creates a new update and stores it in the database. It then returns
// either a unique id to the update, or an error with why it failed.
func (k *Kahinah) NewUpdate(connector, target, name, submitter string, updatetype UpdateType, content *UpdateContent, connectorid, connectorinfo string) (int64, error) {
	// Insert content first
	if err := k.db.Save(content).Error; err != nil {
		return 0, err
	}

	// For each package, insert and create connector
	for _, v := range content.Packages {
		if err := k.db.Save(v).Error; err != nil {
			return 0, err
		}
		connector := &ConnectUpdateContentUpdatePackage{
			UpdateContentId: content.Id,
			UpdatePackageId: v.Id,
		}
		if err := k.db.Save(connector).Error; err != nil {
			return 0, err
		}
	}

	// Do the same for each changelog entry
	for _, v := range content.Changes {
		if err := k.db.Save(v).Error; err != nil {
			return 0, err
		}
		connector := &ConnectUpdateContentUpdateChange{
			UpdateContentId: content.Id,
			UpdateChangeId:  v.Id,
		}
		if err := k.db.Save(connector).Error; err != nil {
			return 0, err
		}
	}

	record := &Update{
		Connector:     connector,
		For:           target,
		Name:          name,
		Submitter:     submitter,
		Type:          updatetype,
		Content:       content,
		ContentId:     content.Id,
		ConnectorId:   connectorid,
		ConnectorInfo: connectorinfo,
	}

	if err := k.db.Save(record).Error; err != nil {
		return 0, err
	}

	// store in the cache
	k.cache.Set(updateCacheID(record.Id), record, CacheUpdateExp)

	return record.Id, nil
}

// RetrieveUpdate retrieves an update from the database (or cache).
func (k *Kahinah) RetrieveUpdate(id int64) (*Update, error) {
	if cached, found := k.cache.Get(updateCacheID(id)); found {
		record := cached.(*Update)
		return record, nil
	}

	record := &Update{}

	if k.db.First(record, id).RecordNotFound() {
		return nil, ErrNoSuchUpdate
	}

	record.Content = &UpdateContent{}

	// Retrieve connecting content
	if k.db.First(record.Content, record.ContentId).RecordNotFound() {
		return nil, ErrCorruptUpdate
	}

	// Then retrieve connecting packages to content
	var pkgconnectors []ConnectUpdateContentUpdatePackage
	if err := k.db.Where(&ConnectUpdateContentUpdatePackage{UpdateContentId: record.Content.Id}).Find(&pkgconnectors).Error; err != nil {
		panic(err) // only possibly explaination is db error?
	}
	record.Content.Packages = make([]*UpdatePackage, len(pkgconnectors))
	for i, v := range pkgconnectors {
		pkg := &UpdatePackage{}
		if err := k.db.First(pkg, v.UpdatePackageId).Error; err != nil {
			panic(ErrCorruptUpdate)
		}
		record.Content.Packages[i] = pkg
	}

	// And connecting changelog entries to content
	var clconnectors []ConnectUpdateContentUpdateChange
	if err := k.db.Where(&ConnectUpdateContentUpdateChange{UpdateContentId: record.Content.Id}).Find(&clconnectors).Error; err != nil {
		panic(err) // only possibly explaination is db error?
	}
	record.Content.Changes = make([]*UpdateChange, len(clconnectors))
	for i, v := range clconnectors {
		change := &UpdateChange{}
		if err := k.db.First(change, v.UpdateChangeId).Error; err != nil {
			panic(ErrCorruptUpdate)
		}
		record.Content.Changes[i] = change
	}

	// store in the cache
	k.cache.Set(updateCacheID(record.Id), record, CacheUpdateExp)

	return record, nil
}

// CountUpdates counts the number of updates in the database.
func (k *Kahinah) CountUpdates() int64 {
	// panic if unable to get count [that sounds fatal]

	var count int64

	if err := k.db.Model(&Update{}).Count(&count).Error; err != nil {
		panic(err)
	}

	return count
}

// ListUpdates returns a list of updates in most recent order.
func (k *Kahinah) ListUpdates(from, limit int64) ([]int64, error) {
	var records []int64

	if err := k.db.Model(&Update{}).Order("created_at desc").Limit(limit).Offset(from).Pluck("id", &records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

// FindUpdatesWithConnector looks for an update with the specified connector name, id, and/or info.
// Leave any blank to include all results.
func (k *Kahinah) FindUpdatesWithConnector(name, id, info string) ([]int64, error) {
	var records []int64

	if err := k.db.Model(&Update{}).Where(&Update{Connector: name, ConnectorId: id, ConnectorInfo: info}).Pluck("id", &records).Error; err != nil {
		return nil, err
	}

	return records, nil
}
