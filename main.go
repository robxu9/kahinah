package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"goji.io/middleware"
	"goji.io/pat"

	"goji.io"

	"golang.org/x/net/context"

	"gopkg.in/cas.v1"

	"github.com/goji/ctx-csrf"
	"github.com/robxu9/kahinah/conf"
	"github.com/robxu9/kahinah/controllers"
	"github.com/robxu9/kahinah/data"
	"github.com/robxu9/kahinah/log"
	"github.com/robxu9/kahinah/render"
	"github.com/robxu9/kahinah/util"
	urender "github.com/unrolled/render"
	"github.com/zenazn/goji/bind"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web/mutil"
	"menteslibres.net/gosexy/to"
)

func init() {
	bind.WithFlag()
	graceful.DoubleKickWindow(2 * time.Second)
}

func main() {
	log.Logger.Info("starting kahinah v4")

	// -- mux -----------------------------------------------------------------
	mux := goji.NewMux()

	// -- middleware ----------------------------------------------------------

	// logging middleware (base middleware)
	mux.UseC(func(inner goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
			log.Logger.Debugf("req  (%v): path (%v), user-agent (%v), referrer (%v)", r.RemoteAddr, r.RequestURI, r.UserAgent(), r.Referer())
			wp := mutil.WrapWriter(rw) // proxy the rw for info later
			inner.ServeHTTPC(ctx, wp, r)
			log.Logger.Debugf("resp (%v): status (%v), bytes written (%v)", r.RemoteAddr, wp.Status(), wp.BytesWritten())
		})
	})

	// rendering middleware (required by panic)
	renderer := urender.New(urender.Options{
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

	mux.UseC(func(inner goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
			newCtx := render.NewContext(ctx, renderer)
			inner.ServeHTTPC(newCtx, rw, r)
		})
	})

	// panic middleware
	mux.UseC(controllers.PanicMiddleware)

	// not found middleware
	mux.UseC(func(inner goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
			routeFound := middleware.Pattern(ctx)

			if routeFound != nil {
				inner.ServeHTTPC(ctx, rw, r)
				return
			}

			panic(controllers.ErrNotFound)
		})
	})

	// authentication (cas) middleware
	if enable, ok := conf.Config.GetDefault("authentication.cas.enable", false).(bool); ok && enable {
		url, _ := url.Parse(conf.Config.Get("authentication.cas.url").(string))

		casClient := cas.NewClient(&cas.Options{
			URL: url,
		})

		mux.Use(casClient.Handle)
	}

	// sessions middleware
	// sessionConfig := &sessionmw.Config{
	// 	Secret:      []byte(getRandomString(30)),
	// 	BlockSecret: []byte(getRandomString(30)),
	// 	Store:       kv.NewMemStore(),
	// 	Name:        "kahinah",
	// }
	// mux.UseC(sessionConfig.Handler)

	// csrf middleware
	mux.UseC(csrf.Protect([]byte(getRandomString(30))))

	// data rendering middleware
	mux.UseC(func(inner goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
			newCtx := data.RenderMiddleware(ctx, rw, r)
			inner.ServeHTTPC(newCtx, rw, r)
			data.RenderAfterware(newCtx, rw, r)
		})
	})

	// session middleware - yikes? sessionmw?
	// ctx = sessions.NewStore(ctx, securecookie.GenerateRandomKey(64))
	// kami.Use("/", sessions.SessionMiddleware("kahinah"))
	// kami.After("/", sessions.SessionAfterware("kahinah"))

	// FIXME: set up xsrf (getRandomString(50), expire in 3600 minutes)

	// --------------------------------------------------------------------
	// HANDLERS
	// --------------------------------------------------------------------

	getHandlers := map[string]goji.HandlerFunc{

		// static paths
		"/static/*": controllers.StaticHandler,

		// main page
		"/": controllers.MainHandler,

		// build - specific
		"/builds/:id/": controllers.BuildGetHandler,

		// build - testing
		"/builds/testing": controllers.TestingHandler,

		// build - published
		"/builds/published": controllers.PublishedHandler,

		// build - rejected
		"/builds/rejected": controllers.RejectedHandler,

		// build - all builds
		"/builds": controllers.BuildsHandler,

		// admin
		"/admin": controllers.AdminGetHandler,

		// authentication - login
		"/user/login": controllers.UserLoginHandler,

		// authentication - logout
		"/user/logout": controllers.UserLogoutHandler,
	}

	postHandlers := map[string]goji.HandlerFunc{

		// build - specific
		"/builds/:id/": controllers.BuildPostHandler,

		// admin
		"/admin": controllers.AdminPostHandler,
	}

	for k, v := range getHandlers {
		mux.HandleC(pat.Get(util.GetPrefixString(k)), v)
	}

	for k, v := range postHandlers {
		mux.HandleC(pat.Post(util.GetPrefixString(k)), v)
	}

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

	// bind and listen
	if !flag.Parsed() {
		flag.Parse()
	}

	listener := bind.Default()
	log.Logger.Infof("binding to %v", bind.Sniff())
	graceful.HandleSignals()
	bind.Ready()

	err := graceful.Serve(listener, mux)

	if err != nil {
		log.Logger.Fatalf("unable to serve: %v", err)
	}

	graceful.Wait()
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
