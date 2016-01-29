package main

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"

	"gopkg.in/guregu/kami.v1"

	"github.com/gorilla/securecookie"
	"github.com/robxu9/kahinah/conf"
	"github.com/robxu9/kahinah/controllers"
	"github.com/robxu9/kahinah/data"
	"github.com/robxu9/kahinah/log"
	"github.com/robxu9/kahinah/render"
	"github.com/robxu9/kahinah/sessions"
	"github.com/robxu9/kahinah/util"
	urender "github.com/unrolled/render"
	"github.com/zenazn/goji/web/mutil"
	"menteslibres.net/gosexy/to"
)

func main() {
	// Initialise an empty context
	ctx := context.Background()

	// Set up logging
	kami.LogHandler = func(ctx context.Context, wp mutil.WriterProxy, r *http.Request) {
		log.Logger.Debugf("request: addr (%v), path (%v), user-agent (%v), referrer (%v)", r.RemoteAddr, r.RequestURI, r.UserAgent(), r.Referer())
		log.Logger.Debugf("response: status (%v), bytes written (%v)", wp.Status(), wp.BytesWritten())
	}

	// Set up rendering
	r := urender.New(urender.Options{
		Directory:  "views",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".tpl"},
		Funcs: []template.FuncMap{
			template.FuncMap{
				"since": func(t time.Time) string {
					hrs := time.Since(t).Hours()
					return fmt.Sprintf("%dd %02dhrs", int(hrs)/24, int(hrs)%24)
				},
				"emailat": func(s string) string {
					return strings.Replace(s, "@", " [@T] ", -1)
				},
				"mapaccess": func(s interface{}, m map[string]string) string {
					return m[to.String(s)]
				},
				"url":     util.GetPrefixString,
				"urldata": util.GetPrefixStringWithData,
			},
		},
		IndentJSON:    true,
		IndentXML:     true,
		IsDevelopment: conf.Config.GetDefault("runMode", "dev").(string) == "dev",
	})
	ctx = render.NewContext(ctx, r)

	kami.Use("/", data.RenderMiddleware())
	kami.After("/", data.RenderAfterware()) // merge data into render pkg?

	// Set up sessions
	ctx = sessions.NewStore(ctx, securecookie.GenerateRandomKey(64))
	kami.Use("/", sessions.SessionMiddleware("kahinah"))
	kami.After("/", sessions.SessionAfterware("kahinah"))

	// FIXME: set up xsrf (getRandomString(50), expire in 3600 minutes)

	// Set up the error handlers
	panicHandler := &controllers.PanicHandler{}
	kami.PanicHandler = panicHandler
	kami.NotFound(panicHandler.Err404)
	kami.MethodNotAllowed(panicHandler.Err405)

	// Set it up as the god context
	kami.Context = ctx

	// --------------------------------------------------------------------
	// HANDLERS
	// --------------------------------------------------------------------

	// Handle static paths
	kami.Get(util.GetPrefixString("/static/*path"), controllers.StaticHandler)

	// Show the main page
	kami.Get(util.GetPrefixString("/"), controllers.MainHandler)

	//
	// --------------------------------------------------------------------
	// BUILDS
	// --------------------------------------------------------------------
	//

	// testing
	kami.Get(util.GetPrefixString("/builds/testing"), controllers.TestingHandler)
	// published
	kami.Get(util.GetPrefixString("/builds/published"), controllers.PublishedHandler)
	// rejected
	kami.Get(util.GetPrefixString("/builds/rejected"), controllers.RejectedHandler) // list all rejected updates
	// all builds
	kami.Get(util.GetPrefixString("/builds"), controllers.BuildsHandler) // show all updates sorted by date

	// // specific
	// beego.Router(util.GetPrefixString("/builds/:id:int"), &controllers.BuildController{})
	//
	// //
	// // --------------------------------------------------------------------
	// // ADVISORIES
	// // --------------------------------------------------------------------
	// //
	//
	// // advisories
	// beego.Router(util.GetPrefixString("/advisories"), &controllers.AdvisoryMainController{})
	// //beego.Router(util.GetPrefixString("/advisories/:platform:string"), &controllers.AdvisoryPlatformController{})
	// //beego.Router(util.GetPrefixString("/advisories/:id:int"), &controllers.AdvisoryController{})
	// beego.Router(util.GetPrefixString("/advisories/new"), &controllers.AdvisoryNewController{})
	//
	// //beego.Router("/about", &controllers.AboutController{})
	//
	// //
	// // --------------------------------------------------------------------
	// // AUTHENTICATION [persona]
	// // --------------------------------------------------------------------
	// //
	// beego.Router(util.GetPrefixString("/auth/check"), &models.PersonaCheckController{})
	// beego.Router(util.GetPrefixString("/auth/login"), &models.PersonaLoginController{})
	// beego.Router(util.GetPrefixString("/auth/logout"), &models.PersonaLogoutController{})
	//
	// //
	// // --------------------------------------------------------------------
	// // ADMINISTRATION [crap]
	// // --------------------------------------------------------------------
	// //
	// beego.Router(util.GetPrefixString("/admin"), &controllers.AdminController{})
	//
	// //
	// // --------------------------------------------------------------------
	// // INTEGRATION
	// // --------------------------------------------------------------------
	// //
	// stop := make(chan bool)
	//
	// integration.Integrate(integration.ABF(1))
	// // ping target (for integration)
	// //beego.Router("/ping", &controllers.PingController{})
	//
	// go func() {
	// 	timeout := make(chan bool)
	// 	go func() {
	// 		for {
	// 			timeout <- true
	// 			time.Sleep(1 * time.Hour)
	// 		}
	// 	}()
	// 	for {
	// 		select {
	// 		case <-stop:
	// 			return
	// 		case <-timeout:
	// 			integration.Ping()
	// 		}
	// 	}
	// }()
	//
	// beego.Run()
	// <-stop

	kami.Serve()
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
