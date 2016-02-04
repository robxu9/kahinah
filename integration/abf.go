package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/robxu9/kahinah/conf"
	"github.com/robxu9/kahinah/models"
	"github.com/robxu9/kahinah/util"
	"menteslibres.net/gosexy/dig"
	"menteslibres.net/gosexy/to"
)

var (
	abfEnable        = conf.Config.GetDefault("integration.abf.enable", false).(bool)
	abfURL           = conf.Config.GetDefault("integration.abf.host", "https://abf.io/api/v1").(string)
	abfUser          = conf.Config.Get("integration.abf.user").(string)
	abfAPIKey        = conf.Config.Get("integration.abf.apiKey").(string)
	abfPlatforms     = conf.Config.Get("integration.abf.readPlatforms").([]interface{})
	abfArchWhitelist = conf.Config.Get("integration.abf.archWhitelist").([]interface{})
	abfGitDiff       = conf.Config.Get("integration.abf.gitDiff").(bool)
	abfGitDiffSSH    = conf.Config.Get("integration.abf.gitDiffSSH").(bool)

	abfPlatformsSet     *util.Set
	abfArchWhitelistSet *util.Set
	abfClient           *http.Client
)

func init() {
	abfPlatformsSet = util.NewSet()
	for _, v := range abfPlatforms {
		abfPlatformsSet.Add(v.(string))
	}

	abfArchWhitelistSet = util.NewSet()
	for _, v := range abfArchWhitelist {
		abfArchWhitelistSet.Add(v.(string))
	}

	abfClient = &http.Client{}
}

type ABF struct{}

func (a *ABF) Poll() error {
	if !abfEnable {
		return ErrDisabled
	}

	for v := range abfPlatformsSet.Iterator() {
		a.pollBuildsInTesting(v) // poll testing builds first if we missed any
		a.pollBuildsCompleted(v) // then poll completed builds
	}

	return nil
}

func (a *ABF) Hook(i interface{}) error {
	if !abfEnable {
		return ErrDisabled
	}

	return ErrNotImplemented
}

func (a *ABF) handleResponse(resp *http.Response, testing bool) error {

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

		num, err := o.QueryTable(new(models.BuildList)).Filter("HandleId__icontains", id).Count()
		if num <= 0 || err != nil {
			json, err := a.getJSONList(id)
			if err != nil {
				log.Printf("abf: error retrieving build list json %v: %v\n", id, err)
				continue
			}

			// check if arch is on whitelist
			if abfArchWhitelistSet != nil && !abfArchWhitelistSet.Contains(dig.String(&json, "arch", "name")) {
				// we are ignoring this buildlist
				continue
			}

			// check if platform is on whitelist
			if !abfPlatformsSet.Contains(dig.String(&json, "save_to_repository", "platform", "name")) {
				// we are ignoring this platform
				continue
			}

			// *check for duplicates before continuing
			// *we only check for duuplicates in the same platform; different platforms have different conditions
			var possibleDuplicate models.BuildList
			err = o.QueryTable(new(models.BuildList)).Filter("Platform", dig.String(&json, "save_to_repository", "platform", "name")).Filter("HandleCommitId", dig.String(&json, "commit_hash")).Filter("Status", models.StatusTesting).One(&possibleDuplicate)
			if err == nil { // we found a duplicate... handle and continue
				possibleDuplicate.HandleId = possibleDuplicate.HandleId + ";" + to.String(id)
				possibleDuplicate.Architecture += ";" + dig.String(&json, "arch", "name")

				if testing {
					// send id to testing
					go a.sendToTesting(id)
				}

				o.Update(&possibleDuplicate)

				pkgs := a.makePkgList(json)
				for _, listpkg := range pkgs {
					listpkg.List = &possibleDuplicate
					o.Insert(listpkg)
				}

				// ok, we're done here
				continue
			}

			list, err := a.makeBuildList(json)
			if err != nil {
				log.Printf("abf: Error retrieving build list %v: %v\n", id, err)
				continue
			}

			if testing {
				// Now send it to testing
				go a.sendToTesting(id)
			}

			_, err = o.Insert(list)
			if err != nil {
				log.Printf("abf: Error saving build list %v: %v\n", id, err)
				continue
			}

			for _, listpkg := range list.Packages {
				listpkg.List = list
				o.Insert(listpkg)
			}

			for _, listlinks := range list.Links {
				listlinks.List = list
				o.Insert(listlinks)
			}

			for _, listkarma := range list.Karma {
				listkarma.List = list
				o.Insert(listkarma)
			}

			go util.MailModel(list)
		}
	}

	return nil
}

