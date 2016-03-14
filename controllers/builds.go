package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"goji.io/pat"

	"golang.org/x/net/context"

	"github.com/astaxie/beego/orm"
	"github.com/jinzhu/gorm"
	"github.com/knq/sessionmw"
	"github.com/microcosm-cc/bluemonday"
	"github.com/robxu9/kahinah/conf"
	"github.com/robxu9/kahinah/data"
	"github.com/robxu9/kahinah/integration"
	"github.com/robxu9/kahinah/log"
	"github.com/robxu9/kahinah/models"
	"github.com/robxu9/kahinah/util"
	"github.com/russross/blackfriday"
	"menteslibres.net/gosexy/to"
)

const (
	blockKarma = 9999
	pushKarma  = 9999
)

var (
	maintainerKarma = to.Int64(conf.Config.Get("karma.maintainerKarma"))
	maintainerHours = to.Int64(conf.Config.Get("karma.maintainerHours"))
)

//
// --------------------------------------------------------------------
// LISTS
// --------------------------------------------------------------------
//

// ListsAPIHandler shows the collection of lists, with filters, paginated, JSON'ified.
func ListsAPIHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)
	dataRenderer.Type = data.DataJSON

	platform := r.FormValue("platform")
	channel := r.FormValue("channel")
	status := r.FormValue("status")

	limit := int(to.Int64(r.FormValue("limit")))
	if limit <= 0 {
		limit = 50 // reasonable
	}

	page := int(to.Int64(r.FormValue("page")))
	if page <= 0 {
		page = 1
	}

	baseDB := models.DB.Model(&models.List{})
	if platform != "" {
		baseDB = baseDB.Where("platform = ?", platform)
	}

	if channel != "" {
		baseDB = baseDB.Where("channel = ?", channel)
	}

	if status == models.ListRunning || status == models.ListPending || status == models.ListSuccess || status == models.ListFailed {
		baseDB = baseDB.Where("stage_result = ?", status)
	} else if status != "" {
		panic(ErrBadRequest)
	}

	var cnt int
	if err := baseDB.Count(&cnt).Error; err != nil {
		panic(err)
	}

	totalpages := cnt / limit
	if cnt%limit != 0 {
		totalpages++
	}

	if page > totalpages {
		page = totalpages
	}

	var packages []*models.List
	if err := baseDB.Limit(limit).Offset((page - 1) * limit).Order("updated_at desc").Find(&packages).Error; err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	dataRenderer.Data = map[string]interface{}{
		"lists": packages,
		"pages": map[string]interface{}{
			"prev":    page - 1,
			"current": page,
			"next":    page + 1,
			"total":   totalpages,
		},
	}
}

// ListsHandler shows the collection of lists (HTML).
func ListsHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)

	// filters that will be passed on by Vue.js to the API
	platform := r.FormValue("platform")
	channel := r.FormValue("channel")
	status := r.FormValue("status")
	limit := r.FormValue("limit")
	page := r.FormValue("page")

	dataRenderer.Data = map[string]interface{}{
		"Title":    "Lists",
		"Nav":      2,
		"Platform": platform,
		"Channel":  channel,
		"Status":   status,
		"Limit":    limit,
		"Page":     page,
	}
	dataRenderer.Template = "i/list"
}

//
// --------------------------------------------------------------------
// INDIVIDUAL BUILD
// --------------------------------------------------------------------
//

// BuildGetHandler displays build information for a specific build
func BuildGetHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)

	toRender := map[string]interface{}{}

	id := to.Uint64(pat.Param(ctx, "id"))

	var pkg models.List
	if err := models.DB.Where("id = ?", id).First(&pkg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(ErrNotFound)
		} else {
			panic(err)
		}
	}

	toRender["Title"] = "Build " + to.String(id) + ": " + pkg.Name
	toRender["Nav"] = 2
	toRender["ID"] = to.String(id)

	dataRenderer.Data = toRender
	dataRenderer.Template = "builds/build"
}

type buildGetJSONKarma struct {
	User    string
	Karma   string
	Comment string
	Time    time.Time
}

type buildGetJSON struct {
	ID        uint64
	Platform  string
	Channel   string   // maps to repo
	Arch      []string // maps to architectures
	Name      string
	Submitter string
	Type      string // update type
	Status    string // testing/published/rejected
	Artifacts []*models.ListArtifact
	Links     []*models.ListLink
	BuildDate time.Time
	Updated   time.Time
	Activity  []*buildGetJSONKarma
	Diff      string
	Advisory  *models.Advisory

	TotalKarma int64
	User       string
	UserIsQA   bool
	Maintainer bool

	MaintainerKarma int64
	PushKarma       int
	BlockKarma      int
	Acceptable      bool
	Rejectable      bool
}

