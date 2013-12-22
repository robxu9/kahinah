package controllers

import (
	"github.com/astaxie/beego"
	"html/template"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) Prepare() {
	this.Data["xsrf_token"] = this.XsrfToken()
	this.Data["xsrf_data"] = template.HTML(this.XsrfFormHtml())
}
