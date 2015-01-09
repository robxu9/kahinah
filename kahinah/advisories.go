package kahinah

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

// AdvisoryStatus reflects the current status of the advisory
type AdvisoryStatus int

const (
	// OPEN represents an advisory that is open for discussion
	OPEN AdvisoryStatus = iota
	// PASS represents an advisory that has passed QA
	PASS
	// FAIL represents an advisory that has failed QA
	FAIL
)

var (
	// ErrNoSuchAdvisory - update does not exist.
	ErrNoSuchAdvisory = errors.New("kahinah: advisory doesn't exist")

	// ErrAdvisoryChanged - the advisory has changed
	ErrAdvisoryChanged = errors.New("kahinah: advisory has changed")

	// ErrAdvisoryUpdateAttached - an update has already been attached to an advisory
	ErrAdvisoryUpdateAttached = errors.New("kahinah: update already attached to advisory")

	// CacheAdvisoryExp - the default expiration for updates in the cache
	CacheAdvisoryExp = 1 * time.Hour
)

// Advisory represents a collection of updates that should be
// sent out together. Advisories also contain information about
// who proposed the advisory, karma [for voting], comments, and
// the proposed advisory ID.
type Advisory struct {
	Id int64

	UserId   int64      // Submitter
	Updates  []int64    `sql:"-"` // List of Updates
	Comments []*Comment `sql:"-"` // List of Comments

	References    []string `sql:"-"`
	ReferencesStr string   `sql:"references"` // \n seperated string b/c sql

	Description string // description of the advisory

	Status AdvisoryStatus // Status of the Advisory

	AdvisoryFamily string // Outward-facing Advisory Family
	AdvisoryId     int64  // Outward-facing Advisory ID

	Version int64 // versioning

	CreatedAt time.Time // time this was created
	UpdatedAt time.Time // time this was updated
}

// BeforeSave is solely for gorm operation
func (a *Advisory) BeforeSave() error {
	a.ReferencesStr = strings.Join(a.References, "\n")
	return nil
}

// AfterFind is solely for gorm operation
func (a *Advisory) AfterFind() error {
	a.References = strings.Split(a.ReferencesStr, "\n")
	return nil
}

// CommentVerdict represents a comment's final opinion on an advisory.
type CommentVerdict int

const (
	// NEUTRAL represents no final opinion
	NEUTRAL CommentVerdict = iota
	// NO represents a negative opinion
	NO
	// YES represents a positive opinion
	YES
	// BLOCK represents a final extremely negative opinion
	BLOCK
	// OVERRIDE represents a final positive opinion
	OVERRIDE
)

// Comment represents an opinion on an Advisory.
type Comment struct {
	Id int64

	UserId     int64 // Submitter
	AdvisoryId int64 // Advisory this was attached to

	Verdict  CommentVerdict // verdict
	Thoughts string         // thoughts

	CreatedAt time.Time // time this was created
}

func advisoryCacheId(id int64) string {
	return "advisories/" + strconv.FormatInt(id, 10)
}

func (k *Kahinah) incrementVersion(a *Advisory) error {
	a.Version++
	return k.db.Save(a).Error
}

func (k *Kahinah) nextFamilyID(family string) int64 {
	var id []int64
	if err := k.db.Model(&Advisory{}).Where(&Advisory{AdvisoryFamily: family}).Order("advisory_id desc").Limit(1).Pluck("advisory_id", &id).Error; err != nil {
		panic(err)
	}

	if len(id) == 0 {
		return 1
	}

	return id[0] + 1
}

// NewAdvisory creates and records an advisory in the system.
func (k *Kahinah) NewAdvisory(user int64, updates []int64, references []string, description, family string) (int64, error) {
	// needed for nextFamilyID()
	k.advisoryFamIDmutex.Lock()
	defer k.advisoryFamIDmutex.Unlock()

	family = strings.TrimSpace(strings.ToUpper(family)) // family is always uppercase

	// double check to make sure this user exists
	if _, err := k.RetrieveUser(user); err != nil {
		return 0, err
	}

	// double check all updates to make sure that they exist
	var updateptrs []*Update
	for _, v := range updates {
		update, err := k.RetrieveUpdate(v)
		if err != nil {
			return 0, err
		}

		// check that the update isn't already attributed to an advisory
		if update.AdvisoryId != 0 {
			return 0, ErrAdvisoryUpdateAttached
		}

		updateptrs = append(updateptrs, update)
	}

	// we're good, insert this
	record := &Advisory{
		UserId:         user,
		Updates:        updates,
		References:     references,
		Description:    description,
		Status:         OPEN,
		AdvisoryFamily: family,
		AdvisoryId:     k.nextFamilyID(family),
		Version:        0,
	}

	if err := k.db.Save(record).Error; err != nil {
		return 0, err
	}

	// now update all these updates
	for _, v := range updateptrs {
		if err := k.setUpdateAdvisoryId(v, record.Id); err != nil {
			panic(err) // this should not happen...
		}
	}

	// store in the cache
	k.cache.Set(advisoryCacheId(record.Id), record, CacheAdvisoryExp)

	return record.Id, nil
}

