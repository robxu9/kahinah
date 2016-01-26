package controllers

import (
	"strings"

	"github.com/astaxie/beego/orm"
	"gopkg.in/robxu9/kahinah.v3/models"
	"menteslibres.net/gosexy/to"
)

type SearchController struct {
	BaseController
	Filters    map[string]interface{}
	Parameters map[interface{}]interface{}
	Generic    bool
}

func (this *SearchController) Get() {
	this.Data = this.Parameters
	if this.Generic {
		this.TplName = "generic_list.tpl"
	} else {
		this.TplName = "builds_list.tpl"
	}
}

// call with /path/to/search/:{page+1}?{filterList:filter}&{sortList:sort}&size={size}
func (this *SearchController) JsonGet() {
	page := to.Uint64(this.Ctx.Input.Param(":page"))
	_ = page

	var packages []*models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	for k, v := range this.Filters {
		qt = qt.Filter(k, v)
	}

	if this.Generic {
		// if we're doing a generic search it better not be too large,
		// since I don't plan on making any limitations on this search
		// [usable e.g. for testing only... definitely not for published/rejected]
		// ...at least unless I fix karma
		// (FIXME)

		// filter[0]=Name/Architecture
		// filter[1]=Submitter
		// filter[2]=For (Platform/Repo)
		// filter[3]=Type
		// filter[4]=Karma
		// filter[5]=BuildDate
		if filter := this.GetString("filter[0]"); filter != "" {
			if strings.Contains(filter, "/") {
				strs := strings.Split(filter, "/")
				qt = qt.Filter("Name__icontains", strs[0])
				qt = qt.Filter("Architecture__icontains", strs[1])
			} else { // could be either one, check both
				cond := orm.NewCondition()
				qt = qt.SetCond(cond.AndCond(cond.Or("Name__icontains", filter).Or("Architecture__icontains", filter)))
			}
		}
		if filter := this.GetString("filter[1]"); filter != "" {
			qt = qt.Filter("Submitter__Email__icontains", filter)
		}
		if filter := this.GetString("filter[2]"); filter != "" {
			if strings.Contains(filter, "/") {
				strs := strings.Split(filter, "/")
				qt = qt.Filter("Platform__icontains", strs[0])
				qt = qt.Filter("Repo__icontains", strs[1])
			} else { // could be either one, check both
				cond := orm.NewCondition()
				qt = qt.SetCond(cond.AndCond(cond.Or("Platform__icontains", filter).Or("Repo__icontains", filter)))
			}
		}
		if filter := this.GetString("filter[3]"); filter != "" {
			qt = qt.Filter("Type__icontains", filter)
		}
		if filter := this.GetString("filter[4]"); filter != "" {
			// KARMA ISSUES (RelatedSel)
			// Use with http://beego.me/docs/mvc/model/query.md#relatedsel but it seems incomplete?
			// TODO investigate

			qt = qt.Filter("Submitter__Email__icontains", filter)
		}
		if filter := this.GetString("filter[5]"); filter != "" {
			qt = qt.Filter("BuildDate__icontains", filter)
		}
	} else {
		// FIXME
		// filter[0]=updateid
		// filter[1]=name
		// filter[2]=submitter
		// filter[3]=for
		// filter[4]=type
		// filter[5]=status
		// filter[6]=updated
		if filter := this.GetString("filter[0]"); filter != "" {
			if strings.Contains(filter, "-") {

			} else {

			}
		}
		if filter := this.GetString("filter[1]"); filter != "" {
			if strings.Contains(filter, "/") {
				strs := strings.Split(filter, "/")
				qt = qt.Filter("Name__icontains", strs[0])
				qt = qt.Filter("Architecture__icontains", strs[1])
			} else { // could be either one, check both
				cond := orm.NewCondition()
				qt = qt.SetCond(cond.AndCond(cond.Or("Name__icontains", filter).Or("Architecture__icontains", filter)))
			}
		}
		if filter := this.GetString("filter[1]"); filter != "" {
			qt = qt.Filter("Submitter__Email__icontains", filter)
		}
		if filter := this.GetString("filter[2]"); filter != "" {
			if strings.Contains(filter, "/") {
				strs := strings.Split(filter, "/")
				qt = qt.Filter("Platform__icontains", strs[0])
				qt = qt.Filter("Repo__icontains", strs[1])
			} else { // could be either one, check both
				cond := orm.NewCondition()
				qt = qt.SetCond(cond.AndCond(cond.Or("Platform__icontains", filter).Or("Repo__icontains", filter)))
			}
		}
		if filter := this.GetString("filter[3]"); filter != "" {
			qt = qt.Filter("Type__icontains", filter)
		}
		if filter := this.GetString("filter[4]"); filter != "" {
			// KARMA ISSUES (RelatedSel)
			// Use with http://beego.me/docs/mvc/model/query.md#relatedsel but how to do such thing T-T
			qt = qt.Filter("Submitter__Email__icontains", filter)
		}
		if filter := this.GetString("filter[5]"); filter != "" {
			qt = qt.Filter("BuildDate__icontains", filter)
		}
	}

	qt.All(&packages)
}
