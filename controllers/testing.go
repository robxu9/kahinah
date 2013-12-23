package controllers

import (
	"github.com/astaxie/beego/orm"
	"github.com/robxu9/kahinah/models"
	"log"
	"menteslibres.net/gosexy/to"
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
	BaseController
}

func (this *TestingController) Get() {
	var packages []*models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	num, err := qt.Filter("status", models.STATUS_TESTING).All(&packages)
	if err != nil && err != orm.ErrNoRows {
		log.Println(err)
		this.Abort("500")
	}

	pkgkarma := make(map[string]string)

	for _, v := range packages {
		o.LoadRelated(v, "Submitter")
		o.LoadRelated(v, "Karma")

		totalKarma := 0

		for _, karma := range v.Karma {
			if karma.Vote == models.KARMA_UP {
				totalKarma++
			} else if karma.Vote == models.KARMA_DOWN {
				totalKarma--
			} else if karma.Vote == models.KARMA_MAINTAINER {
				totalKarma += int(maintainer_karma)
			} else if karma.Vote == models.KARMA_BLOCK {
				totalKarma -= int(block_karma)
			}
		}

		pkgkarma[to.String(v.Id)] = to.String(totalKarma)
	}

	sort.Sort(ByBuildDate(packages))

	this.Data["Title"] = "Testing"
	this.Data["Tab"] = 1
	this.Data["Packages"] = packages
	this.Data["PkgKarma"] = pkgkarma
	this.Data["Entries"] = num
	this.TplNames = "generic_list.tpl"
}
