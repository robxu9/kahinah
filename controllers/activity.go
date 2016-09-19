package controllers

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/microcosm-cc/bluemonday"
	"github.com/robxu9/kahinah/data"
	"github.com/robxu9/kahinah/models"
	"github.com/robxu9/kahinah/render"
	"github.com/russross/blackfriday"

	"menteslibres.net/gosexy/to"

	"golang.org/x/net/context"
)

func ActivityHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)

	// we more or less have Vue.js show recent activity, so just render the template
	dataRenderer.Data = map[string]interface{}{
		"Title": "Recent Activity",
		"Nav":   1,
	}
	dataRenderer.Template = "i/activity"
}

type activityJSON struct {
	ListId  uint
	User    string
	Comment string
	Time    time.Time
	URL     string
}

func ActivityJSONHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)

	page := int(to.Int64(r.FormValue("page")))
	if page <= 0 {
		page = 1
	}

	limit := int(to.Int64(r.FormValue("limit")))
	if limit <= 0 {
		limit = 50
	}

	var cnt int
	if err := models.DB.Model(&models.ListActivity{}).Count(&cnt).Error; err != nil {
		panic(err)
	}

	totalpages := cnt / 50
	if cnt%50 != 0 {
		totalpages++
	}

	if page > totalpages {
		page = totalpages
	}

	var activities []models.ListActivity
	if err := models.DB.Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&activities).Error; err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	// render a better karma view
	var rendered []*activityJSON
	for _, v := range activities {
		// load the username...
		rendered = append(rendered, &activityJSON{
			ListId:  v.ListID,
			User:    models.FindUserByID(v.UserID).Username,
			Comment: string(bluemonday.UGCPolicy().SanitizeBytes(blackfriday.MarkdownCommon([]byte(v.Activity)))),
			Time:    v.CreatedAt,
			URL:     render.ConvertURL("/b/" + to.String(v.ListID)),
		})
	}

	dataRenderer.Data = map[string]interface{}{
		"totalpages": totalpages,
		"page":       page,
		"activities": rendered,
	}
	dataRenderer.Type = data.DataJSON
}
