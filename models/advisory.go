package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

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
