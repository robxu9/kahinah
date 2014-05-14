package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	ADVISORY_UP         = "+"
	ADVISORY_DOWN       = "-"
	ADVISORY_MAINTAINER = "*"
	ADVISORY_BLOCK      = "v"
	ADVISORY_NONE       = "_"
)

type Advisory struct {
	Id uint64 `orm:"auto;pk"`

	Prefix     string
	AdvisoryId uint64

	Creator     *User  `orm:"rel(fk)"`
	Description string `orm:"type(text)"`

	Type string // type of advisory (bugfix, security, etc)

	BugsFixed string `orm:"type(text)"` // format: 123;234;534;123;133

	Requested time.Time `orm:"auto_now_add"`
	Issued    time.Time // will not be filled until advisory is approved
	Updated   time.Time `orm:"auto_now"`

	Updates []*BuildList `orm:"reverse(many)"`

	// lists comments
	Comments []*AdvisoryKarma `orm:"reverse(many)"`
}

type AdvisoryKarma struct {
	Id       uint64    `orm:"auto;pk"`
	Advisory *Advisory `orm:"rel(fk)"`
	User     *User     `orm:"rel(fk)"`
	Vote     string
	Comment  string `orm:"type(text)"`
}

func NextAdvisoryId(prefix string) uint64 {
	o := orm.NewOrm()

	var adv Advisory
	err := o.QueryTable(new(Advisory)).Filter("Prefix", prefix).OrderBy("-AdvisoryId").Limit(1).One(&adv)
	if err == orm.ErrNoRows {
		return 1
	} else if err != nil {
		panic(err)
	}

	return adv.AdvisoryId + 1
}
