package controllers

import "net/http"

// RedirectHandler redirects paths without a trailing slash to
// those with a trailing slash...
func RedirectHandler(rw http.ResponseWriter, r *http.Request) {
	redirectTo := r.URL.String() + "/"
	http.Redirect(rw, r, redirectTo, http.StatusMovedPermanently)
}