// BuildGetJSONHandler displays build information in JSON for a specific build.
func BuildGetJSONHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)
	id := to.Uint64(pat.Param(ctx, "id"))

	// load the requested build list
	var pkg models.List
	if err := models.DB.Where("id = ?", id).First(&pkg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(ErrNotFound)
		} else {
			panic(err)
		}
	}

	o.LoadRelated(&pkg, "Submitter")
	o.LoadRelated(&pkg, "Packages")
	o.LoadRelated(&pkg, "Links")
	o.LoadRelated(&pkg, "Advisory")
	o.LoadRelated(&pkg, "Karma")

	// load karma
	totalKarma := getTotalKarma(id) // get total karma
	var renderedKarma []*buildGetJSONKarma
	for _, v := range pkg.Karma {
		o.LoadRelated(v, "User")
		renderedKarma = append(renderedKarma, &buildGetJSONKarma{
			User:    v.User.Username,
			Karma:   v.Vote,
			Comment: string(bluemonday.UGCPolicy().SanitizeBytes(blackfriday.MarkdownCommon([]byte(v.Comment)))),
			Time:    v.Time,
		})
	}

	user := Authenticated(r)
	isQA := PermCheck(r, PermissionQA)
	maintainerAllowed := pkg.Submitter.Username == user && time.Since(pkg.BuildDate).Hours() >= float64(maintainerHours)

	// check if we can accept or reject
	acceptable := false
	rejectable := false

	if pkg.Status == models.StatusTesting {
		upperThreshold := conf.Config.GetDefault("karma.upperKarma", 3).(int64)
		lowerThreshold := conf.Config.GetDefault("karma.lowerKarma", -3).(int64)

		if totalKarma >= upperThreshold {
			acceptable = true
		} else if totalKarma <= lowerThreshold {
			rejectable = true
		}
	}

	// render the data in a nice way

	dataRenderer.Data = &buildGetJSON{
		ID:        pkg.Id,
		Platform:  pkg.Platform,
		Channel:   pkg.Repo,
		Arch:      strings.Split(pkg.Architecture, ";"),
		Name:      pkg.Name,
		Submitter: pkg.Submitter.Username,
		Type:      pkg.Type,
		Status:    pkg.Status,
		Artifacts: pkg.Packages,
		Links:     pkg.Links,
		BuildDate: pkg.BuildDate,
		Updated:   pkg.Updated,
		Activity:  renderedKarma,
		Diff:      pkg.Diff,
		Advisory:  pkg.Advisory,

		TotalKarma: totalKarma,
		User:       user,
		UserIsQA:   isQA,
		Maintainer: maintainerAllowed,

		MaintainerKarma: maintainerKarma,
		PushKarma:       pushKarma,
		BlockKarma:      blockKarma,
		Acceptable:      acceptable,
		Rejectable:      rejectable,
	}
	dataRenderer.Type = data.DataJSON
}

// BuildPostHandler handles post actions that occur.
func BuildPostHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	// check for authentication
	user := MustAuthenticate(r)

	// load parameters
	dataRenderer := data.FromContext(ctx)

	id := to.Uint64(pat.Param(ctx, "id"))
	action := r.FormValue("type")
	comment := r.FormValue("comment")

	// find the build list
	var pkg models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	err := qt.Filter("Id", id).One(&pkg)
	if err == orm.ErrNoRows {
		panic(ErrNotFound)
	} else if err != nil {
		panic(err)
	}

	o.LoadRelated(&pkg, "Submitter")

	// FIXME: check for whitelist

	switch action {
	case "Neutral":
	case "Up":
	case "Down":
	case "Maintainer":
		if pkg.Submitter.Username != user {
			panic(ErrForbidden)
		}
		if time.Since(pkg.BuildDate).Hours() < float64(maintainerHours) {
			panic(ErrBadRequest)
		}
	case "QABlock":
		fallthrough
	case "QAPush":
		PermAbortCheck(r, PermissionQA)
	case "Accept":
	case "Reject":
	default:
		panic(ErrBadRequest)
	}

	var userkarma models.Karma

	userkarma.List = &pkg
	userkarma.User = models.FindUser(user)

	switch action {
	case "Up":
		userkarma.Vote = models.KARMA_UP
	case "Maintainer":
		userkarma.Vote = models.KARMA_MAINTAINER
	case "QABlock":
		userkarma.Vote = models.KARMA_BLOCK
	case "QAPush":
		userkarma.Vote = models.KARMA_PUSH
	case "Down":
		userkarma.Vote = models.KARMA_DOWN
	default:
		userkarma.Vote = models.KARMA_NONE
	}

	userkarma.Comment = comment
	o.Insert(&userkarma)

	if action == "Accept" || action == "Reject" {
		karmaTotal := getTotalKarma(id)

		upperThreshold := conf.Config.GetDefault("karma.upperKarma", 3).(int64)
		lowerThreshold := conf.Config.GetDefault("karma.lowerKarma", -3).(int64)

		if karmaTotal >= upperThreshold {
			pkg.Status = models.StatusPublished
			o.Update(&pkg)
			go func() {
				err := integration.Accept(&pkg)
				if err != nil {
					log.Logger.Critical("Unable to accept update %v: %v", id, err)
				}
			}()
		} else if karmaTotal <= lowerThreshold {
			pkg.Status = models.StatusRejected
			o.Update(&pkg)

			go func() {
				err := integration.Reject(&pkg)
				if err != nil {
					log.Logger.Critical("Unable to reject update %v: %v", id, err)
				}
			}()
		}
	}

	dataRenderer.Type = data.DataNoRender
	sessionmw.Set(ctx, data.FlashInfo, fmt.Sprintf("Committed \"%v\" with comment \"%v\".", action, comment))
	http.Redirect(rw, r, r.URL.String(), http.StatusTemporaryRedirect)
}

func getTotalKarma(id uint64) int64 {
	o := orm.NewOrm()
	kt := o.QueryTable(new(models.Karma))

	var karma []*models.Karma
	kt.Filter("List__Id", id).OrderBy("-Time").All(&karma)

	set := util.NewSet()
	var totalKarma int64

	// only count most recent votes
	for _, v := range karma {
		o.LoadRelated(v, "User")

		if set.Contains(v.User.Username) {
			continue // we've already counted this person's most recent vote
		}

		switch v.Vote {
		case models.KARMA_UP:
			totalKarma++
		case models.KARMA_DOWN:
			totalKarma--
		case models.KARMA_MAINTAINER:
			totalKarma += maintainerKarma
		case models.KARMA_BLOCK:
			totalKarma -= blockKarma
		case models.KARMA_PUSH:
			totalKarma += pushKarma
		}

		set.Add(v.User.Username)
	}

	return totalKarma
}
