package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/robxu9/kahinah/integration"
	"github.com/robxu9/kahinah/models"
	"io/ioutil"
	"log"
	"menteslibres.net/gosexy/to"
	"net/http"
	"sort"
	"time"
)

const (
	block_karma = 9999
)

var (
	maintainer_karma = to.Int64(beego.AppConfig.String("karma::maintainerkarma"))
	maintainer_hours = to.Int64(beego.AppConfig.String("karma::maintainerhours"))
)

type ByUpdateDate []*models.BuildList

func (b ByUpdateDate) Len() int {
	return len(b)
}
func (b ByUpdateDate) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b ByUpdateDate) Less(i, j int) bool {
	return b[i].Updated.Unix() > b[j].Updated.Unix()
}

type BuildsController struct {
	beego.Controller
}

func (this *BuildsController) Get() {
	this.Data["xsrf_token"] = this.XsrfToken()

	page, err := this.GetInt("page")
	if err != nil {
		page = 1
	} else if page <= 0 {
		page = 1
	}

	var packages []*models.BuildList

	o := orm.NewOrm()

	qt := o.QueryTable(new(models.BuildList))

	cnt, err := qt.Count()
	if err != nil {
		log.Println(err)
		this.Abort("500")
	}

	totalpages := cnt / 50
	if cnt%50 != 0 {
		totalpages++
	}

	if page > totalpages {
		page = totalpages
	}

	_, err = qt.Limit(50, (page-1)*50).OrderBy("-BuildDate").All(&packages)
	if err != nil && err != orm.ErrNoRows {
		log.Println(err)
		this.Abort("500")
	}

	for _, v := range packages {
		o.LoadRelated(v, "Submitter")
	}

	sort.Sort(ByUpdateDate(packages))

	this.Data["Title"] = "Builds"
	this.Data["Tab"] = 4
	this.Data["Packages"] = packages
	this.Data["PrevPage"] = page - 1
	this.Data["Page"] = page
	this.Data["NextPage"] = page + 1
	this.Data["Pages"] = totalpages
	this.TplNames = "builds_list.tpl"
}

type BuildController struct {
	beego.Controller
}

func (this *BuildController) Get() {
	this.Data["xsrf_token"] = this.XsrfToken()

	id := to.Uint64(this.Ctx.Input.Param(":id"))

	var pkg models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	err := qt.Filter("Id", id).One(&pkg)
	if err == orm.ErrNoRows {
		this.Abort("404")
	} else if err != nil {
		log.Println(err)
		this.Abort("500")
	}

	if pkg.Changelog != "" {
		resp, _ := http.Get(pkg.Changelog)
		if err == nil {
			defer resp.Body.Close()
			changelog, _ := ioutil.ReadAll(resp.Body)
			this.Data["Changelog"] = this.processChangelog(string(changelog))
		}
	}

	this.Data["Url"] = integration.Url(&pkg)

	// karma controls
	totalKarma := 0
	votes := make(map[string]bool)

	o.LoadRelated(&pkg, "Submitter")
	o.LoadRelated(&pkg, "Karma")
	o.LoadRelated(&pkg, "Packages")

	for _, karma := range pkg.Karma {
		o.LoadRelated(karma, "User")
		if karma.Vote == models.KARMA_UP {
			totalKarma++
			votes[karma.User.Email] = true
		} else if karma.Vote == models.KARMA_DOWN {
			totalKarma--
			votes[karma.User.Email] = false
		} else if karma.Vote == models.KARMA_MAINTAINER {
			totalKarma += int(maintainer_karma)
			votes[karma.User.Email] = true
		} else if karma.Vote == models.KARMA_BLOCK {
			totalKarma -= int(block_karma)
			votes[karma.User.Email] = false
		}
	}

	this.Data["Votes"] = votes
	this.Data["Karma"] = totalKarma

	user := models.IsLoggedIn(&this.Controller)
	if user != "" {
		kt := o.QueryTable(new(models.Karma))
		var userkarma models.Karma
		err = kt.Filter("User__Email", user).Filter("List__Id", id).One(&userkarma)
		if err != orm.ErrNoRows && err != nil {
			log.Println(err)
		} else if err == nil {
			if userkarma.Vote == models.KARMA_UP {
				this.Data["KarmaUpYes"] = true
			} else if userkarma.Vote == models.KARMA_MAINTAINER {
				this.Data["KarmaMaintainerYes"] = true
			} else {
				this.Data["KarmaDownYes"] = true
			}
		}

		if models.PermCheck(&this.Controller, PERMISSION_QA) {
			this.Data["QAControls"] = true
		}

	}

	// karma controls end

	this.Data["Title"] = "Build " + to.String(id) + ": " + pkg.Name
	if pkg.Status == models.STATUS_TESTING {
		this.Data["Tab"] = 1
		if user != "" {
			this.Data["KarmaControls"] = true
			if pkg.Submitter != nil && pkg.Submitter.Email == user {
				if time.Since(pkg.BuildDate).Hours() >= float64(maintainer_hours) {
					this.Data["MaintainerControls"] = true
				}
			}
		}
	} else if pkg.Status == models.STATUS_PUBLISHED {
		this.Data["Tab"] = 2
	} else if pkg.Status == models.STATUS_REJECTED {
		this.Data["Tab"] = 3
	} else {
		this.Data["Tab"] = 4
	}
	this.Data["Package"] = pkg
	this.TplNames = "build.tpl"
}

