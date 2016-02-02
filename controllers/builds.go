package controllers

import (
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"

	"goji.io/pat"

	"golang.org/x/net/context"

	"github.com/astaxie/beego/orm"
	"github.com/robxu9/kahinah/conf"
	"github.com/robxu9/kahinah/data"
	"github.com/robxu9/kahinah/integration"
	"github.com/robxu9/kahinah/models"
	"github.com/robxu9/kahinah/util"
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
		"Loc":      1,
		"Tab":      4,
		"Packages": packages,
		"PrevPage": page - 1,
		"Page":     page,
		"NextPage": page + 1,
		"Pages":    totalpages,
	}
	dataRenderer.Template = "builds/builds_list"
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

	cnt, err := qt.Filter("status", models.STATUS_REJECTED).Count()
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

	_, err = qt.Limit(50, (page-1)*50).OrderBy("-Updated").Filter("status", models.STATUS_REJECTED).All(&packages)
	if err != nil && err != orm.ErrNoRows {
		panic(err)
	}

	for _, v := range packages {
		o.LoadRelated(v, "Submitter")
	}

	sort.Sort(sortByUpdateDate(packages))

	dataRenderer.Data = map[string]interface{}{
		"Title":    "Rejected",
		"Loc":      1,
		"Tab":      3,
		"Packages": packages,
		"PrevPage": page - 1,
		"Page":     page,
		"NextPage": page + 1,
		"Pages":    totalpages,
	}
	dataRenderer.Template = "builds/builds_list"
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

	cnt, err := qt.Filter("status", models.STATUS_PUBLISHED).Count()
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

	_, err = qt.Limit(50, (page-1)*50).OrderBy("-Updated").Filter("status", models.STATUS_PUBLISHED).All(&packages)
	if err != nil && err != orm.ErrNoRows {
		panic(err)
	}

	for _, v := range packages {
		o.LoadRelated(v, "Submitter")
	}

	sort.Sort(sortByUpdateDate(packages))

	dataRenderer.Data = map[string]interface{}{
		"Title":    "Published",
		"Loc":      1,
		"Tab":      2,
		"Packages": packages,
		"PrevPage": page - 1,
		"Page":     page,
		"NextPage": page + 1,
		"Pages":    totalpages,
	}
	dataRenderer.Template = "builds/builds_list"
}

