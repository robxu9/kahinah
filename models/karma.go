package models

import "time"

const (
	KARMA_UP         = "+"
	KARMA_DOWN       = "-"
	KARMA_MAINTAINER = "*"
	KARMA_BLOCK      = "v"
	KARMA_PUSH       = "^"
	KARMA_NONE       = "_"
)

type Karma struct {
	Id      uint64     `orm:"auto;pk"`
	List    *BuildList `orm:"rel(fk)"`
	User    *User      `orm:"rel(fk)"`
	Vote    string
	Comment string    `orm:"type(text)"`
	Time    time.Time `orm:"auto_now;type(datetime)"`
}
