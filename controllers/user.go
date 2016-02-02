package controllers

import (
	"net/http"

	"gopkg.in/cas.v1"

	"github.com/robxu9/kahinah/data"
	"github.com/robxu9/kahinah/util"

	"golang.org/x/net/context"
)

func UserLoginHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	if !cas.IsAuthenticated(r) {

		// shut off rendering
		dataRenderer := data.FromContext(ctx)
		dataRenderer.Type = data.DataNoRender

		// and redirect
		cas.RedirectToLogin(rw, r)
	} else {

		// shut off rendering
		dataRenderer := data.FromContext(ctx)
		dataRenderer.Type = data.DataNoRender

		// and redirect
		http.Redirect(rw, r, util.GetPrefixString("/"), http.StatusTemporaryRedirect)
	}

}

func UserLogoutHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	// CAS logouts are always one-way.
	cas.RedirectToLogout(rw, r)
}
