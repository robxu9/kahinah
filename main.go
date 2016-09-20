package main

import (
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
	"github.com/gorilla/securecookie"
	"github.com/knq/kv"
	"github.com/knq/sessionmw"
	"github.com/robfig/cron"
	"github.com/robxu9/kahinah/common/conf"
	"github.com/robxu9/kahinah/common/klog"
	"github.com/robxu9/kahinah/controllers"
	"github.com/robxu9/kahinah/data"
	"github.com/robxu9/kahinah/integration"
	"github.com/robxu9/kahinah/job"
	"github.com/robxu9/kahinah/models"
	"github.com/robxu9/kahinah/render"
	urender "github.com/unrolled/render"
	"github.com/zenazn/goji/web/mutil"
	"menteslibres.net/gosexy/to"
)

func init() {
	http.DefaultClient = &http.Client{
		Timeout: time.Second * 30, // we need a sane default timeout
	}
}

func main() {
	klog.Info("starting kahinah v4")

	// -- mux -----------------------------------------------------------------
	mux := goji.NewMux()

	// -- middleware ----------------------------------------------------------

	// logging middleware (base middleware)
	mux.UseC(func(inner goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
			klog.Debugf("req  (%v): path (%v), user-agent (%v), referrer (%v)", r.RemoteAddr, r.RequestURI, r.UserAgent(), r.Referer())
			wp := mutil.WrapWriter(rw) // proxy the rw for info later
			inner.ServeHTTPC(ctx, wp, r)
			klog.Debugf("resp (%v): status (%v), bytes written (%v)", r.RemoteAddr, wp.Status(), wp.BytesWritten())
		})
	})

	// rendering middleware (required by panic)
	renderer := urender.New(urender.Options{
		Directory:  "views",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".tpl"},
		Funcs: []template.FuncMap{
			template.FuncMap{
				"rfc3339": func(t time.Time) string {
					return t.Format(time.RFC3339)
				},
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
				"url":     render.ConvertURL,
				"urldata": render.ConvertURLWithData,
			},
		},
		IndentJSON:    true,
		IndentXML:     true,
		IsDevelopment: conf.GetDefault("runMode", "dev").(string) == "dev",
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
	if enable, ok := conf.GetDefault("authentication.cas.enable", false).(bool); ok && enable {
		url, _ := url.Parse(conf.Get("authentication.cas.url").(string))

		casClient := cas.NewClient(&cas.Options{
			URL: url,
		})

		mux.Use(casClient.Handle)
	}

	// sessions middleware
	sessionConfig := &sessionmw.Config{
		Secret:      []byte(securecookie.GenerateRandomKey(64)),
		BlockSecret: []byte(securecookie.GenerateRandomKey(32)),
		Store:       kv.NewMemStore(),
		Name:        "kahinah",
	}
	mux.UseC(sessionConfig.Handler)

	// csrf middleware
	mux.UseC(csrf.Protect(securecookie.GenerateRandomKey(64), csrf.Secure(false)))

	// data rendering middleware
	mux.UseC(func(inner goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
			newCtx := data.RenderMiddleware(ctx, rw, r)
			inner.ServeHTTPC(newCtx, rw, r)
			data.RenderAfterware(newCtx, rw, r)
		})
	})

	// --------------------------------------------------------------------
	// HANDLERS
	// --------------------------------------------------------------------

	getHandlers := map[string]goji.HandlerFunc{

		// static paths
		"/static/*": controllers.StaticHandler,

		// main page
		"/": controllers.MainHandler,

		// builds - json
		"/i/list/json/": controllers.ListsAPIHandler,

		// builds
		"/i/list/": controllers.ListsHandler,

		// build - specific
		"/b/:id/": controllers.BuildGetHandler,

		// build - specific - json
		"/b/:id/json/": controllers.BuildGetJSONHandler,

		// activity - json
		"/i/activity/json/": controllers.ActivityJSONHandler,

		// activity - html
		"/i/activity/": controllers.ActivityHandler,

		// admin
		"/admin/": controllers.AdminGetHandler,

		// authentication - login
		"/u/login/": controllers.UserLoginHandler,

		// authentication - logout
		"/u/logout/": controllers.UserLogoutHandler,
	}

	postHandlers := map[string]goji.HandlerFunc{

		// webhooks
		"/hook/*": controllers.IntegrationHandler,

		// build - specific
		"/b/:id/": controllers.BuildPostHandler,

		// admin
		"/admin/": controllers.AdminPostHandler,
	}

	for k, v := range getHandlers {
		if len(k) > 1 && strings.HasSuffix(k, "/") {
			getHandlerRedirectName := render.ConvertURLRelative(k[:len(k)-1])
			klog.Debugf("get handler setup: redirecting %v", getHandlerRedirectName)
			mux.HandleFunc(pat.Get(getHandlerRedirectName), controllers.RedirectHandler)
		}
		getHandlerUseName := render.ConvertURLRelative(k)
		klog.Debugf("get handler setup: using %v", getHandlerUseName)
		mux.HandleC(pat.Get(getHandlerUseName), v)
	}

	for k, v := range postHandlers {
		if len(k) > 1 && strings.HasSuffix(k, "/") {
			postHandlerRedirectName := render.ConvertURLRelative(k[:len(k)-1])
			klog.Debugf("post handler setup: redirecting %v", postHandlerRedirectName)
			mux.HandleFunc(pat.Post(postHandlerRedirectName), controllers.RedirectHandler)
		}
		postHandlerUseName := render.ConvertURLRelative(k)
		klog.Debugf("post handler setup: using %v", postHandlerUseName)
		mux.HandleC(pat.Post(postHandlerUseName), v)
	}

	// -- cronjobs ----------------------------------------------------------

	cronRunner := cron.New()

	// integration polling
	if pollRate, ok := conf.Get("integration.poll").(string); ok && pollRate != "" {
		pollFunc := func() {
			pollAllErr := integration.PollAll()
			for name, err := range pollAllErr {
				klog.Warningf("integration polling failed for %v: %v", name, err)
			}
		}
		cronRunner.AddFunc(pollRate, pollFunc)

		// and do an initial poll
		pollFunc()
	}

	// process new stages/check processes every 10 seconds
	cronRunner.AddFunc("@every 10s", func() {
		models.CheckAllListStages()
	})

	// run the job scheduler every minute
	cronRunner.AddFunc("@every 1m", func() {
		job.ProcessQueue()
	})

	// start cron
	cronRunner.Start()

	// -- server setup --------------------------------------------------------

	// bind and listen
	listenAddr := conf.GetDefault("listenAddr", "0.0.0.0").(string)
	listenPort := conf.GetDefault("listenPort", 3000).(int64)

	klog.Infof("listening to %v:%v", listenAddr, listenPort)
	if err := http.ListenAndServe(fmt.Sprintf("%v:%v", listenAddr, listenPort), mux); err != nil {
		klog.Fatalf("unable to serve: %v", err)
	}

	cronRunner.Stop()

	klog.Info("processing leftover jobs...")
	close(job.Queue)
	for len(job.Queue) > 0 {
		job.ProcessQueue()
	}
}