func (this *BuildController) Post() {
	id := to.Uint64(this.Ctx.Input.Param(":id"))

	postType := this.GetString("type")
	if postType != "Up" && postType != "Down" && postType != "Maintainer" && postType != "QABlock" {
		this.Abort("400")
	}

	user := models.IsLoggedIn(&this.Controller)
	if user == "" {
		this.Abort("403") // MUST be logged in
	}

	var pkg models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	err := qt.Filter("Id", id).One(&pkg)
	if err == orm.ErrNoRows {
		this.Abort("404")
	} else if err != nil {
		log.Println(err)
		this.Abort("500")
	}

	o.LoadRelated(&pkg, "Submitter")

	if postType == "Maintainer" {
		if pkg.Submitter.Email != user {
			this.Abort("403")
		} else {
			if time.Since(pkg.BuildDate).Hours() < float64(maintainer_hours) { // week
				this.Abort("400")
			}
		}
	} else if postType == "QABlock" {
		models.PermAbortCheck(&this.Controller, PERMISSION_QA)
	} else {
		// whitelist stuff
		if Whitelist {
			perm := models.PermCheck(&this.Controller, PERMISSION_WHITELIST)
			if !perm {
				flash := beego.NewFlash()
				flash.Warning("Sorry, the whitelist is on and you are not allowed to vote.")
				flash.Store(&this.Controller)
				this.Get()
				return
			}
		}
	}

	kt := o.QueryTable(new(models.Karma))

	var userkarma models.Karma
	err = kt.Filter("List__Id", id).Filter("User__Email", user).One(&userkarma)
	if err != orm.ErrNoRows && err != nil {
		log.Println(err)
	} else if err == nil { // already has entry
		if postType == "Up" {
			if userkarma.Vote == models.KARMA_UP {
				o.Delete(&userkarma)
			} else {
				userkarma.Vote = models.KARMA_UP
				o.Update(&userkarma)
			}
		} else if postType == "Maintainer" {
			if userkarma.Vote == models.KARMA_MAINTAINER {
				o.Delete(&userkarma)
			} else {
				userkarma.Vote = models.KARMA_MAINTAINER
				o.Update(&userkarma)
			}
		} else if postType == "QABlock" {
			if userkarma.Vote == models.KARMA_BLOCK {
				o.Delete(&userkarma)
			} else {
				userkarma.Vote = models.KARMA_BLOCK
				o.Update(&userkarma)
			}
		} else {
			if userkarma.Vote == models.KARMA_DOWN {
				o.Delete(&userkarma)
			} else {
				userkarma.Vote = models.KARMA_DOWN
				o.Update(&userkarma)
			}
		}
	} else {
		userkarma.List = &pkg
		userkarma.User = models.FindUser(user)
		if postType == "Up" {
			userkarma.Vote = models.KARMA_UP
		} else if postType == "Maintainer" {
			userkarma.Vote = models.KARMA_MAINTAINER
		} else if postType == "QABlock" {
			userkarma.Vote = models.KARMA_BLOCK
		} else {
			userkarma.Vote = models.KARMA_DOWN
		}
		o.Insert(&userkarma)
	}

	karmaup, _ := kt.Filter("List__Id", id).Filter("Vote", models.KARMA_UP).Count()
	karmadown, _ := kt.Filter("List__Id", id).Filter("Vote", models.KARMA_DOWN).Count()
	karmamaintainer, _ := kt.Filter("List__Id", id).Filter("Vote", models.KARMA_MAINTAINER).Count()
	karmablock, _ := kt.Filter("List__Id", id).Filter("Vote", models.KARMA_BLOCK).Count()

	karmaTotal := karmaup - karmadown + (maintainer_karma * karmamaintainer) - (block_karma * karmablock)

	upthreshold, err := beego.AppConfig.Int64("karma::upperkarma")
	if err != nil {
		panic(err)
	}

	downthreshold, err := beego.AppConfig.Int64("karma::lowerkarma")
	if err != nil {
		panic(err)
	}

	if karmaTotal >= upthreshold {
		pkg.Status = models.STATUS_PUBLISHED
		o.Update(&pkg)
		go integration.Publish(&pkg)
	} else if karmaTotal <= downthreshold {
		pkg.Status = models.STATUS_REJECTED
		o.Update(&pkg)
		go integration.Reject(&pkg)
	}

	this.Get()
}

func (this *BuildController) processChangelog(changelog string) string {
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
