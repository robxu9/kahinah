package controllers

import (
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/microcosm-cc/bluemonday"
	"github.com/robxu9/kahinah/data"
	"github.com/russross/blackfriday"
	"golang.org/x/net/context"
)

// MainHandler shows the main page.
func MainHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)

	bte, err := ioutil.ReadFile("news.md")
	markdown := []byte("_Couldn't retrieve the latest news._")
	if err == nil {
		markdown = bte
	}

	output := blackfriday.MarkdownCommon(markdown)

	dataRenderer.Data = map[string]interface{}{
		"Title": "Main",
		"News":  template.HTML(bluemonday.UGCPolicy().SanitizeBytes(output)),
		"Loc":   0,
	}
	dataRenderer.Template = "index"
}
