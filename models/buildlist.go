package models

import (
	"time"
)

const (
	StatusTesting   = "testing"
	StatusRejected  = "rejected"
	StatusPublished = "published"

	LinkMainURL  = "_mainURL"
	ChangelogURL = "_changelogURL"
	SCMLogURL    = "_scmlogURL"
)

type BuildList struct {
	// represent unique id
	Id uint64 `xml:"id,attr" orm:"auto;pk"`
	// represent the platform it came from (e.g. RHEL7)
	Platform string `xml:"info>platform"`
	// represent the sub-platform it came from (if any) (e.g. stable/testing/dev)
	Repo string `xml:"info>platform>repo"`
	// represent the variations (a.k.a. architectures) (e.g. x86_64/arm64)
	Architecture string `xml:"info>platform>arch"`
	// represent its name (e.g. project name like "Mantra")
	Name string `xml:"info>name"`
	// represents the submitter
	Submitter *User `xml:"info>submitter" orm:"rel(fk)"`
	// represents the update type (bugfix, new package, security, etc)
	Type string `xml:"type"`
	// represents the current status of this update
	Status string `xml:"status"`
	// lists packages
	Packages []*BuildListPkg `xml:"packages" orm:"reverse(many)"`
	// lists links
	Links []*BuildListLink `xml:"links" orm:"reverse(many)"`
	// lists the builddate
	BuildDate time.Time `xml:"time"`
	// lists the updatetime
	Updated time.Time `xml:"updated" orm:"auto_now"`
	// lists the karma (or recent updates)
	Karma []*Karma `xml:"karma" orm:"reverse(many)"`
	// shows the diff (defined by the build system)
	Diff string `xml:"diff" orm:"type(text)"`
	// shows the link to an advisory (when published)
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

type BuildListLink struct {
	Id   uint64     `orm:"auto;pk"`
	List *BuildList `orm:"rel(fk)"`
	Name string
	Url  string
}
