package models

import (
	"time"
)

const (
	STATUS_TESTING   = "testing"
	STATUS_REJECTED  = "rejected"
	STATUS_PUBLISHED = "published"
)

type BuildList struct {
	// represent unique id
	Id uint64 `xml:"id,attr" orm:"auto;pk"`
	// represent platform
	Platform string `xml:"info>platform"`
	// represent the repo it's being saved to
	Repo string `xml:"info>platform>repo"`
	// represent the arch
	Architecture string `xml:"info>platform>arch"`
	// represent its name
	Name string `xml:"info>name"`
	// represents the user
	Submitter *User `xml:"info>submitter" orm:"rel(fk)"`
	// represents the update type (bugfix, new package, security, etc)
	Type string `xml:"type"`
	// represents the current status of this update
	Status string `xml:"status"`
	// lists packages
	Packages []*BuildListPkg `xml:"packages" orm:"reverse(many)"`
	// lists the changelog url
	Changelog string `xml:"changelog" orm:"type(text)"` // url
	// lists the builddate
	BuildDate time.Time `xml:"time"`
	// lists the updatetime
	Updated time.Time `xml:"updated" orm:"auto_now"`
	// lists the karma
	Karma []*Karma `xml:"karma" orm:"reverse(many)"`
	// shows the diff
	Diff string `xml:"diff" orm:"type(text)"`
	// shows the link to an advisory, if any
	Advisory *Advisory `xml:"advisory" orm:"null;rel(fk);on_delete(set_null)"`

	// abf specifics (abf is represented as the handler)
	HandleId       string `xml:"handle>id,attr"` // for the handler to identify the package in the buildsystem
	HandleProject  string `xml:"handle>project" orm:"type(text)"`
	HandleCommitId string `xml:"handle>commitid" orm:"type(text)"`
	// end handler specifics
}

type BuildListPkg struct {
	Id      uint64     `orm:"auto;pk"`
	List    *BuildList `orm:"rel(fk)"`
	Name    string
	Type    string
	Epoch   int64
	Version string
	Release string
	Arch    string
	Url     string `orm:"type(text)"`
}
