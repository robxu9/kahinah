package controllers

import (
	"net/http"

	"gopkg.in/cas.v1"

	"github.com/knq/sessionmw"
	"github.com/robxu9/kahinah/data"
	"github.com/robxu9/kahinah/render"

	"golang.org/x/net/context"
)

const (
	CASReferrer = "_cas_refer"
)

func UserLoginHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	if !cas.IsAuthenticated(r) {
		// save the referrer
		sessionmw.Set(ctx, CASReferrer, r.Referer())

		// shut off rendering
		dataRenderer := data.FromContext(ctx)
		dataRenderer.Type = data.DataNoRender

		// and redirect
		cas.RedirectToLogin(rw, r)
	} else {
		// get the referrer
		referrer, has := sessionmw.Get(ctx, CASReferrer)
		sessionmw.Delete(ctx, CASReferrer)

		// shut off rendering
		dataRenderer := data.FromContext(ctx)
		dataRenderer.Type = data.DataNoRender

		// and redirect
		if !has {
			http.Redirect(rw, r, render.ConvertURL("/"), http.StatusTemporaryRedirect)
		} else {
			http.Redirect(rw, r, referrer.(string), http.StatusTemporaryRedirect)
		}
	}
}

func UserLogoutHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	// shut off rendering
	dataRenderer := data.FromContext(ctx)
	dataRenderer.Type = data.DataNoRender

	// Destroy the session
	sessionmw.Destroy(ctx, rw)

	// CAS logouts are always one-way.
	cas.RedirectToLogout(rw, r)
}
