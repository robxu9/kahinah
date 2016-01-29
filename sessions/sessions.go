// Package sessions implements a session handler on top of kami with
// gorilla/sessions.
package sessions

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/zenazn/goji/web/mutil"

	gcontext "golang.org/x/net/context"
	"gopkg.in/guregu/kami.v1"
)

type key int

var sessionKey key
var storeKey key = 1

func NewStore(ctx gcontext.Context, secret []byte) gcontext.Context {
	cookieStore := sessions.NewCookieStore(secret)
	return gcontext.WithValue(ctx, storeKey, cookieStore)
}

func StoreFromContext(ctx gcontext.Context) *sessions.CookieStore {
	s, ok := ctx.Value(storeKey).(*sessions.CookieStore)
	if !ok {
		panic("unable to retrieve session store from context")
	}
	return s
}

func FromContext(ctx gcontext.Context) *sessions.Session {
	s, ok := ctx.Value(sessionKey).(*sessions.Session)
	if !ok {
		panic("unable to retrieve session from context")
	}
	return s
}

func SessionMiddleware(name string) kami.Middleware {
	return func(ctx gcontext.Context, rw http.ResponseWriter, r *http.Request) gcontext.Context {
		store := StoreFromContext(ctx)
		session, _ := store.Get(r, name)
		return gcontext.WithValue(ctx, sessionKey, session)
	}
}

func SessionAfterware(name string) kami.Afterware {
	return func(ctx gcontext.Context, wp mutil.WriterProxy, r *http.Request) gcontext.Context {
		session := FromContext(ctx)
		session.Save(r, wp)
		context.Clear(r)
		return gcontext.WithValue(ctx, sessionKey, nil)
	}
}
