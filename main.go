package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gopkg.in/robxu9/kahinah.v3/util"

	"github.com/astaxie/beego"
	"gopkg.in/robxu9/kahinah.v3/controllers"
	"gopkg.in/robxu9/kahinah.v3/integration"
	"gopkg.in/robxu9/kahinah.v3/models"
	"menteslibres.net/gosexy/to"
)

func main() {
	beego.BConfig.WebConfig.Session.SessionOn = true

	beego.BConfig.WebConfig.EnableXSRF = true
	beego.BConfig.WebConfig.XSRFKey = getRandomString(50)
	beego.BConfig.WebConfig.XSRFExpire = 3600

	beego.SetStaticPath(util.GetPrefixString("/static"), "static")

	beego.Router(util.GetPrefixString("/"), &controllers.MainController{})

	//
	// --------------------------------------------------------------------
	// BUILDS
	// --------------------------------------------------------------------
	//

	// testing
	beego.Router(util.GetPrefixString("/builds/testing"), &controllers.TestingController{}) // lists testing updates
	// published
	beego.Router(util.GetPrefixString("/builds/published"), &controllers.PublishedController{})
	// rejected
	beego.Router(util.GetPrefixString("/builds/rejected"), &controllers.RejectedController{}) // lists all rejected updates
	// all builds
	beego.Router(util.GetPrefixString("/builds"), &controllers.BuildsController{}) // show all testing, published, rejected (all sorted by date, linking respectively to above)

	// specific
	beego.Router(util.GetPrefixString("/builds/:id:int"), &controllers.BuildController{})

	//
	// --------------------------------------------------------------------
	// ADVISORIES
	// --------------------------------------------------------------------
	//

	// advisories
	beego.Router(util.GetPrefixString("/advisories"), &controllers.AdvisoryMainController{})
	//beego.Router(util.GetPrefixString("/advisories/:platform:string"), &controllers.AdvisoryPlatformController{})
	//beego.Router(util.GetPrefixString("/advisories/:id:int"), &controllers.AdvisoryController{})
	beego.Router(util.GetPrefixString("/advisories/new"), &controllers.AdvisoryNewController{})

	//beego.Router("/about", &controllers.AboutController{})

	//
	// --------------------------------------------------------------------
	// AUTHENTICATION [persona]
	// --------------------------------------------------------------------
	//
	beego.Router(util.GetPrefixString("/auth/check"), &models.PersonaCheckController{})
	beego.Router(util.GetPrefixString("/auth/login"), &models.PersonaLoginController{})
	beego.Router(util.GetPrefixString("/auth/logout"), &models.PersonaLogoutController{})

	//
	// --------------------------------------------------------------------
	// ADMINISTRATION [crap]
	// --------------------------------------------------------------------
	//
	beego.Router(util.GetPrefixString("/admin"), &controllers.AdminController{})

	//
	// --------------------------------------------------------------------
	// TEMPLATING FUNCTIONS
	// --------------------------------------------------------------------
	//
	beego.AddFuncMap("since", func(t time.Time) string {
		hrs := time.Since(t).Hours()
		return fmt.Sprintf("%dd %02dhrs", int(hrs)/24, int(hrs)%24)
	})

	beego.AddFuncMap("emailat", func(s string) string {
		return strings.Replace(s, "@", " [@T] ", -1)
	})

	beego.AddFuncMap("mapaccess", func(s interface{}, m map[string]string) string {
		return m[to.String(s)]
	})

	beego.AddFuncMap("url", util.GetPrefixString)
	beego.AddFuncMap("urldata", util.GetPrefixStringWithData)

	//
	// --------------------------------------------------------------------
	// ERROR HANDLERS
	// --------------------------------------------------------------------
	//
	beego.ErrorHandler("400", func(rw http.ResponseWriter, r *http.Request) {

		templateName := "errors/400.tpl"

		data := make(map[string]interface{})
		data["Title"] = "huh wut"
		data["Loc"] = -2
		data["Tab"] = -1
		data["copyright"] = time.Now().Year()

		if beego.BConfig.RunMode == "dev" {
			beego.BuildTemplate(beego.BConfig.WebConfig.ViewsPath)
		}

		newbytes := bytes.NewBufferString("")
		if _, ok := beego.BeeTemplates[templateName]; !ok {
			panic("can't find templatefile in the path:" + templateName)
		}
		err := beego.BeeTemplates[templateName].ExecuteTemplate(newbytes, templateName, data)
		if err != nil {
			panic("template Execute err: " + err.Error())
		}
		tplcontent, _ := ioutil.ReadAll(newbytes)
		fmt.Fprint(rw, template.HTML(string(tplcontent)))
	})

	beego.ErrorHandler("403", func(rw http.ResponseWriter, r *http.Request) {

		templateName := "errors/403.tpl"

		data := make(map[string]interface{})
		data["Title"] = "bzzzt..."
		data["Loc"] = -2
		data["Tab"] = -1
		data["copyright"] = time.Now().Year()

		if beego.BConfig.RunMode == "dev" {
			beego.BuildTemplate(beego.BConfig.WebConfig.ViewsPath)
		}

		newbytes := bytes.NewBufferString("")
		if _, ok := beego.BeeTemplates[templateName]; !ok {
			panic("can't find templatefile in the path:" + templateName)
		}
		err := beego.BeeTemplates[templateName].ExecuteTemplate(newbytes, templateName, data)
		if err != nil {
			panic("template Execute err: " + err.Error())
		}
		tplcontent, _ := ioutil.ReadAll(newbytes)
		fmt.Fprint(rw, template.HTML(string(tplcontent)))
	})

	beego.ErrorHandler("404", func(rw http.ResponseWriter, r *http.Request) {

		templateName := "errors/404.tpl"

		data := make(map[string]interface{})
		data["Title"] = "i have no idea what i'm doing"
		data["Loc"] = -2
		data["Tab"] = -1
		data["copyright"] = time.Now().Year()

		if j, found := util.Cache.Get("404_xkcd_json"); found {

			v := j.(map[string]interface{})
			data["xkcd_today"] = v["img"]
			data["xkcd_today_title"] = v["alt"]

		} else {
			resp, err := http.Get("http://xkcd.com/info.0.json")
			if err == nil {
				defer resp.Body.Close()
				bte, err := ioutil.ReadAll(resp.Body)

				if err == nil {
					var v map[string]interface{}

					if json.Unmarshal(bte, &v) == nil {
						util.Cache.Set("404_xkcd_json", v, 0)

						data["xkcd_today"] = v["img"]
						data["xkcd_today_title"] = v["alt"]
					}
				}
			}
		}

		if beego.BConfig.RunMode == "dev" {
			beego.BuildTemplate(beego.BConfig.WebConfig.ViewsPath)
		}

		newbytes := bytes.NewBufferString("")
		if _, ok := beego.BeeTemplates[templateName]; !ok {
			panic("can't find templatefile in the path:" + templateName)
		}
		err := beego.BeeTemplates[templateName].ExecuteTemplate(newbytes, templateName, data)
		if err != nil {
			panic("template Execute err: " + err.Error())
		}
		tplcontent, _ := ioutil.ReadAll(newbytes)
		fmt.Fprint(rw, template.HTML(string(tplcontent)))
	})

	beego.ErrorHandler("500", func(rw http.ResponseWriter, r *http.Request) {

		templateName := "errors/500.tpl"

		data := make(map[string]interface{})
		data["Title"] = "eek fire FIRE"
		data["Loc"] = -2
		data["Tab"] = -1
		data["copyright"] = time.Now().Year()

		if beego.BConfig.RunMode == "dev" {
			beego.BuildTemplate(beego.BConfig.WebConfig.ViewsPath)
		}

		newbytes := bytes.NewBufferString("")
		if _, ok := beego.BeeTemplates[templateName]; !ok {
			panic("can't find templatefile in the path:" + templateName)
		}
		err := beego.BeeTemplates[templateName].ExecuteTemplate(newbytes, templateName, data)
		if err != nil {
			panic("template Execute err: " + err.Error())
		}
		tplcontent, _ := ioutil.ReadAll(newbytes)
		fmt.Fprint(rw, template.HTML(string(tplcontent)))
	})

	beego.ErrorHandler("550", func(rw http.ResponseWriter, r *http.Request) {

		templateName := "errors/550.tpl"

		data := make(map[string]interface{})
		data["Title"] = "Oh No!"
		data["Permission"] = r.Form.Get("permission")
		data["Loc"] = -2
		data["Tab"] = -1
		data["copyright"] = time.Now().Year()

		data["xsrf_token"] = r.Form.Get("xsrf")

		if beego.BConfig.RunMode == "dev" {
			beego.BuildTemplate(beego.BConfig.WebConfig.ViewsPath)
		}

		newbytes := bytes.NewBufferString("")
		if _, ok := beego.BeeTemplates[templateName]; !ok {
			panic("can't find templatefile in the path:" + templateName)
		}
		err := beego.BeeTemplates[templateName].ExecuteTemplate(newbytes, templateName, data)
		if err != nil {
			panic("template Execute err: " + err.Error())
		}
		tplcontent, _ := ioutil.ReadAll(newbytes)
		fmt.Fprint(rw, template.HTML(string(tplcontent)))
	})

	//
	// --------------------------------------------------------------------
	// INTEGRATION
	// --------------------------------------------------------------------
	//
	stop := make(chan bool)

	integration.Integrate(integration.ABF(1))
	// ping target (for integration)
	//beego.Router("/ping", &controllers.PingController{})

	go func() {
		timeout := make(chan bool)
		go func() {
			for {
				timeout <- true
				time.Sleep(1 * time.Hour)
			}
		}()
		for {
			select {
			case <-stop:
				return
			case <-timeout:
				integration.Ping()
			}
		}
	}()

	beego.Run()
	<-stop
}

func getRandomString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
