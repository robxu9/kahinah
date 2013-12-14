package models

const (
	KARMA_UP         = "+"
	KARMA_DOWN       = "-"
	KARMA_MAINTAINER = "*"
	KARMA_BLOCK      = "v"
)

type Karma struct {
	Id   uint64     `orm:"auto;pk"`
	List *BuildList `orm:"rel(fk)"`
	User *User      `orm:"rel(fk)"`
	Vote string
}
