package main

import (
	"github.com/astaxie/beego"
	"github.com/robxu9/kahinah/controllers"
	"github.com/robxu9/kahinah/integration"
	"github.com/robxu9/kahinah/util"
	"strings"
	"time"
)

func main() {

	beego.Router("/", &controllers.MainController{})

	// testing
	beego.Router("/testing", &controllers.TestingController{}) // lists testing updates
	//beego.Router("/testing/:buildid:int", &controllers.TestingPkgController{}) // shows specific testing update
	// ^ now use below BuildSpecificController

	// published
	beego.Router("/published", &controllers.PublishedController{}) // lists years

	// below: TODO in future update
	//beego.Router("/published/:year:int", &controllers.PublishedListController{})        // lists updates in said years
	//beego.Router("/published/:year:int/:id:int", &controllers.PublishedPkgController{}) // shows specific update
	// ^ redirect to below build specific controller

	// rejected
	beego.Router("/rejected", &controllers.RejectedController{}) // lists all rejected updates
	//beego.Router("/rejected/:before:int", &controllers.RejectedController{}) // list all rejected updates before this date

	// builds
	beego.Router("/builds", &controllers.BuildsController{}) // show all testing, published, rejected (all sorted by date, linking respectively to above)
	beego.Router("/builds/:buildid:int", &controllers.BuildController{})

	// platform
	//beego.Router("/platforms", &controllers.PlatformsController{})
	//beego.Router("/platforms/:platform:string/", &controllers.PlatformSpecificController{}) // show by platform

	// about
	//beego.Router("/about", &controllers.AboutController{})

	// ping target (for integration)
	//beego.Router("/ping", &controllers.PingController{})

	// persona
	beego.Router("/auth/check", &util.PersonaCheckController{})
	beego.Router("/auth/login", &util.PersonaLoginController{})
	beego.Router("/auth/logout", &util.PersonaLogoutController{})

	// templating
	beego.AddFuncMap("since", func(t time.Time) time.Duration {
		return time.Since(t)
	})

	beego.AddFuncMap("emailat", func(s string) string {
		return strings.Replace(s, "@", " [aT] ", -1)
	})

	// integration
	stop := make(chan bool)

	integration.Integrate("abf", integration.ABF(1))

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
