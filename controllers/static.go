package controllers

import (
	"net/http"
	"path"
	"strings"

	"goji.io/pattern"

	"github.com/robxu9/kahinah/data"

	"golang.org/x/net/context"
)

var (
	staticDir = http.Dir("static")
	indexFile = "index.html"
)

func StaticHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	file := pattern.Path(ctx)

	f, err := staticDir.Open(file)
	if err != nil {
		panic(ErrNotFound) // assume not found
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		panic(err) // this is weird
	}

	// try to serve index file
	if fi.IsDir() {
		// redirect if missing trailing slash
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(rw, r, r.URL.Path+"/", http.StatusFound)
			return
		}

		file = path.Join(file, indexFile)
		f, err = staticDir.Open(file)
		if err != nil {
			panic(ErrNotFound) // just hide its existence
		}
		defer f.Close()

		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			panic(ErrNotFound) // go away
		}
	}

	// okay, we don't need to render this.
	dataRender := data.FromContext(ctx)
	dataRender.Type = data.DataNoRender

	http.ServeContent(rw, r, file, fi.ModTime(), f)
}
