package controllers

import (
	"gopkg.in/robxu9/kahinah.v3/models"
)

type UserController struct {
	BaseController
}

func (u *UserController) Get() {
	userStr := models.IsLoggedIn(&u.Controller)
	if userStr == "" {
		u.Abort("403")
	}

	user := models.FindUser(userStr)
	_ = user

	u.TplName = ""
}
