package data

import (
	"net/http"
	"time"

	"gopkg.in/cas.v1"

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

	switch ret.Type {
	case DataNoRender:
		break
	case DataHTML:
		if m, ok := ret.Data.(map[string]interface{}); ok {
			// Set the copyright on all pages
			m["copyright"] = time.Now().Year()

			// FIXME: xsrf tokens

			// Add authentication information
			if cas.IsAuthenticated(r) {
				m["authenticated"] = cas.Username(r)
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