// TestingHandler shows the list of builds that have yet to be approved, paginated.
func TestingHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)

	var packages []*models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	num, err := qt.Filter("status", models.STATUS_TESTING).All(&packages)
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
		"Loc":      1,
		"Tab":      1,
		"Packages": packages,
		"PkgKarma": pkgkarma,
		"Entries":  num,
	}
	dataRenderer.Template = "builds/generic_list"
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

	if pkg.Changelog != "" {
		resp, err := http.Get(pkg.Changelog)
		if err == nil {
			defer resp.Body.Close()
			changelog, _ := ioutil.ReadAll(resp.Body)
			toRender["Changelog"] = processChangelog(string(changelog))
		} else {
			toRender["Changelog"] = "Failed to retrieve changelog: " + err.Error()
		}
	}

	toRender["Commits"] = integration.Commits(&pkg)

	toRender["Url"] = integration.Url(&pkg)

	o.LoadRelated(&pkg, "Submitter")
	o.LoadRelated(&pkg, "Packages")

	// karma controls
	totalKarma := getTotalKarma(id) // get total karma

	var votes []util.Pair // *models.Karma, int

	// load karma totals
	var inOrder []*models.Karma
	kt := o.QueryTable(new(models.Karma))
	kt.Filter("List__Id", id).OrderBy("Time").All(&inOrder)

	// only count most recent votes
	for _, v := range inOrder {
		o.LoadRelated(v, "User")

		pair := util.Pair{}
		pair.Key = v

		switch v.Vote {
		case models.KARMA_UP:
			pair.Value = 1
		case models.KARMA_DOWN:
			pair.Value = 2
		case models.KARMA_MAINTAINER:
			pair.Value = 1
		case models.KARMA_BLOCK:
			pair.Value = 2
		case models.KARMA_PUSH:
			pair.Value = 1
		case models.KARMA_NONE:
			if v.Comment != "" {
				pair.Value = 0
			} else {
				continue // no karma and no comment? useless
			}
		}

		votes = append(votes, pair)
	}

	toRender["Votes"] = votes
	toRender["Karma"] = totalKarma

	toRender["UserVote"] = 0
	user := Authenticated(r)
	if user != "" {
		kt := o.QueryTable(new(models.Karma))
		var userkarma models.Karma
		err = kt.Filter("User__Email", user).Filter("List__Id", id).OrderBy("-Time").Limit(1).One(&userkarma)
		if err != orm.ErrNoRows && err != nil {
			log.Println(err)
		} else if err == nil {
			if userkarma.Vote == models.KARMA_UP {
				toRender["UserVote"] = 1
			} else if userkarma.Vote == models.KARMA_MAINTAINER {
				toRender["UserVote"] = 2
			} else if userkarma.Vote == models.KARMA_DOWN {
				toRender["UserVote"] = -1
			}

			toRender["KarmaCommentPrev"] = userkarma.Comment
		}

		if PermCheck(r, PermissionQA) {
			toRender["QAControls"] = true
		}

	}

	// karma controls end

	toRender["Title"] = "Build " + to.String(id) + ": " + pkg.Name
	toRender["Loc"] = 1
	if pkg.Status == models.STATUS_TESTING {
		toRender["Tab"] = 1
		if user != "" {
			toRender["KarmaControls"] = true
			if pkg.Submitter != nil && pkg.Submitter.Email == user {
				toRender["MaintainerControls"] = true
				toRender["MaintainerHoursNeeded"] = maintainerHours
				if time.Since(pkg.BuildDate).Hours() >= float64(maintainerHours) {
					toRender["MaintainerTime"] = true
					delete(toRender, "MaintainerHoursNeeded")
				}
			}
		}
	} else if pkg.Status == models.STATUS_PUBLISHED {
		toRender["Tab"] = 2
	} else if pkg.Status == models.STATUS_REJECTED {
		toRender["Tab"] = 3
	} else {
		toRender["Tab"] = 4
	}
	toRender["Package"] = pkg

	dataRenderer.Data = toRender
	dataRenderer.Template = "builds/build"
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

		upperThreshold := conf.Config.GetDefault("karma.upperKarma", 3).(int)
		lowerThreshold := conf.Config.GetDefault("karma.lowerKarma", -3).(int)

		if karmaTotal >= upperThreshold {
			pkg.Status = models.STATUS_PUBLISHED
			o.Update(&pkg)
			go integration.Publish(&pkg)
		} else if karmaTotal <= lowerThreshold {
			pkg.Status = models.STATUS_REJECTED
			o.Update(&pkg)
			go integration.Reject(&pkg)
		}
	}

	dataRenderer.Type = data.DataNoRender
	http.Redirect(rw, r, r.URL.String(), http.StatusTemporaryRedirect)
}

func getTotalKarma(id uint64) int {
	o := orm.NewOrm()
	kt := o.QueryTable(new(models.Karma))

	var karma []*models.Karma
	kt.Filter("List__Id", id).OrderBy("-Time").All(&karma)

	set := util.NewSet()
	totalKarma := 0

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
			totalKarma += int(maintainerKarma)
		case models.KARMA_BLOCK:
			totalKarma -= int(blockKarma)
		case models.KARMA_PUSH:
			totalKarma += int(pushKarma)
		}

		set.Add(v.User.Email)
	}

	return totalKarma
}

func processChangelog(changelog string) string {
	toreturn := ""
	open := true
	for _, c := range changelog {
		if c == '<' {
			toreturn += string(c)
			toreturn += "email hidden"
			open = false
		} else if c == '>' {
			toreturn += string(c)
			open = true
		} else if open {
			toreturn += string(c)
		}
	}
	return toreturn
}
