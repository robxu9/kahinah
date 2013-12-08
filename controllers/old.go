package controllers

import (
	"github.com/astaxie/beego"
)

type OldMainController struct {
	beego.Controller
}

func (this *OldMainController) Get() {
	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie@gmail.com"
	this.TplNames = "index.tpl"
}
