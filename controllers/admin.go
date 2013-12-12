package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/robxu9/kahinah/models"
	"github.com/robxu9/kahinah/util"
	"menteslibres.net/gosexy/to"
	"strings"
)

const (
	PERMISSION_ADMIN     = "kahinah.admin"
	PERMISSION_QA        = "kahinah.qa"
	PERMISSION_WHITELIST = "kahinah.whitelist"
)

var (
	adminWhitelist = strings.Split(beego.AppConfig.String("adminwhitelist"), ";")
	Whitelist      = to.Bool(beego.AppConfig.String("whitelist"))
)

func init() {
	models.PermRegister(PERMISSION_ADMIN)
	models.PermRegister(PERMISSION_QA)
	models.PermRegister(PERMISSION_WHITELIST)
}

func adminCheck(this *beego.Controller) {
	user := util.IsLoggedIn(this)

	if user == "" {
		this.Abort("403")
	}

	loggedin := false

	for _, v := range adminWhitelist {
		if v == user {
			loggedin = true
		}
	}

	if !loggedin {
		models.PermAbortCheck(this, PERMISSION_ADMIN)
	}
}

type AdminController struct {
	beego.Controller
}

func (this *AdminController) Get() {
	adminCheck(&this.Controller)

	if this.GetString("email") != "" {
		user := models.FindUserNoCreate(this.GetString("email"))
		if user == nil {
			flash := beego.NewFlash()
			flash.Error("No such email " + this.GetString("email") + "!")
			flash.Store(&this.Controller)
		} else {
			o := orm.NewOrm()
			m2m := o.QueryM2M(user, "Permissions")
			this.Data["User"] = user
			if this.GetString("add") != "" {
				addperm := this.GetString("add")
				addpermobj := models.PermGet(addperm)
				if addpermobj == nil {
					flash := beego.NewFlash()
					flash.Error("No such permission " + addperm + "!")
					flash.Store(&this.Controller)
				} else {
					if !m2m.Exist(addpermobj) {
						_, err := m2m.Add(addpermobj)
						if err != nil {
							panic(err)
						}
					}
				}
			}
			if this.GetString("rm") != "" {
				rmperm := this.GetString("rm")
				rmpermobj := models.PermGet(rmperm)
				if rmpermobj == nil {
					flash := beego.NewFlash()
					flash.Error("No such permission " + rmperm + "!")
					flash.Store(&this.Controller)
				} else {
					_, err := m2m.Remove(rmpermobj)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}

	this.Data["Tab"] = -1

	perms := make(map[string][]string)

	for _, perm := range models.PermGetAll() {
		perms[perm.Permission] = make([]string, 0)

		for _, v := range perm.Users {
			perms[perm.Permission] = append(perms[perm.Permission], v.Email)
		}
	}

	this.Data["Title"] = "Admin"
	this.Data["Permissions"] = perms
	this.TplNames = "admin.tpl"
}
