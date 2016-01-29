package controllers

import "github.com/robxu9/kahinah/models"

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
