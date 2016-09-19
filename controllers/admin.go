package controllers

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/jinzhu/gorm"
	"github.com/robxu9/kahinah/common/conf"
	"github.com/robxu9/kahinah/common/set"
	"github.com/robxu9/kahinah/data"
	"github.com/robxu9/kahinah/models"
	"github.com/robxu9/kahinah/render"
)

var (
	adminWhitelist  = conf.Get("admin.administrators").([]interface{})
	enableWhitelist = conf.Get("admin.whitelist").(bool)

	adminWhitelistSet *set.Set
)

func init() {
	adminWhitelistSet = set.NewSet()
	for _, v := range adminWhitelist {
		adminWhitelistSet.Add(v.(string))
	}
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

	var perms []models.UserPermission
	if err := models.DB.Find(&perms).Error; err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	rendered := map[string][]string{}

	for _, perm := range perms {
		if _, ok := rendered[perm.Permission]; !ok {
			rendered[perm.Permission] = []string{}
		}
		rendered[perm.Permission] = append(rendered[perm.Permission], models.FindUserByID(perm.UserID).Username)
	}

	toRender["Permissions"] = rendered

	dataRenderer.Data = toRender
	dataRenderer.Template = "admin"
}

// AdminPostHandler manipulates the central dashboard for kahinah
func AdminPostHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	adminCheck(r)

	user := r.FormValue("username")
	action := r.FormValue("action")
	permission := r.FormValue("permission")

	if user == "" || (action != "add" && action != "rm") || permission == "" {
		panic(ErrBadRequest)
	}

	modelUser := models.FindUser(user)

	if action == "add" {
		if err := models.DB.Model(modelUser).Association("Permissions").Append(models.UserPermission{
			Permission: permission,
		}).Error; err != nil {
			panic(err)
		}
	} else {
		if err := models.DB.Model(modelUser).Association("Permissions").Delete(models.UserPermission{
			Permission: permission,
		}).Error; err != nil {
			panic(err)
		}
	}

	http.Redirect(rw, r, render.ConvertURL("/admin"), http.StatusTemporaryRedirect)
}
