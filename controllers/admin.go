package controllers

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/astaxie/beego/orm"
	"github.com/robxu9/kahinah/conf"
	"github.com/robxu9/kahinah/data"
	"github.com/robxu9/kahinah/models"
	"github.com/robxu9/kahinah/sessions"
	"github.com/robxu9/kahinah/util"
)

var (
	adminWhitelist  = conf.Config.Get("admin.administrators").([]interface{})
	enableWhitelist = conf.Config.Get("admin.whitelist").(bool)

	adminWhitelistSet *util.Set
)

func init() {
	adminWhitelistSet = util.NewSet()
	for _, v := range adminWhitelist {
		adminWhitelistSet.Add(v.(string))
	}

	models.PermRegister(PermissionAdmin)
	models.PermRegister(PermissionQA)
	models.PermRegister(PermissionWhitelist)
}

func adminCheck(r *http.Request) {
	user := MustAuthenticate(r)

	if adminWhitelistSet.Contains(user) {
		return
	}

	PermAbortCheck(r, PermissionAdmin)
}

// AdminGetHandler controls the central dashboard for Kahinah.
func AdminGetHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	adminCheck(r)

	dataRenderer := data.FromContext(ctx)

	toRender := map[string]interface{}{
		"Loc":   0,
		"Tab":   -1,
		"Title": "Admin",
	}

	user := r.FormValue("username")
	if user != "" {
		userModel := models.FindUser(user)
		toRender["User"] = userModel
	}

	perms := map[string][]string{}
	for _, perm := range models.PermGetAll() {
		perms[perm.Permission] = []string{}
		for _, u := range perm.Users {
			perms[perm.Permission] = append(perms[perm.Permission], u.Username)
		}
	}
	toRender["Permissions"] = perms

	dataRenderer.Data = toRender
	dataRenderer.Template = "admin"
}

// AdminPostHandler manipulates the central dashboard for kahinah
func AdminPostHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	adminCheck(r)

	sess := sessions.FromContext(ctx)

	user := r.FormValue("username")
	if user != "" {
		model := models.FindUser(user)
		o := orm.NewOrm()
		m2m := o.QueryM2M(model, "Permissions")

		add := r.FormValue("add")
		rm := r.FormValue("rm")

		if add != "" {
			addpermobj := models.PermGet(add)
			if addpermobj == nil {
				sess.AddFlash(sessions.FlashError, "No such permission "+add+"!")
			} else {
				if !m2m.Exist(addpermobj) {
					_, err := m2m.Add(addpermobj)
					if err != nil {
						panic(err)
					}
				}
			}
		}

		if rm != "" {
			rmpermobj := models.PermGet(rm)
			if rmpermobj == nil {
				sess.AddFlash(sessions.FlashError, "No such permission "+rm+"!")
			} else {
				_, err := m2m.Remove(rmpermobj)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	http.Redirect(rw, r, util.GetPrefixString("/admin"), http.StatusTemporaryRedirect)
}
