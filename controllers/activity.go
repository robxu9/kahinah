package controllers

import (
	"net/http"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/robxu9/kahinah/data"
	"github.com/robxu9/kahinah/models"
	"github.com/robxu9/kahinah/util"

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
	ListId  uint64
	User    string
	Karma   int
	Comment string
	Time    time.Time
	URL     string
}

func ActivityJSONHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)

	page := to.Int64(r.FormValue("page"))
	if page <= 0 {
		page = 1
	}

	limit := to.Int64(r.FormValue("limit"))
	if limit <= 0 {
		limit = 50
	}

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.Karma))
	var karma []*models.Karma

	cnt, err := qt.Count()
	if err != nil {
		panic(err)
	}

	totalpages := cnt / 50
	if cnt%50 != 0 {
		totalpages++
	}

	if page > totalpages {
		page = totalpages
	}

	_, err = qt.Limit(limit, (page-1)*limit).OrderBy("-Time").All(&karma)
	if err != nil && err != orm.ErrNoRows {
		panic(err)
	}

	for _, v := range karma {
		o.LoadRelated(v, "List")
		o.LoadRelated(v, "User")
	}

	// render a better karma view
	var renderedKarma []*activityJSON
	for _, v := range karma {
		renderedKarma = append(renderedKarma, &activityJSON{
			ListId:  v.List.Id,
			User:    v.User.Username,
			Karma:   getTotalKarma(v.List.Id),
			Comment: v.Comment,
			Time:    v.Time,
			URL:     util.GetPrefixString("/b/" + to.String(v.List.Id)),
		})
	}

	dataRenderer.Data = map[string]interface{}{
		"totalpages": totalpages,
		"page":       page,
		"activities": renderedKarma,
	}
	dataRenderer.Type = data.DataJSON
}
