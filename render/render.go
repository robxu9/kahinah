package render

import (
	"github.com/unrolled/render"
	"golang.org/x/net/context"
)

// via godoc for x/net/context: this allows us to distinguish types from
// different packages -- so no collisions will occur.
type key int

var renderKey key

// NewContext adds the renderer to the context.
func NewContext(ctx context.Context, r *render.Render) context.Context {
	return context.WithValue(ctx, renderKey, r)
}

// FromContext retrieves the renderer from the context.
func FromContext(ctx context.Context) *render.Render {
	r, ok := ctx.Value(renderKey).(*render.Render)
	if !ok {
		panic("unable to retrieve renderer from context")
	}
	return r
}
