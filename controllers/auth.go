package controllers

import (
	"net/http"

	"github.com/robxu9/kahinah/models"

	"gopkg.in/cas.v1"
)

const (
	PermissionAdmin     = "kahinah.admin"
	PermissionQA        = "kahinah.qa"
	PermissionWhitelist = "kahinah.whitelist"
)

func Authenticated(r *http.Request) string {
	apiKey := r.Header.Get("X-Kahinah-Key")
	if apiKey != "" {
		user := models.FindUserByAPI(apiKey)
		if user != nil {
			return user.Username
		}
	}

	if cas.IsAuthenticated(r) {
		return cas.Username(r)
	}

	return ""
}

func MustAuthenticate(r *http.Request) string {
	if u := Authenticated(r); u != "" {
		return u
	}
	panic(ErrForbidden)
}

func PermAbortCheck(r *http.Request, perm string) {
	user := Authenticated(r)
	if user != "" {

		model := models.FindUser(user)
		for _, v := range model.Permissions {
			if v.Permission == perm {
				return
			}
		}

	}

	panic(ErrPermission)
}

func PermCheck(r *http.Request, perm string) bool {
	user := Authenticated(r)
	if user != "" {

		model := models.FindUser(user)
		for _, v := range model.Permissions {
			if v.Permission == perm {
				return true
			}
		}

	}
	return false
}
