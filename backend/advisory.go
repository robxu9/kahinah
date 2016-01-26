package kahinah

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	rt "gopkg.in/dancannon/gorethink.v1"
)

const (
	// TableAdvisory represents the advisory table. See the Advisory struct
	// for the layout.
	TableAdvisory = "advisories"

	// TableCounter represents the counter table. It is structured like so:
	// "group": "GROUP" # group is the primary key
	// "YEAR": counter # counter is an int
	TableCounter = "counter"
)

// Advisory represents a group of updates that should be evaluated together.
// They present a unified "patch-set" that users can be informed from, as well
// as a unified interface for testers to identify common bugs and provide
// feedback for the system.
type Advisory struct {
	// rethinkdb unique id
	ID string `gorethink:"id,omitempty"`

	// advisory attribute (indexes)
	// these form the compound index 'advisory_id'
	Group       string `gorethink:"group"`
	Year        int    `gorethink:"year"`
	AdvisoryNum int    `gorethink:"advisory_num"`

	// Ruleset dictates what category this advisory falls into.
	Ruleset string `gorethink:"ruleset"`

	// advisories have revisions that are attached to them. this is so that
	// changes can be tracked as well as reverted to as necessary.
	Revisions map[time.Time]*Revision `gorethink:"revisions"`

	// Verdicts and Comments are the same thing, really. (Comment is +- 0)
	Verdicts []*Verdict `gorethink:"verdicts"`

	// Pushed means this advisory was pushed, Deprecated means that you really
	// should loo
	Pushed     bool `gorethink:"pushed"`
	Deprecated bool `gorethink:"deprecated"`
}

// Revision represents a specific revision of an advisory. It can be considered
// a snapshot of an advisory, and updates to advisories usually involve creating
// a new revision instead.
type Revision struct {
	Description string   `gorethink:"description"`
	References  []string `gorethink:"references"`

	Updates []string `gorethink:"update_ids"` // uuids of updates
}

// Verdict represents an opinion on an advisory. A comment is a verdict with
// no influence.
type Verdict struct {
	Time      time.Time `gorethink:"created_at"`
	Submitter string    `gorethink:"submitter"` // this is system dependent so we just expose a string

	Type      string `gorethink:"type"`
	Influence int    `gorethink:"influence"`
}

type AdvisoryFilter struct {
}

func (k *K) GetAdvisories(filter AdvisoryFilter) []*Advisory {

}

// GetNextAdvisoryNumber retrieves the next advisory number for the group and
// year from the database. It uses the counters table, which keeps track of the
// latest numbers and atomically sets to increment and get the next one.
func (k *K) GetNextAdvisoryNumber(group string, year int) int {
	strYear := strconv.Itoa(year)

	result, err := k.RunWrite(k.DB().Table(TableCounter, rt.TableOpts{}).Get(group).Update(map[string]interface{}{
		strYear: rt.Row.Field(strYear).Add(1),
	}, rt.UpdateOpts{
		ReturnChanges: true,
	}))

	if err != nil { // we traditionally don't get errors?
		panic(err)
	}

	if result.Skipped != 0 { // case: group doesn't exist
		_, err = k.RunWrite(k.DB().Table(TableCounter).Insert(map[string]interface{}{
			"group": group,
			strYear: 0,
		}))
		// Something went wrong? Panic.
		if err != nil {
			panic(err)
		}

		return k.GetNextAdvisoryNumber(group, year) // recurse to try again
	}

	if result.Errors > 0 {
		if strings.Contains(result.FirstError, "No attribute") {
			// no such year, so let's add it
			//
			// made from this confirmed reql query:
			// r.db('test').table('counters').get('omv').do(
			//   r.branch(
			//     r.row.hasFields('2014'),
			//     r.object(),
			//     r.db('test').table('counters').get('omv').update({2014: 0}, {returnChanges: true})
			//   )
			// )
			_, err := k.Run(k.DB().Table(TableCounter).Get(group).Do(
				rt.Branch(
					rt.Row.HasFields(strYear),
					rt.Object(),
					k.DB().Table(TableCounter).Get(group).Update(map[string]interface{}{
						strYear: 0,
					}))))

			if err != nil {
				panic(err)
			}

			return k.GetNextAdvisoryNumber(group, year)
		}

		panic(fmt.Errorf("error from rethinkdb that we don't know: %v", result.FirstError))
	}

	if len(result.Changes) != 1 {
		panic("should not have gotten more than one change!?")
	}

	newValues := result.Changes[0].NewValue.(map[string]interface{})

	return newValues[strYear].(int)
}
