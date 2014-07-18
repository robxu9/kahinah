package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/astaxie/beego"
	"menteslibres.net/gosexy/to"
)

var (
	outwardUrl = beego.AppConfig.String("outwardloc")
)

type PersonaResponse struct {
	Status   string `json: "status"`
	Email    string `json: "email,omitempty"`
	Audience string `json: "audience,omitempty"`
	Expires  int64  `json: "expires,omitempty"`
	Issuer   string `json: "issuer,omitempty"`
	Reason   string `json: "reason,omitempty"`
}

type PersonaCheckController struct {
	beego.Controller
}

func (this *PersonaCheckController) Get() {
	session := this.GetSession("persona")
	if session != nil {
		pr := PersonaResponse{}
		json.Unmarshal(to.Bytes(session), &pr)
		this.Ctx.WriteString(pr.Email)
	} else {
		this.Ctx.WriteString("")
	}
}

type PersonaLogoutController struct {
	beego.Controller
}

func (this *PersonaLogoutController) Get() {
	this.DelSession("persona")
	this.DestroySession()
	this.Ctx.WriteString("OK")
}

type PersonaLoginController struct {
	beego.Controller
}

func (this *PersonaLoginController) Post() {
	assertion := this.GetString("assertion")
	if assertion == "" {
		this.Abort("400")
	}

	data := url.Values{"assertion": {assertion}, "audience": {outwardUrl}}

	resp, err := http.PostForm("https://verifier.login.persona.org/verify", data)
	if err != nil {
		log.Println(err)
		this.Abort("400")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		this.Abort("500")
	}

	pr := PersonaResponse{}
	err = json.Unmarshal(body, &pr)
	if err != nil {
		log.Println(err)
		this.Abort("403")
	}

	if pr.Status != "okay" {
		this.Abort("403")
	}

	go FindUser(pr.Email)

	this.SetSession("persona", body)

	this.Data["json"] = &pr
	this.ServeJson()
}
