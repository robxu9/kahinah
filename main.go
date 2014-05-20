package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"github.com/robxu9/kahinah/util"

	"github.com/astaxie/beego"
	"github.com/robxu9/kahinah/controllers"
	"github.com/robxu9/kahinah/integration"
	"github.com/robxu9/kahinah/models"
	"menteslibres.net/gosexy/to"
)

func main() {
	beego.SessionOn = true

	beego.EnableXSRF = true
	beego.XSRFKEY = getRandomString(50)
	beego.XSRFExpire = 3600

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
	// ERROR HANDLERS [beego doesn't want to work though e-e]
	// --------------------------------------------------------------------
	//
	beego.Errorhandler("550", func(rw http.ResponseWriter, r *http.Request) {

		templateName := "errors/550.tpl"

		data := make(map[string]interface{})
		data["Title"] = "Oh No!"
		data["Permission"] = r.Form.Get("permission")
		data["Loc"] = -2
		data["Tab"] = -1
		data["copyright"] = time.Now().Year()

		data["xsrf_token"] = r.Form.Get("xsrf")

		if beego.RunMode == "dev" {
			beego.BuildTemplate(beego.ViewsPath)
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
