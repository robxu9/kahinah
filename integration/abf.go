package integration

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/robxu9/kahinah/models"
	"github.com/robxu9/kahinah/util"
	"io/ioutil"
	"log"
	"menteslibres.net/gosexy/dig"
	"menteslibres.net/gosexy/to"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	ABF_URL = "https://abf.io/api/v1"
)

var (
	user        = beego.AppConfig.String("abf::abf_user")
	pass        = beego.AppConfig.String("abf::abf_pass")
	platforms   *util.Set
	platformids *util.Set
	client      = &http.Client{}
)

func init() {
	platforms = util.NewSet()
	platformids = util.NewSet()
	for _, v := range strings.Split(beego.AppConfig.String("abf::abf_platforms"), ";") {
		platform := strings.Split(v, ":")
		platforms.Add(platform[0])
		platformids.Add(platform[1])
	}
}

type ABF byte

func (a ABF) Ping() error {

	for v := range platformids.Iterator() {
		go a.pingBuildCompleted(v)
		go a.pingTestingBuilds(v)
	}

	// FUTURE TODO: ping published & rejected builds to ensure consistency

	//if err != nil && err2 != nil {
	//	return fmt.Errorf("abf: two errors: %s | %s", err, err2)
	//} else if err != nil {
	//	return err
	//} else if err2 != nil {
	//	return err2
	//}

	return nil
}