// RetrieveAdvisory retrieves an advisory from the cache or db
func (k *Kahinah) RetrieveAdvisory(id int64) (*Advisory, error) {
	if cached, found := k.cache.Get(advisoryCacheId(id)); found {
		record := cached.(*Advisory)
		return record, nil
	}

	return k.ForceRetrieveAdvisory(id)
}

// ForceRetrieveAdvisory circumvents the cache and forcefully retrieves
// an advisory from the db
func (k *Kahinah) ForceRetrieveAdvisory(id int64) (*Advisory, error) {
	record := &Advisory{}

	if k.db.First(record, id).RecordNotFound() {
		return nil, ErrNoSuchUpdate
	}

	// FIXME GORM WHY
	// get a list of all updates
	if err := k.db.Model(&Update{}).Where(&Update{AdvisoryId: record.Id}).Pluck("id", &record.Updates).Error; err != nil {
		panic(err)
	}

	// get a list of all comments, sorted by time
	var comments []Comment
	if err := k.db.Where(&Comment{AdvisoryId: record.Id}).Order("created_at").Find(&comments).Error; err != nil {
		panic(err)
	}

	// ugh I wanted []*Comment not []Comment
	commentsptr := make([]*Comment, len(comments))
	for k, v := range comments {
		commentsptr[k] = &v
	}

	// store in the cache
	k.cache.Set(advisoryCacheId(record.Id), record, CacheAdvisoryExp)

	return record, nil
}

// CountAdvisories lists the number of advisories in the system.
func (k *Kahinah) CountAdvisories() int64 {
	// panic if unable to get count

	var count int64

	if err := k.db.Model(&Advisory{}).Count(&count).Error; err != nil {
		panic(err)
	}

	return count
}

// ListAdvisories retrusn a list of advisories in descending order
func (k *Kahinah) ListAdvisories(from, limit int64) ([]int64, error) {
	var records []int64

	if err := k.db.Model(&Advisory{}).Order("created_at desc").Limit(limit).Offset(from).Pluck("id", &records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (k *Kahinah) processAdvisory() {
	defer k.advisoryProcessRoutines.Done()

	processRoutines := &sync.WaitGroup{}

	for advisory := range k.advisoryProcessQueue {
		// check the status of the advisory - we may have already
		// processed it?
		if advisory.Status != OPEN {
			// already processed, skip
			continue
		}

		// check whether we've met the config requirements
		advisory.Status = k.AdvisoryProcessFunc(advisory)

		// save that new status
		k.incrementVersion(advisory)

		if advisory.Status != OPEN {
			// for each update
			for _, update := range advisory.Updates {
				// get updateptr
				updateptr, err := k.RetrieveUpdate(update)
				if err != nil {
					// FIXME: log?
					continue
				}

				// get the associated connector
				connector, ok := k.connectors[updateptr.Connector]
				if !ok {
					// FIXME: log?
					continue
				}

				// send to connector
				processRoutines.Add(1)
				go func(a *Advisory, u *Update) {
					defer processRoutines.Done()

					if a.Status == FAIL {
						connector.Fail(u)
					} else {
						connector.Pass(u)
					}

				}(advisory, updateptr)
			}
		}
	}

	processRoutines.Wait()
}

// NewComment adds a new comment to an advisory.
func (k *Kahinah) NewComment(a *Advisory, user int64, verdict CommentVerdict, thoughts string) error {
	// verify that the user exists
	if _, err := k.RetrieveUser(user); err != nil {
		return err
	}

	// verify that the advisory hasn't changed before you add a new comment
	forceAdvisory, err := k.ForceRetrieveAdvisory(a.Id)
	if err != nil {
		return err
	}

	if forceAdvisory.Version != a.Version {
		return ErrAdvisoryChanged
	}

	// insert this comment
	comment := &Comment{
		UserId:     user,
		AdvisoryId: a.Id,
		Verdict:    verdict,
		Thoughts:   thoughts,
	}

	if err := k.db.Save(comment).Error; err != nil {
		return err
	}

	// append to existing advisory
	a.Comments = append(a.Comments, comment)

	// up the version and save
	if err := k.incrementVersion(a); err != nil {
		return err
	}

	// send to queue
	k.advisoryProcessQueue <- a

	return nil
}
