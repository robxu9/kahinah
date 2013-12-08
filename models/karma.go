package models

const (
	KARMA_UP   = "+"
	KARMA_DOWN = "-"
)

type Karma struct {
	Id     uint64 `orm:"auto;pk"`
	ListId uint64
	User   string `orm:"type(text)"`
	Vote   string
}
