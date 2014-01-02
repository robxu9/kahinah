package models

type Advisory struct {
	Id uint64 `orm:"auto;pk"`

	Prefix     string
	AdvisoryId uint64

	Creator     *User     `orm:"rel(fk)"`
	Description string    `orm:"type(text)"`
	Issued      time.Time `orm:"auto_now_add"`
	Updated     time.Time `orm:"auto_now"`

	Updates []*BuildList `orm:"reverse(many)"`
}
