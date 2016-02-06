package controllers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"goji.io/pat"

	"golang.org/x/net/context"

	"github.com/astaxie/beego/orm"
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

// Sorters

type sortByUpdateDate []*models.BuildList

func (b sortByUpdateDate) Len() int {
	return len(b)
}
func (b sortByUpdateDate) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b sortByUpdateDate) Less(i, j int) bool {
	return b[i].Updated.Unix() > b[j].Updated.Unix()
}

type sortByBuildDate []*models.BuildList

func (b sortByBuildDate) Len() int {
	return len(b)
}
func (b sortByBuildDate) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b sortByBuildDate) Less(i, j int) bool {
	return b[i].BuildDate.Unix() > b[j].BuildDate.Unix()
}

//
// --------------------------------------------------------------------
// LISTS
// --------------------------------------------------------------------
//

// BuildsHandler shows the list of all available builds, paginated.
func BuildsHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)

	page := to.Int64(r.FormValue("page"))
	if page <= 0 {
		page = 1
	}

	var packages []*models.BuildList

	o := orm.NewOrm()

	qt := o.QueryTable(new(models.BuildList))

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

	_, err = qt.Limit(50, (page-1)*50).OrderBy("-Updated").All(&packages)
	if err != nil && err != orm.ErrNoRows {
		panic(err)
	}

	for _, v := range packages {
		o.LoadRelated(v, "Submitter")
	}

	sort.Sort(sortByUpdateDate(packages))

	dataRenderer.Data = map[string]interface{}{
		"Title":    "Builds",
		"Nav":      5,
		"Packages": packages,
		"PrevPage": page - 1,
		"Page":     page,
		"NextPage": page + 1,
		"Pages":    totalpages,
	}
	dataRenderer.Template = "i/list_paginated"
}

// RejectedHandler shows the list of all rejected builds, paginated.
func RejectedHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)

	page := to.Int64(r.FormValue("page"))
	if page <= 0 {
		page = 1
	}

	var packages []*models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	cnt, err := qt.Filter("status", models.StatusRejected).Count()
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

	_, err = qt.Limit(50, (page-1)*50).OrderBy("-Updated").Filter("status", models.StatusRejected).All(&packages)
	if err != nil && err != orm.ErrNoRows {
		panic(err)
	}

	for _, v := range packages {
		o.LoadRelated(v, "Submitter")
	}

	sort.Sort(sortByUpdateDate(packages))

	dataRenderer.Data = map[string]interface{}{
		"Title":    "Rejected",
		"Nav":      4,
		"Packages": packages,
		"PrevPage": page - 1,
		"Page":     page,
		"NextPage": page + 1,
		"Pages":    totalpages,
	}
	dataRenderer.Template = "i/list_paginated"
}

// PublishedHandler shows the list of all published builds, paginated.
func PublishedHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)

	filterPlatform := r.FormValue("platform")

	page := to.Int64(r.FormValue("page"))
	if page <= 0 {
		page = 1
	}

	var packages []*models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	if filterPlatform != "" {
		qt = qt.Filter("platform", filterPlatform)
	}

	cnt, err := qt.Filter("status", models.StatusPublished).Count()
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

	_, err = qt.Limit(50, (page-1)*50).OrderBy("-Updated").Filter("status", models.StatusPublished).All(&packages)
	if err != nil && err != orm.ErrNoRows {
		panic(err)
	}

	for _, v := range packages {
		o.LoadRelated(v, "Submitter")
	}

	sort.Sort(sortByUpdateDate(packages))

	dataRenderer.Data = map[string]interface{}{
		"Title":    "Accepted",
		"Nav":      3,
		"Packages": packages,
		"PrevPage": page - 1,
		"Page":     page,
		"NextPage": page + 1,
		"Pages":    totalpages,
	}
	dataRenderer.Template = "i/list_paginated"
}

// TestingHandler shows the list of builds that have yet to be approved, paginated.
func TestingHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)

	var packages []*models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	num, err := qt.Filter("status", models.StatusTesting).All(&packages)
	if err != nil && err != orm.ErrNoRows {
		panic(err)
	}

	pkgkarma := make(map[string]string)

	for _, v := range packages {
		totalKarma := getTotalKarma(v.Id)

		pkgkarma[to.String(v.Id)] = to.String(totalKarma)

		o.LoadRelated(v, "Submitter")
	}

	sort.Sort(sortByBuildDate(packages))

	dataRenderer.Data = map[string]interface{}{
		"Title":    "Testing",
		"Nav":      2,
		"Packages": packages,
		"PkgKarma": pkgkarma,
		"Entries":  num,
	}
	dataRenderer.Template = "i/list_sortable"
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

	var pkg models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	err := qt.Filter("Id", id).One(&pkg)
	if err == orm.ErrNoRows {
		panic(ErrNotFound)
	} else if err != nil {
		panic(err)
	}

	toRender["Title"] = "Build " + to.String(id) + ": " + pkg.Name
	if pkg.Status == models.StatusTesting {
		toRender["Nav"] = 2
	} else if pkg.Status == models.StatusPublished {
		toRender["Nav"] = 3
	} else if pkg.Status == models.StatusRejected {
		toRender["Nav"] = 4
	} else {
		toRender["Nav"] = 5
	}
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
	Artifacts []*models.BuildListPkg
	Links     []*models.BuildListLink
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
