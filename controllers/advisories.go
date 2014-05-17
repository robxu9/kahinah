package controllers

import (
	"log"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/robxu9/kahinah/models"
)

var (
	enabledPlatforms = make(map[string]string) // [platform]PREFIX
)

func init() {
	configPlatforms := strings.Split(beego.AppConfig.String("advisory::platforms"), ";")
	for _, v := range configPlatforms {
		parts := strings.Split(v, ":")
		enabledPlatforms[parts[0]] = parts[1]
	}
}

//
// our base controller
//

type AdvisoryBaseController struct {
	BaseController
}

func (this *AdvisoryBaseController) Prepare() {
	this.BaseController.Prepare()
	this.Data["Loc"] = 2
}

//
// main controller
// shows recent advisories for enabled platforms
//

type AdvisoryMainController struct {
	AdvisoryBaseController
}

func (this *AdvisoryMainController) Get() {
	platforms := make(map[string][]*models.Advisory)

	o := orm.NewOrm()

	for k, _ := range enabledPlatforms {
		qt := o.QueryTable(new(models.Advisory)).Filter("Platform", k).OrderBy("-Issued").Limit(5)

		var advisories []*models.Advisory

		_, err := qt.All(&advisories)
		if err != nil && err != orm.ErrNoRows {
			log.Printf("error occured trying to get advisories: %s", err)
			this.Abort("500")
		}
		platforms[k] = advisories
	}

	this.Data["Tab"] = 1
	this.Data["Title"] = "Advisories"
	this.TplNames = "advisories/main.tpl"

	this.Data["Platforms"] = platforms
}
