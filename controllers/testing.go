package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/robxu9/kahinah/models"
	"log"
	"sort"
)

type ByBuildDate []*models.BuildList

func (b ByBuildDate) Len() int {
	return len(b)
}
func (b ByBuildDate) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b ByBuildDate) Less(i, j int) bool {
	return b[i].BuildDate.Unix() > b[j].BuildDate.Unix()
}

type TestingController struct {
	beego.Controller
}

func (this *TestingController) Get() {
	this.Data["xsrf_token"] = this.XsrfToken()

	var packages []*models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	num, err := qt.Filter("status", models.STATUS_TESTING).All(&packages)
	if err != nil && err != orm.ErrNoRows {
		log.Println(err)
		this.Abort("500")
	}

	for _, v := range packages {
		o.LoadRelated(v, "Submitter")
	}

	sort.Sort(ByBuildDate(packages))

	this.Data["Title"] = "Testing"
	this.Data["Tab"] = 1
	this.Data["Packages"] = packages
	this.Data["Entries"] = num
	this.TplNames = "generic_list.tpl"
}
