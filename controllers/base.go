package controllers

import (
	"html/template"
	"time"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) Prepare() {
	this.Data["xsrf_token"] = this.XSRFToken()
	this.Data["xsrf_data"] = template.HTML(this.XSRFFormHTML())

	this.Data["copyright"] = time.Now().Year()
}
