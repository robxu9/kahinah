package models

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"menteslibres.net/gosexy/to"
)

func IsLoggedIn(controller *beego.Controller) string {
	// check api key header first
	apiKey := controller.Ctx.Input.Header("X-Kahinah-Key")
	if apiKey != "" {
		user := FindUserApi(apiKey)
		if user != nil {
			return user.Email
		}
	}

	// check persona
	session := controller.GetSession("persona")
	if session == nil {
		return ""
	}
	pr := PersonaResponse{}
	json.Unmarshal(to.Bytes(session), &pr)
	return pr.Email
}
