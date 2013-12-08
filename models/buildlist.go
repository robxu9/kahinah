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
	Id            uint64    `xml:"id,attr" orm:"auto;pk"`
	ListId        uint64    `xml:"info>id,attr"`
	Platform      string    `xml:"info>platform"`
	Repo          string    `xml:"info>platform>repo"`
	Architecture  string    `xml:"info>platform>arch"`
	Name          string    `xml:"info>name"`
	Submitter     string    `xml:"info>submitter"`
	Type          string    `xml:"type"`
	Status        string    `xml:"status"`
	Url           string    `xml:"url" orm:"type(text)"`
	Packages      string    `xml:"packages" orm:"type(text)"`  // semicolon-separated (damn orm)
	Changelog     string    `xml:"changelog" orm:"type(text)"` // url
	PublishHandle string    `xml:"-"`
	RejectHandle  string    `xml:"-"`
	BuildDate     time.Time `xml:"time"`
	Updated       time.Time `xml:"updated" orm:"auto_now"`
}
