package data

import (
	"net/http"
	"time"

	"gopkg.in/cas.v1"

	"github.com/goji/ctx-csrf"
	"github.com/knq/sessionmw"
	"github.com/robxu9/kahinah/common/conf"
	"github.com/robxu9/kahinah/render"

	"golang.org/x/net/context"
)

type key int

var datakey key

const (
	DataNoRender int = iota
	DataHTML
	DataBinary
	DataText
	DataJSON
	DataJSONP
	DataXML
)

const (
	FlashErr  = "_flash_err"
	FlashWarn = "_flash_warn"
	FlashInfo = "_flash_info"
)

var (
	runMode = conf.GetDefault("runMode", "dev").(string)
)

type Render struct {
	Type     int
	Status   int
	Template string // for DataHTML
	Callback string // for DataJSONP
	Data     interface{}
}

func FromContext(ctx context.Context) *Render {
	r, ok := ctx.Value(datakey).(*Render)
	if !ok {
		panic("unable to retrieve render struct from context")
	}
	return r
}

func RenderMiddleware(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
	return context.WithValue(ctx, datakey, &Render{
		Type:   DataHTML,
		Status: 200,
	})
}

func RenderAfterware(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	ret := FromContext(ctx)
	renderer := render.FromContext(ctx)

	// Set the CSRF token
	rw.Header().Set("X-CSRF-Token", csrf.Token(ctx, r))

	switch ret.Type {
	case DataNoRender:
		break
	case DataHTML:
		if ret.Template == "" {
			// guess we're not rendering anything
			break
		}
		if m, ok := ret.Data.(map[string]interface{}); ok {
			// Set the copyright on all pages
			m["copyright"] = time.Now().Year()

			// Add xsrf tokens
			m["xsrf_token"] = csrf.Token(ctx, r)
			m["xsrf_data"] = csrf.TemplateField(ctx, r)

			// Add environment declaration
			m["environment"] = runMode

			// Add Nav info if it doesn't already exist
			if _, ok := m["Nav"]; !ok {
				m["Nav"] = -1
			}

			// Add authentication information
			m["authenticated"] = cas.Username(r)

			// Add session flash stuff
			if f, has := sessionmw.Get(ctx, FlashErr); has {
				m["flash_err"] = f
				sessionmw.Delete(ctx, FlashErr)
			}
			if f, has := sessionmw.Get(ctx, FlashWarn); has {
				m["flash_warn"] = f
				sessionmw.Delete(ctx, FlashWarn)
			}
			if f, has := sessionmw.Get(ctx, FlashInfo); has {
				m["flash_info"] = f
				sessionmw.Delete(ctx, FlashInfo)
			}
		}
		renderer.HTML(rw, ret.Status, ret.Template, ret.Data)
	case DataJSON:
		renderer.JSON(rw, ret.Status, ret.Data)
	case DataBinary:
		renderer.Data(rw, ret.Status, ret.Data.([]byte))
	case DataText:
		renderer.Text(rw, ret.Status, ret.Data.(string))
	case DataJSONP:
		renderer.JSONP(rw, ret.Status, ret.Callback, ret.Data)
	case DataXML:
		renderer.XML(rw, ret.Status, ret.Data)
	default:
		panic("no such data type")
	}
}
