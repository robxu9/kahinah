package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/cas.v1"
	"gopkg.in/errors.v0"

	"goji.io"

	"github.com/robxu9/kahinah/common/conf"
	"github.com/robxu9/kahinah/common/klog"
	"github.com/robxu9/kahinah/render"
	"github.com/wunderlist/ttlcache"

	"golang.org/x/net/context"
)

var (
	ErrNotFound         = errors.New("kahinah: path not found")
	ErrMethodNotAllowed = errors.New("kahinah: method not allowed")
	ErrBadRequest       = errors.New("kahinah: bad request")
	ErrForbidden        = errors.New("kahinah: forbidden")
	ErrPermission       = errors.New("kahinah: permission error")

	runMode = conf.GetDefault("runMode", "dev").(string)
)

func PanicMiddleware(inner goji.Handler) goji.Handler {
	panicHandler := &PanicHandler{
		Cache: ttlcache.NewCache(24 * time.Hour),
	}

	return goji.HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				panicHandler.ServeHTTPC(err, ctx, rw, r)
			}
		}()

		inner.ServeHTTPC(ctx, rw, r)
	})
}

type PanicHandler struct {
	Cache *ttlcache.Cache
}

func (p *PanicHandler) ServeHTTPC(ex interface{}, ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	switch ex {
	case ErrNotFound:
		p.Err404(ctx, rw, r)
	case ErrMethodNotAllowed:
		p.Err405(ctx, rw, r)
	case ErrBadRequest:
		p.Err400(ctx, rw, r)
	case ErrForbidden:
		p.Err403(ctx, rw, r)
	case ErrPermission:
		p.Err550(ctx, rw, r)
	default:
		p.Err500(ex, ctx, rw, r)
	}
}

func (p *PanicHandler) Err404(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	renderer := render.FromContext(ctx)

	result := make(map[string]interface{})
	result["Title"] = "i have no idea what i'm doing"
	result["Nav"] = -1
	result["environment"] = runMode
	result["authenticated"] = cas.Username(r)
	result["copyright"] = time.Now().Year()

	if xImg, found := p.Cache.Get("404-xkcd-img"); found {
		result["xkcd_today"] = xImg
		xAlt, _ := p.Cache.Get("404-xkcd-alt")
		result["xkcd_today_alt"] = xAlt
	} else {
		if resp, err := http.Get("https://www.xkcd.com/info.0.json"); err == nil {
			defer resp.Body.Close()

			bte, err := ioutil.ReadAll(resp.Body)

			if err == nil {
				var v map[string]interface{}

				if json.Unmarshal(bte, &v) == nil {
					p.Cache.Set("404-xkcd-img", v["img"].(string))
					p.Cache.Set("404-xkcd-alt", v["alt"].(string))

					result["xkcd_today"] = v["img"]
					result["xkcd_today_title"] = v["alt"]
				}
			}
		}
	}

	renderer.HTML(rw, 404, "errors/404", result)
}

func (p *PanicHandler) Err400(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	renderer := render.FromContext(ctx)

	data := make(map[string]interface{})
	data["Title"] = "huh wut"
	data["Nav"] = -1
	data["environment"] = runMode
	data["authenticated"] = cas.Username(r)
	data["copyright"] = time.Now().Year()

	renderer.HTML(rw, 400, "errors/400", data)
}

func (p *PanicHandler) Err403(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	renderer := render.FromContext(ctx)

	data := make(map[string]interface{})
	data["Title"] = "bzzzt..."
	data["Nav"] = -1
	data["environment"] = runMode
	data["authenticated"] = cas.Username(r)
	data["copyright"] = time.Now().Year()

	renderer.HTML(rw, 403, "errors/403", data)
}

func (p *PanicHandler) Err405(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	renderer := render.FromContext(ctx)

	renderer.HTML(rw, 405, "errors/405", map[string]interface{}{
		"Title":         "Method Not Allowed",
		"Nav":           -1,
		"environment":   runMode,
		"authenticated": cas.Username(r),
		"copyright":     time.Now().Year(),
	})
}

func (p *PanicHandler) Err500(ex interface{}, ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	renderer := render.FromContext(ctx)

	data := make(map[string]interface{})
	data["Title"] = "eek fire FIRE"
	data["Nav"] = -1
	data["environment"] = runMode
	data["authenticated"] = cas.Username(r)
	data["copyright"] = time.Now().Year()
	data["error"] = fmt.Sprintf("%v", ex)

	stackTrace := errors.Wrap(ex, 4).ErrorStack()

	klog.Critical("err  (%v): Internal Server Error: %v", r.RemoteAddr, ex)
	klog.Critical("err  (%v): Stacktrace: %v", r.RemoteAddr, stackTrace)

	if mode := conf.GetDefault("runMode", "dev").(string); mode == "dev" {
		// dump the stacktrace out on the page too
		data["stacktrace"] = stackTrace
	}

	renderer.HTML(rw, 500, "errors/500", data)
}

func (p *PanicHandler) Err550(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	renderer := render.FromContext(ctx)

	data := make(map[string]interface{})
	data["Title"] = "Oh No!"
	//data["Permission"] = r.Form.Get("permission") // FIXME
	data["Nav"] = -1
	data["environment"] = runMode
	data["authenticated"] = cas.Username(r)
	data["copyright"] = time.Now().Year()

	renderer.HTML(rw, 550, "errors/550", data)
}
