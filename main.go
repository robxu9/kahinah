package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/robxu9/kahinah/controllers"
	"github.com/robxu9/kahinah/integration"
	"github.com/robxu9/kahinah/models"
	"html/template"
	"menteslibres.net/gosexy/to"
	"net/http"
	"strings"
	"time"
)

var (
	PREFIX = beego.AppConfig.String("urlprefix")
)

func main() {
	beego.SessionOn = true

	beego.EnableXSRF = true
	beego.XSRFKEY = getRandomString(50)
	beego.XSRFExpire = 3600

	beego.SetStaticPath(getPrefixString("/static"), "static")

	beego.Router(getPrefixString("/"), &controllers.MainController{})

	// testing
	beego.Router(getPrefixString("/testing"), &controllers.TestingController{}) // lists testing updates
	//beego.Router("/testing/:buildid:int", &controllers.TestingPkgController{}) // shows specific testing update
	// ^ now use below BuildSpecificController

	// published
	beego.Router(getPrefixString("/published"), &controllers.PublishedController{}) // lists years

	// below: TODO in future update
	//beego.Router("/published/:year:int", &controllers.PublishedListController{})        // lists updates in said years
	//beego.Router("/published/:year:int/:id:int", &controllers.PublishedPkgController{}) // shows specific update
	// ^ redirect to below build specific controller

	// rejected
	beego.Router(getPrefixString("/rejected"), &controllers.RejectedController{}) // lists all rejected updates
	//beego.Router("/rejected/:before:int", &controllers.RejectedController{}) // list all rejected updates before this date

	// builds
	beego.Router(getPrefixString("/builds"), &controllers.BuildsController{}) // show all testing, published, rejected (all sorted by date, linking respectively to above)
	beego.Router(getPrefixString("/builds/:id:int"), &controllers.BuildController{})

	// platform
	//beego.Router("/platforms", &controllers.PlatformsController{})
	//beego.Router("/platforms/:platform:string/", &controllers.PlatformSpecificController{}) // show by platform

	// about
	//beego.Router("/about", &controllers.AboutController{})

	// ping target (for integration)
	//beego.Router("/ping", &controllers.PingController{})

	// persona
	beego.Router(getPrefixString("/auth/check"), &models.PersonaCheckController{})
	beego.Router(getPrefixString("/auth/login"), &models.PersonaLoginController{})
	beego.Router(getPrefixString("/auth/logout"), &models.PersonaLogoutController{})

	// admin
	beego.Router(getPrefixString("/admin"), &controllers.AdminController{})

	// templating
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

	beego.AddFuncMap("url", getPrefixString)
	beego.AddFuncMap("urldata", getPrefixStringWithData)

	// error handling
	beego.Errorhandler("550", func(rw http.ResponseWriter, r *http.Request) {
		t := template.Must(template.New("permerror").ParseFiles(beego.ViewsPath + "/permissions_error.tpl"))
		data := make(map[string]interface{})
		data["Permission"] = r.Form.Get("permission")
		t.Execute(rw, data)
	})

	// integration
	stop := make(chan bool)

	integration.Integrate(integration.ABF(1))

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

func getPrefixStringWithData(dest string, data interface{}) string {
	// no need to prefix if the dest has no / before it
	temp := template.Must(template.New("prefixTemplate").Parse(dest))
	var b bytes.Buffer

	err := temp.Execute(&b, data)
	if err != nil {
		panic(err)
	}

	result := b.String()
	return getPrefixString(result)
}

func getPrefixString(dest string) string {
	if PREFIX == "" {
		return dest
	}

	return "/" + PREFIX + dest
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