func (a *ABF) pollBuildsCompleted(platformId string) error {
	// regular usage: use 0 (Build has been completed)
	// below: use 12000 ([testing] build has been published)
	req, err := http.NewRequest("GET", abfURL+"/build_lists.json?per_page=100&filter[status]=0&filter[ownership]=index&filter[platform_id]="+platformId, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(abfUser, abfAPIKey)

	resp, err := abfClient.Do(req)
	if err != nil {
		return err
	}

	return a.handleResponse(resp, true)
}

func (a *ABF) pollBuildsInTesting(platformId string) error {
	req, err := http.NewRequest("GET", abfURL+"/build_lists.json?per_page=100&filter[status]=12000&filter[ownership]=index&filter[platform_id]="+platformId, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(abfUser, abfAPIKey)

	resp, err := abfClient.Do(req)
	if err != nil {
		return err
	}

	return a.handleResponse(resp, false)
}

func (a *ABF) Accept(m *models.BuildList) error {
	go util.MailModel(m)

	for _, v := range strings.Split(m.HandleId, ";") {

		id := to.Uint64(v)
		req, err := http.NewRequest("PUT", abfURL+"/build_lists/"+to.String(id)+"/publish.json", nil)
		if err != nil {
			return err
		}

		req.SetBasicAuth(abfUser, abfAPIKey)
		req.Header.Add("Content-Length", "0")

		resp, err := abfClient.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
	}
	return nil
}

func (a *ABF) Reject(m *models.BuildList) error {
	go util.MailModel(m)

	for _, v := range strings.Split(m.HandleId, ";") {
		id := to.Uint64(v)
		req, err := http.NewRequest("PUT", abfURL+"/build_lists/"+to.String(id)+"/reject_publish.json", nil)
		if err != nil {
			return err
		}

		req.SetBasicAuth(abfUser, abfAPIKey)
		req.Header.Add("Content-Length", "0")

		resp, err := abfClient.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
	}

	return nil
}

func (a *ABF) sendToTesting(id uint64) error {
	req, err := http.NewRequest("PUT", abfURL+"/build_lists/"+to.String(id)+"/publish_into_testing.json", nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(abfUser, abfAPIKey)
	req.Header.Add("Content-Length", "0")

	resp, err := abfClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to send %d to testing: %s\n", id, err.Error())
		return err
	}

	defer resp.Body.Close()
	// we are not doing this again...
	//bte, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf("sending %d to testing yielded %s\n", id, bte)
	return nil
}

func (a *ABF) getJSONList(id uint64) (list map[string]interface{}, err error) {
	req, err := http.NewRequest("GET", abfURL+"/build_lists/"+to.String(id)+".json", nil)
	if err != nil {
		return
	}

	req.SetBasicAuth(abfUser, abfAPIKey)

	resp, err := abfClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var result map[string]interface{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}

	dig.Get(&result, &list, "build_list")
	return
}

func (a *ABF) makeBuildList(list map[string]interface{}) (*models.BuildList, error) {
	pkg := a.makePkgList(list)

	changelog := ""

	var logs []interface{}
	dig.Get(&list, &logs, "logs")

	for _, v := range logs {
		asserted := v.(map[string]interface{})
		if dig.String(&asserted, "file_name") == "changelog.log" {
			changelog = dig.String(&asserted, "url")
			break
		}
	}

	user := a.getUser(dig.Uint64(&list, "user", "id"))
	handleId := to.String(dig.Uint64(&list, "id"))
	handleProject := dig.String(&list, "project", "fullname")
	platform := dig.String(&list, "save_to_repository", "platform", "name")

	bl := models.BuildList{
		HandleId:       handleId,
		HandleProject:  handleProject,
		HandleCommitId: dig.String(&list, "commit_hash"),
		Diff:           a.makeDiff(dig.String(&list, "project", "git_url"), dig.String(&list, "last_published_commit_hash"), dig.String(&list, "commit_hash")),

		//Platform:     dig.String(&list, "build_for_platform", "name"),
		Platform:     platform,
		Repo:         dig.String(&list, "save_to_repository", "name"),
		Architecture: dig.String(&list, "arch", "name"),
		Name:         dig.String(&list, "project", "name"),
		Submitter:    user,
		Type:         dig.String(&list, "update_type"),
		Status:       models.StatusTesting,
		Links: []*models.BuildListLink{
			&models.BuildListLink{
				Name: models.LinkMainURL,
				Url:  fmt.Sprintf("https://abf.io/build_lists/%v", handleId),
			},
			&models.BuildListLink{
				Name: models.ChangelogURL,
				Url:  changelog,
			},
			&models.BuildListLink{
				Name: models.SCMLogURL,
				Url:  fmt.Sprintf("https://abf.io/%v/commits/%v", handleProject, platform),
			},
		},
		Packages:  pkg,
		BuildDate: time.Unix(dig.Int64(&list, "updated_at"), 0),
		Karma: []*models.Karma{
			&models.Karma{
				User:    models.FindUser(models.UserSystem),
				Vote:    models.KARMA_NONE,
				Comment: "Imported this build from ABF.",
			},
		},
	}

	return &bl, nil
}

func (a *ABF) makeDiff(gitUrl, fromHash, toHash string) string {
	// make sure it's not disabled
	if !abfGitDiff {
		return "Diff creation disabled."
	}

	// ugh, looks like we'll have to do this the sadly hard way
	tmpdir, err := ioutil.TempDir("", "kahinah_")
	if err != nil {
		return "Error creating directory for diff creation: " + err.Error()
	}

	defer os.RemoveAll(tmpdir)

	if strings.Contains(gitUrl, "@") {
		gitUrl = gitUrl[:strings.Index(gitUrl, "//")+2] + gitUrl[strings.Index(gitUrl, "@")+1:]
	}

	urlToUse := gitUrl
	if abfGitDiffSSH {
		urlToUse = strings.Replace("git@"+gitUrl[strings.Index(gitUrl, "//")+2:], "/", ":", 1)
	}

	// reusable bytes output
	var b bytes.Buffer

	gitclonecmd := exec.Command("git", "clone", urlToUse, tmpdir)
	gitclonecmd.Stderr = &b
	gitclonecmd.Stdout = &b
	gitclonecmd.Start()

	gitresult := make(chan error, 1)
	go func() {
		gitresult <- gitclonecmd.Wait()
	}()
	select {
	case <-time.After(10 * time.Minute):
		if err := gitclonecmd.Process.Kill(); err != nil {
			fmt.Fprintf(os.Stderr, "couldn't kill git process: %s\n", err.Error())
		}
		<-gitresult // allow goroutine to exit
		log.Println("process killed")
	case err := <-gitresult:
		if err != nil { // git exited with non-zero status
			fmt.Fprintf(os.Stderr, "git errored: %s\n", b.String())
			return "Repository could not be cloned for diff: " + err.Error()
		}
	}

	if fromHash == "" || fromHash == toHash {
		gitdiffcmd := exec.Command("git", "show", "--format=fuller", "--patch-with-stat", "--summary", toHash)
		gitdiffcmd.Dir = tmpdir

		gitdiff, err := gitdiffcmd.CombinedOutput()
		if err != nil {
			fmt.Printf("%s", gitdiff)
			return "No diff available: " + err.Error()
		}

		return fmt.Sprintf("$ git show --format=fuller --patch-with-stat --summary %s\n\n%s", toHash, string(gitdiff))
	} else {
		gitdiffcmd := exec.Command("git", "diff", "--patch-with-stat", "--summary", fromHash+".."+toHash)
		gitdiffcmd.Dir = tmpdir

		gitdiff, err := gitdiffcmd.CombinedOutput()
		if err != nil {
			fmt.Printf("%s", gitdiff)
			return "No diff available: " + err.Error()
		}

		return fmt.Sprintf("$ git diff --patch-with-stat --summary %s\n\n%s", fromHash+".."+toHash, string(gitdiff))
	}
}

func (a *ABF) makePkgList(json map[string]interface{}) []*models.BuildListPkg {
	pkgs := make([]interface{}, 0)
	dig.Get(&json, &pkgs, "packages")

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
			Arch:    dig.String(&json, "arch", "name"),
			Url:     dig.String(&m, "url"),
		}
	}

	for _, v := range pkgs {
		asserted := v.(map[string]interface{})
		pkg = append(pkg, pkgrep(asserted))
	}

	return pkg
}

func (a *ABF) getUser(id uint64) *models.User {

	o := orm.NewOrm()

	var userModel *models.User

	var userIntegration models.User
	err := o.QueryTable(new(models.User)).Filter("Integration", id).One(&userIntegration)
	if err == orm.ErrNoRows {

		req, err := http.NewRequest("GET", abfURL+"/users/"+to.String(id)+".json", nil)
		if err != nil {
			return nil
		}

		req.SetBasicAuth(abfUser, abfAPIKey)

		resp, err := abfClient.Do(req)
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

		name := dig.String(&result, "user", "name")

		userModel = models.FindUser(name)
		userModel.Integration = to.String(id)
		userModel.Save()

	} else if err != nil {
		panic(err)
	} else {
		userModel = &userIntegration
	}

	return userModel
}
