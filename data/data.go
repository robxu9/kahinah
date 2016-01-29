package data

import (
	"net/http"
	"time"

	"github.com/robxu9/kahinah/render"
	"github.com/zenazn/goji/web/mutil"

	"golang.org/x/net/context"
	"gopkg.in/guregu/kami.v1"
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

func RenderMiddleware() kami.Middleware {
	return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) context.Context {
		return context.WithValue(ctx, datakey, &Render{
			Type:   DataHTML,
			Status: 200,
		})
	}
}

func RenderAfterware() kami.Afterware {
	return func(ctx context.Context, wp mutil.WriterProxy, r *http.Request) context.Context {
		ret := FromContext(ctx)
		renderer := render.FromContext(ctx)

		switch ret.Type {
		case DataNoRender:
			break
		case DataHTML:
			if m, ok := ret.Data.(map[string]interface{}); ok {
				m["copyright"] = time.Now().Year()
			}
			renderer.HTML(wp, ret.Status, ret.Template, ret.Data)
		case DataJSON:
			renderer.JSON(wp, ret.Status, ret.Data)
		case DataBinary:
			renderer.Data(wp, ret.Status, ret.Data.([]byte))
		case DataText:
			renderer.Text(wp, ret.Status, ret.Data.(string))
		case DataJSONP:
			renderer.JSONP(wp, ret.Status, ret.Callback, ret.Data)
		case DataXML:
			renderer.XML(wp, ret.Status, ret.Data)
		default:
			panic("no such data type")
		}

		return ctx
	}
}
