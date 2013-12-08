package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/robxu9/kahinah/models"
	"log"
	"sort"
)

type ByTimeTP []models.BuildList

func (b ByTimeTP) Len() int {
	return len(b)
}
func (b ByTimeTP) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b ByTimeTP) Less(i, j int) bool {
	return b[i].BuildDate.Unix() > b[j].BuildDate.Unix()
}

type TestingController struct {
	beego.Controller
}

func (this *TestingController) Get() {
	var packages []models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	num, err := qt.Filter("status", models.STATUS_TESTING).All(&packages)
	if err != nil && err != orm.ErrNoRows {
		log.Println(err)
		this.Abort("500")
	}

	sort.Sort(ByTimeTP(packages))

	this.Data["Title"] = "Testing"
	this.Data["Tab"] = 1
	this.Data["Packages"] = packages
	this.Data["Entries"] = num
	this.TplNames = "generic_list.tpl"
}