func (a ABF) pingBuildCompleted(platformId string) error {
	// regular usage: use 0 (Build has been completed)
	// below: use 12000 ([testing] build has been published)
	// FIXME - remove hardcoded openmandriva2013.0 filter (but find alternative, we need said filter)
	req, err := http.NewRequest("GET", ABF_URL+"/build_lists.json?per_page=100&filter[status]=0&filter[ownership]=index&filter[platform_id]="+platformId, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]interface{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	var lists []interface{}

	err = dig.Get(&result, &lists, "build_lists")
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	for _, v := range lists {
		asserted := v.(map[string]interface{})
		id := dig.Uint64(&asserted, "id")

		num, err := o.QueryTable(new(models.BuildList)).Filter("HandleId", id).Count()
		if num <= 0 || err != nil {
			list, err := a.getBuildList(id)
			if err != nil {
				log.Printf("abf: Error retrieving build list %s: %s\n", id, err)
			}

			if !platforms.Contains(list.Platform) {
				// ignore
				continue
			}

			// Now send it to testing
			go a.sendToTesting(id)

			_, err = o.Insert(list)
			if err != nil {
				log.Printf("abf: Error saving build list %s: %s\n", id, err)
			}

			for _, listpkg := range list.Packages {
				listpkg.List = list
				o.Insert(listpkg)
			}
		}
	}

	return nil
}

func (a ABF) pingTestingBuilds(platformId string) error {
	// FIXME - remove hardcoded openmandriva2013.0 filter (but find alternative, we need said filter)
	req, err := http.NewRequest("GET", ABF_URL+"/build_lists.json?per_page=100&filter[status]=12000&filter[ownership]=index&filter[platform_id]="+platformId, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]interface{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	var lists []interface{}

	err = dig.Get(&result, &lists, "build_lists")
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	for _, v := range lists {
		asserted := v.(map[string]interface{})
		id := dig.Uint64(&asserted, "id")

		num, err := o.QueryTable(new(models.BuildList)).Filter("HandleId", id).Count()
		if num <= 0 || err != nil {
			list, err := a.getBuildList(id)
			if err != nil {
				log.Printf("abf: Error retrieving build list %s: %s\n", id, err)
			}

			if !platforms.Contains(list.Platform) {
				// ignore
				continue
			}

			_, err = o.Insert(list)
			if err != nil {
				log.Printf("abf: Error saving build list %s: %s\n", id, err)
			}

			for _, listpkg := range list.Packages {
				listpkg.List = list
				o.Insert(listpkg)
			}
		}
	}

	return nil
}

func (a ABF) PingParams(m map[string]string) error {
	return fmt.Errorf("abf: PingParams not supported yet.")
}

func (a ABF) Url(m *models.BuildList) string {
	return "https://abf.io/build_lists/" + m.HandleId
}

func (a ABF) Publish(m *models.BuildList) error {
	id := to.Uint64(m.HandleId)
	req, err := http.NewRequest("PUT", ABF_URL+"/build_lists/"+to.String(id)+"/publish.json", nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(user, pass)
	req.Header.Add("Content-Length", "0")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (a ABF) Reject(m *models.BuildList) error {
	id := to.Uint64(m.HandleId)
	req, err := http.NewRequest("PUT", ABF_URL+"/build_lists/"+to.String(id)+"/reject_publish.json", nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(user, pass)
	req.Header.Add("Content-Length", "0")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (a ABF) sendToTesting(id uint64) error {
	req, err := http.NewRequest("PUT", ABF_URL+"/build_lists/"+to.String(id)+"/publish_into_testing.json", nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(user, pass)
	req.Header.Add("Content-Length", "0")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (a ABF) getBuildList(id uint64) (*models.BuildList, error) {
	req, err := http.NewRequest("GET", ABF_URL+"/build_lists/"+to.String(id)+".json", nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	var list map[string]interface{}

	dig.Get(&result, &list, "build_list")

	pkgs := make([]interface{}, 0)
	dig.Get(&list, &pkgs, "packages")

	pkg := make([]*models.BuildListPkg, 0)

	pkgrep := func(m map[string]interface{}) *models.BuildListPkg {
		pkgType := dig.String(&m, "type")
		if strings.HasSuffix(dig.String(&m, "name"), "-debuginfo") {
			pkgType = "debuginfo"
		}

		return &models.BuildListPkg{
			Name:    dig.String(&m, "name"),
			Type:    pkgType,
			Epoch:   dig.Int64(&m, "epoch"),
			Version: dig.String(&m, "version"),
			Release: dig.String(&m, "release"),
			Url:     dig.String(&m, "url"),
		}
	}

	for _, v := range pkgs {
		asserted := v.(map[string]interface{})
		pkg = append(pkg, pkgrep(asserted))
	}

	changelog := ""

	logs := make([]interface{}, 0)
	dig.Get(&list, &logs, "logs")

	for _, v := range logs {
		asserted := v.(map[string]interface{})
		if dig.String(&asserted, "file_name") == "changelog.log" {
			changelog = dig.String(&asserted, "url")
			break
		}
	}

	user := a.getUser(dig.Uint64(&list, "user", "id"))

	bl := models.BuildList{
		HandleId:      to.String(dig.Uint64(&list, "id")),
		HandleProject: dig.String(&list, "project", "fullname"),
		Diff:          a.getDiff(dig.String(&list, "project", "git_url"), dig.String(&list, "last_published_commit_hash"), dig.String("commit_hash")),

		//Platform:     dig.String(&list, "build_for_platform", "name"),
		Platform:     dig.String(&list, "save_to_repository", "platform", "name"),
		Repo:         dig.String(&list, "save_to_repository", "name"),
		Architecture: dig.String(&list, "arch", "name"),
		Name:         dig.String(&list, "project", "name"),
		Submitter:    user,
		Type:         dig.String(&list, "update_type"),
		Status:       models.STATUS_TESTING,
		Changelog:    changelog, // url
		Packages:     pkg,
		BuildDate:    time.Unix(dig.Int64(&list, "updated_at"), 0),
	}

	return &bl, nil

}

func (a ABF) getDiff(gitUrl, fromHash, toHash string) string {
	// ugh, looks like we'll have to do this the sadly hard way
	tmpdir, err := ioutil.TempDir("", "kahinah_")
	if err != nil {
		return "Error creating directory for diff creation: " + err.Error()
	}

	defer os.RemoveAll(tmpdir)

	gitresult := exec.Command("git", "clone", gitUrl, tmpdir).Run()
	if gitresult != nil { // git better exit with status zero
		return "Repository could not be cloned for diff: " + gitresult.Error()
	}

	gitdiffcmd := exec.Command("git", "diff", "--patch-with-stat", "--summary", fromHash+".."+toHash)
	gitdiffcmd.Dir = tmpdir

	gitdiff, err := gitdiffcmd.Output()
	if err != nil {
		return "No diff available: " + err.Error()
	}

	return fmt.Sprintf("$ git diff --patch-with-stat --summary %s..%s\n\n%s", fromHash, toHash, string(gitdiff))
}

func (a ABF) getUser(id uint64) *models.User {

	o := orm.NewOrm()

	var userModel *models.User

	var userIntegration models.User
	err := o.QueryTable(new(models.User)).Filter("Integration", id).One(&userIntegration)
	if err == orm.ErrNoRows {

		req, err := http.NewRequest("GET", ABF_URL+"/users/"+to.String(id)+".json", nil)
		if err != nil {
			return nil
		}

		req.SetBasicAuth(user, pass)

		resp, err := client.Do(req)
		if err != nil {
			return nil
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil
		}

		var result map[string]interface{}

		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil
		}

		email := dig.String(&result, "user", "email")

		userModel = models.FindUser(email)
		userModel.Integration = to.String(id)
		userModel.Save()

	} else if err != nil {
		panic(err)
	} else {
		userModel = &userIntegration
	}

	return userModel
}
