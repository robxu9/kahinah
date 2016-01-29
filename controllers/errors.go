package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"

	"github.com/robxu9/kahinah/conf"
	"github.com/robxu9/kahinah/data"
	"github.com/robxu9/kahinah/render"
	"github.com/robxu9/kahinah/util"
	"gopkg.in/guregu/kami.v1"

	"golang.org/x/net/context"
)

var (
	ErrNotFound         = errors.New("kami: path not found")
	ErrMethodNotAllowed = errors.New("kami: method not allowed")
	ErrBadRequest       = errors.New("kami: bad request")
	ErrForbidden        = errors.New("kami: forbidden")
	ErrPermission       = errors.New("kami: permission error")
)

type PanicHandler struct{}

func (p *PanicHandler) ServeHTTPContext(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	ex := kami.Exception(ctx) // retrieve the panic

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
		p.Err500(ctx, rw, r)
	}
}

func (p *PanicHandler) Err404(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	renderer := render.FromContext(ctx)

	result := make(map[string]interface{})
	result["Title"] = "i have no idea what i'm doing"
	result["Loc"] = -2
	result["Tab"] = -1

	if j, found := util.Cache.Get("404_xkcd_json"); found {

		v := j.(map[string]interface{})
		result["xkcd_today"] = v["img"]
		result["xkcd_today_title"] = v["alt"]

	} else {
		resp, err := http.Get("http://xkcd.com/info.0.json")
		if err == nil {
			defer resp.Body.Close()
			bte, err := ioutil.ReadAll(resp.Body)

			if err == nil {
				var v map[string]interface{}

				if json.Unmarshal(bte, &v) == nil {
					util.Cache.Set("404_xkcd_json", v, 0)

					result["xkcd_today"] = v["img"]
					result["xkcd_today_title"] = v["alt"]
				}
			}
		}
	}

	dataRenderer := data.FromContext(ctx)
	dataRenderer.Type = data.DataNoRender

	renderer.HTML(rw, 404, "errors/404", result)
}

func (p *PanicHandler) Err400(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	renderer := render.FromContext(ctx)

	data := make(map[string]interface{})
	data["Title"] = "huh wut"
	data["Loc"] = -2
	data["Tab"] = -1

	renderer.HTML(rw, 400, "errors/400", data)
}

func (p *PanicHandler) Err403(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	renderer := render.FromContext(ctx)

	data := make(map[string]interface{})
	data["Title"] = "bzzzt..."
	data["Loc"] = -2
	data["Tab"] = -1

	renderer.HTML(rw, 403, "errors/403", data)
}

func (p *PanicHandler) Err405(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	renderer := render.FromContext(ctx)

	dataRenderer := data.FromContext(ctx)
	dataRenderer.Type = data.DataNoRender

	renderer.HTML(rw, 405, "errors/405", map[string]interface{}{
		"Title": "Method Not Allowed",
		"Loc":   -2,
		"Tab":   -1,
	})
}

func (p *PanicHandler) Err500(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	renderer := render.FromContext(ctx)

	data := make(map[string]interface{})
	data["Title"] = "eek fire FIRE"
	data["Loc"] = -2
	data["Tab"] = -1
	data["error"] = fmt.Sprintf("%v", kami.Exception(ctx))

	if mode := conf.Config.GetDefault("runMode", "dev").(string); mode == "dev" {
		// dump the stacktrace out on the page too
		var trace []byte
		runtime.Stack(trace, false)
		data["stacktrace"] = string(trace)
	}

	renderer.HTML(rw, 500, "errors/500", data)
}

func (p *PanicHandler) Err550(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	renderer := render.FromContext(ctx)

	data := make(map[string]interface{})
	data["Title"] = "Oh No!"
	//data["Permission"] = r.Form.Get("permission") // FIXME
	data["Loc"] = -2
	data["Tab"] = -1

	renderer.HTML(rw, 550, "errors/550", data)
}
