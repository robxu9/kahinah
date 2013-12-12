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
	Id           uint64          `xml:"id,attr" orm:"auto;pk"`
	Platform     string          `xml:"info>platform"`
	Repo         string          `xml:"info>platform>repo"`
	Architecture string          `xml:"info>platform>arch"`
	Name         string          `xml:"info>name"`
	Submitter    *User           `xml:"info>submitter" orm:"rel(fk)"`
	Type         string          `xml:"type"`
	Status       string          `xml:"status"`
	Packages     []*BuildListPkg `xml:"packages" orm:"reverse(many)"` // semicolon-separated (damn orm)
	Changelog    string          `xml:"changelog" orm:"type(text)"`   // url
	BuildDate    time.Time       `xml:"time"`
	Updated      time.Time       `xml:"updated" orm:"auto_now"`
	Karma        []*Karma        `xml:"karma" orm:"reverse(many)"`

	// handler specifics
	Handler  string `xml:"info>handle,attr"` // what handler should this use?
	HandleId string `xml:"info>id,attr"`     // for the handler to identify the package in the buildsystem
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
	Url     string `orm:"type(text)"`
}
