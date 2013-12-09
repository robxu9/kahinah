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
	"strings"
	"time"
)

const (
	ABF_URL = "https://abf.io/api/v1"
)

var (
	user      = beego.AppConfig.String("abf_user")
	pass      = beego.AppConfig.String("abf_pass")
	platforms util.Set
	client    = &http.Client{}
)

func init() {
	platforms = util.NewSet()
	for _, v := range strings.Split(beego.AppConfig.String("abf_platforms"), ";") {
		platforms.Add(v)
	}
}

type ABF byte

func (a ABF) Ping() error {
	err := a.pingBuildCompleted()
	err2 := a.pingTestingBuilds()

	// FUTURE TODO: ping published & rejected builds to ensure consistency

	if err != nil && err2 != nil {
		return fmt.Errorf("abf: two errors: %s | %s", err, err2)
	} else if err != nil {
		return err
	} else if err2 != nil {
		return err2
	}

	return nil
}

func (a ABF) pingBuildCompleted() error {
	// regular usage: use 0 (Build has been completed)
	// below: use 12000 ([testing] build has been published)
	// FIXME - remove hardcoded openmandriva2013.0 filter (but find alternative, we need said filter)
	req, err := http.NewRequest("GET", ABF_URL+"/build_lists.json?per_page=100&filter[status]=0&filter[ownership]=index&filter[platform_id]=668", nil)
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

		var list models.BuildList
		err := o.QueryTable(new(models.BuildList)).Filter("ListId", id).One(&list)
		if err != nil {
			list, err = a.getBuildList(id)
			if err != nil {
				log.Printf("abf: Error retrieving build list %s: %s\n", id, err)
			}

			if !platforms.Contains(list.Platform) {
				// ignore
				continue
			}

			// Now send it to testing
			go a.sendToTesting(id)

			_, err = o.Insert(&list)
			if err != nil {
				log.Printf("abf: Error saving build list %s: %s\n", id, err)
			}
		}
	}

	return nil
}

func (a ABF) pingTestingBuilds() error {
	// FIXME - remove hardcoded openmandriva2013.0 filter (but find alternative, we need said filter)
	req, err := http.NewRequest("GET", ABF_URL+"/build_lists.json?per_page=100&filter[status]=12000&filter[ownership]=index&filter[platform_id]=668", nil)
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

		var list models.BuildList
		err := o.QueryTable(new(models.BuildList)).Filter("ListId", id).One(&list)
		if err != nil {
			list, err = a.getBuildList(id)
			if err != nil {
				log.Printf("abf: Error retrieving build list %s: %s\n", id, err)
			}

			if !platforms.Contains(list.Platform) {
				// ignore
				continue
			}

			_, err = o.Insert(&list)
			if err != nil {
				log.Printf("abf: Error saving build list %s: %s\n", id, err)
			}
		}
	}

	return nil
}

func (a ABF) PingParams(m map[string]string) error {
	return fmt.Errorf("abf: PingParams not supported yet.")
}

func (a ABF) Publish(id uint64) error {
	req, err := http.NewRequest("PUT", ABF_URL+"/build_lists/"+to.String(id)+"/publish.json", nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (a ABF) Reject(id uint64) error {
	req, err := http.NewRequest("PUT", ABF_URL+"/build_lists/"+to.String(id)+"/reject_publish.json", nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(user, pass)

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

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (a ABF) getBuildList(id uint64) (models.BuildList, error) {
	req, err := http.NewRequest("GET", ABF_URL+"/build_lists/"+to.String(id)+".json", nil)
	if err != nil {
		return models.BuildList{}, err
	}

	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		return models.BuildList{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.BuildList{}, err
	}

	var result map[string]interface{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return models.BuildList{}, err
	}

	var list map[string]interface{}

	dig.Get(&result, &list, "build_list")

	pkgs := make([]interface{}, 0)
	dig.Get(&list, &pkgs, "packages")

	pkg := ""

	pkgrep := func(m map[string]interface{}) string {
		res := dig.String(&m, "name") + " "
		if dig.String(&m, "epoch") != "" {
			res += dig.String(&m, "epoch")
			res += ":"
		}
		res += dig.String(&m, "version") + "-"
		res += dig.String(&m, "release")
		return res
	}

	for _, v := range pkgs {
		asserted := v.(map[string]interface{})
		if dig.String(&asserted, "type") == "source" {
			continue // no source packages
		}
		if strings.HasSuffix(dig.String(&asserted, "name"), "-debuginfo") {
			continue // no debuginfo packages
		}
		if pkg == "" {
			pkg = pkgrep(asserted)
		} else {
			pkg = pkg + ";" + pkgrep(asserted)
		}
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

	bl := models.BuildList{ListId: dig.Uint64(&list, "id"),
		//Platform:     dig.String(&list, "build_for_platform", "name"),
		Platform:      dig.String(&list, "save_to_repository", "platform", "name"),
		Repo:          dig.String(&list, "save_to_repository", "name"),
		Architecture:  dig.String(&list, "arch", "name"),
		Name:          dig.String(&list, "project", "name"),
		Submitter:     dig.String(&list, "user", "name"),
		Type:          dig.String(&list, "update_type"),
		Status:        models.STATUS_TESTING,
		Url:           "https://abf.io/build_lists/" + to.String(dig.Uint64(&list, "id")),
		Changelog:     changelog, // url
		PublishHandle: "abf",
		RejectHandle:  "abf",
		Packages:      pkg,
		BuildDate:     time.Unix(dig.Int64(&list, "updated_at"), 0)}

	return bl, nil

}
