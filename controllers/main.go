package controllers

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/robxu9/kahinah/data"
	"github.com/russross/blackfriday"
	"golang.org/x/net/context"
)

// MainHandler shows the main page.
func MainHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)

	stat, err := os.Stat("news.md")
	var time time.Time
	if err == nil {
		time = stat.ModTime()
	}

	bte, err := ioutil.ReadFile("news.md")
	markdown := []byte("_Couldn't retrieve the latest news._")
	if err == nil {
		markdown = bte
	}

	output := blackfriday.MarkdownCommon(markdown)

	dataRenderer.Data = map[string]interface{}{
		"Title": "Main",
		"News":  template.HTML(bluemonday.UGCPolicy().SanitizeBytes(output)),
		"Time":  time,
		"Nav":   0,
	}
	dataRenderer.Template = "index"
}
