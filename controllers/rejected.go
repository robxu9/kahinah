package controllers

import (
	"log"
	"sort"

	"github.com/astaxie/beego/orm"
	"github.com/robxu9/kahinah/models"
)

type RejectedController struct {
	BaseController
}

func (this *RejectedController) Get() {
	page, err := this.GetInt("page")
	if err != nil {
		page = 1
	} else if page <= 0 {
		page = 1
	}

	var packages []*models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	cnt, err := qt.Filter("status", models.STATUS_REJECTED).Count()
	if err != nil {
		log.Println(err)
		this.Abort("500")
	}

	totalpages := cnt / 50
	if cnt%50 != 0 {
		totalpages++
	}

	if page > totalpages {
		page = totalpages
	}

	_, err = qt.Limit(50, (page-1)*50).OrderBy("-Updated").Filter("status", models.STATUS_REJECTED).All(&packages)
	if err != nil && err != orm.ErrNoRows {
		log.Println(err)
		this.Abort("500")
	}

	for _, v := range packages {
		o.LoadRelated(v, "Submitter")
	}

	sort.Sort(ByUpdateDate(packages))

	this.Data["Title"] = "Rejected"
	this.Data["Loc"] = 1
	this.Data["Tab"] = 3
	this.Data["Packages"] = packages
	this.Data["PrevPage"] = page - 1
	this.Data["Page"] = page
	this.Data["NextPage"] = page + 1
	this.Data["Pages"] = totalpages
	this.TplNames = "builds_list.tpl"
}
