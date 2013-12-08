package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/robxu9/kahinah/integration"
	"github.com/robxu9/kahinah/models"
	"github.com/robxu9/kahinah/util"
	"io/ioutil"
	"log"
	"menteslibres.net/gosexy/to"
	"net/http"
	"sort"
	"strings"
)

type BuildsController struct {
	beego.Controller
}

func (this *BuildsController) Get() {
	page, err := this.GetInt("page")
	if err != nil {
		page = 1
	} else if page <= 0 {
		page = 1
	}

	var packages []models.BuildList

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

	_, err = qt.Limit(50, (page-1)*50).All(&packages)
	if err != nil && err != orm.ErrNoRows {
		log.Println(err)
		this.Abort("500")
	}

	sort.Sort(ByTimeTP(packages))

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
	id := to.Uint64(this.Ctx.Input.Param(":buildid"))

	var pkg models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	err := qt.Filter("ListId", id).One(&pkg)
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

	// karma controls
	kt := o.QueryTable(new(models.Karma))

	var upkarma []models.Karma
	karmaup, _ := kt.Filter("ListId", id).Filter("Vote", models.KARMA_UP).All(&upkarma)

	this.Data["YayVotes"] = upkarma

	var downkarma []models.Karma
	karmadown, _ := kt.Filter("ListId", id).Filter("Vote", models.KARMA_DOWN).All(&downkarma)

	this.Data["NayVotes"] = downkarma

	this.Data["Karma"] = karmaup - karmadown

	user := util.IsLoggedIn(this.Controller)
	if user != "" {
		var userkarma models.Karma
		err = kt.Filter("ListId", id).One(&userkarma)
		if err != orm.ErrNoRows && err != nil {
			log.Println(err)
		} else if err == nil {
			if userkarma.Vote == models.KARMA_UP {
				this.Data["KarmaUpYes"] = true
			} else {
				this.Data["KarmaDownYes"] = true
			}
		}

	}

	// karma controls end

	this.Data["Title"] = "Build " + to.String(id) + ": " + pkg.Name
	if pkg.Status == models.STATUS_TESTING {
		this.Data["Tab"] = 1
		this.Data["Header"] = "Testing"
		if user != "" {
			this.Data["KarmaControls"] = true
		}
	} else if pkg.Status == models.STATUS_PUBLISHED {
		this.Data["Tab"] = 2
		this.Data["Header"] = "Published"
	} else if pkg.Status == models.STATUS_REJECTED {
		this.Data["Tab"] = 3
		this.Data["Header"] = "Rejected"
	} else {
		this.Data["Tab"] = 4
		this.Data["Header"] = "Unknown"
	}
	this.Data["Package"] = pkg
	this.Data["Packages"] = strings.Split(pkg.Packages, ";")
	this.TplNames = "build.tpl"
}

func (this *BuildController) Post() {
	id := to.Uint64(this.Ctx.Input.Param(":buildid"))

	postType := this.GetString("type")
	if postType != "Up" && postType != "Down" {
		this.Abort("400")
	}

	user := util.IsLoggedIn(this.Controller)
	if user == "" {
		this.Abort("403") // MUST be logged in
	}

	var pkg models.BuildList

	o := orm.NewOrm()
	qt := o.QueryTable(new(models.BuildList))

	err := qt.Filter("ListId", id).One(&pkg)
	if err == orm.ErrNoRows {
		this.Abort("404")
	} else if err != nil {
		log.Println(err)
		this.Abort("500")
	}

	kt := o.QueryTable(new(models.Karma))

	var userkarma models.Karma
	err = kt.Filter("ListId", id).Filter("User", user).One(&userkarma)
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
		} else {
			if userkarma.Vote == models.KARMA_DOWN {
				o.Delete(&userkarma)
			} else {
				userkarma.Vote = models.KARMA_DOWN
				o.Update(&userkarma)
			}
		}
	} else {
		userkarma.ListId = id
		userkarma.User = user
		if postType == "Up" {
			userkarma.Vote = models.KARMA_UP
		} else {
			userkarma.Vote = models.KARMA_DOWN
		}
		o.Insert(&userkarma)
	}

	karmaup, _ := kt.Filter("ListId", id).Filter("Vote", models.KARMA_UP).Count()
	karmadown, _ := kt.Filter("ListId", id).Filter("Vote", models.KARMA_DOWN).Count()

	this.Data["Karma"] = karmaup - karmadown

	upthreshold, err := beego.AppConfig.Int64("upperkarma")
	if err != nil {
		// assume reasonable default
		upthreshold = 3
	}

	downthreshold, err := beego.AppConfig.Int64("lowerkarma")
	if err != nil {
		// assume reasonable default
		downthreshold = -3
	}

	if karmaup-karmadown >= upthreshold {
		pkg.Status = models.STATUS_PUBLISHED
		o.Update(&pkg)
		go integration.Publish(pkg.PublishHandle, pkg.ListId)
	} else if karmaup-karmadown <= downthreshold {
		pkg.Status = models.STATUS_REJECTED
		o.Update(&pkg)
		go integration.Reject(pkg.RejectHandle, pkg.ListId)
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
