package models

import "time"

const (
	KARMA_UP         = "+"
	KARMA_DOWN       = "-"
	KARMA_MAINTAINER = "*"
	KARMA_BLOCK      = "v"
	KARMA_PUSH       = "^"
	KARMA_NONE       = "_"

	KARMA_ACCEPTED = "O"
	KARMA_REJECTED = "X"
)

type Karma struct {
	Id      uint64     `orm:"auto;pk" json:"-"`
	List    *BuildList `orm:"rel(fk)" json:",omitempty"`
	User    *User      `orm:"rel(fk)" json:",omitempty"`
	Vote    string
	Comment string    `orm:"type(text)"`
	Time    time.Time `orm:"auto_now;type(datetime)"`
}
