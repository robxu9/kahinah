package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/robxu9/kahinah/common/conf"
	"github.com/robxu9/kahinah/common/klog"
	"github.com/robxu9/kahinah/common/set"
	"github.com/robxu9/kahinah/models"
	"menteslibres.net/gosexy/dig"
	"menteslibres.net/gosexy/to"
)

var (
	abfEnable        = conf.GetDefault("integration.abf.enable", false).(bool)
	abfURL           = conf.GetDefault("integration.abf.host", "https://abf.io").(string)
	abfUser          = conf.Get("integration.abf.user").(string)
	abfAPIKey        = conf.Get("integration.abf.apiKey").(string)
	abfPlatforms     = conf.Get("integration.abf.readPlatforms").([]interface{})
	abfArchWhitelist = conf.Get("integration.abf.archWhitelist").([]interface{})
	abfGitDiff       = conf.Get("integration.abf.gitDiff").(bool)
	abfGitDiffSSH    = conf.Get("integration.abf.gitDiffSSH").(bool)

	abfPlatformsSet     *set.Set
	abfArchWhitelistSet *set.Set
	abfClient           *http.Client
)

func init() {
	abfPlatformsSet = set.NewSet()
	for _, v := range abfPlatforms {
		abfPlatformsSet.Add(v.(string))
	}

	abfArchWhitelistSet = set.NewSet()
	for _, v := range abfArchWhitelist {
		abfArchWhitelistSet.Add(v.(string))
	}

	abfClient = &http.Client{
		Timeout: time.Second * 60,
	}

	if abfEnable {
		Implementations["abf"] = &ABF{}
	}
}

// ABF is the integration type for the Automated Build Farm.
// It uses the following fields:
// IntegrationName = "abf"
// IntegrationOne = Build List IDs ([id];[id];[id])
// IntegrationTwo = Commit Hash
// IntegrationThree
// IntegrationFour
// IntegrationFive
type ABF struct{}

// Poll polls ABF for new build lists with the status
// [testing] build completed and build completed.
func (a *ABF) Poll() error {
	for v := range abfPlatformsSet.Iterator() {
		a.pollBuildsInTesting(v) // poll testing builds first if we missed any
		a.pollBuildsCompleted(v) // then poll completed builds
	}

	return nil
}

// Hook is not implemented by ABF, because ABF does not support webhooks.
func (a *ABF) Hook(r *http.Request) {
	// Not implemented
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

	for _, v := range lists {
		asserted := v.(map[string]interface{})
		id := dig.Uint64(&asserted, "id")
		strID := "[" + to.String(id) + "]"

		var num int
		err = models.DB.Model(&models.List{}).Where("integration_name = ? AND integration_one LIKE ?", "abf", "%"+strID+"%").Count(&num).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			klog.Criticalf("abf: couldn't check db for existing (list %v): %v", id, err)
			continue
		}

		if num != 0 {
			klog.Infof("abf: ignoring list, already processed (list %v)", id)
			continue
		}

		json, err := a.getJSONList(id)
		if err != nil {
			klog.Criticalf("abf: couldn't retrieve the build list JSON (list %v): %v", id, err)
			continue
		}

		// check if arch is on whitelist
		if abfArchWhitelistSet != nil && !abfArchWhitelistSet.Contains(dig.String(&json, "arch", "name")) {
			klog.Infof("abf: ignoring list, arch not on whitelist (list %v)", id)
			continue
		}

		// check if platform is on whitelist
		if !abfPlatformsSet.Contains(dig.String(&json, "save_to_repository", "platform", "name")) {
			klog.Infof("abf: ignoring list, platform not on whitelist (list %v)", id)
			continue
		}

		// *check for duplicates before continuing
		// *we only check for duplicates in the same platform; different platforms have different conditions

		// check for duplicates
		var duplicate models.List
		err = models.DB.Where("platform = ? AND integration_two = ? AND stage_current <> ?",
			dig.String(&json, "save_to_repository", "platform", "name"),
			dig.String(&json, "commit_hash"),
			models.ListStageFinished).First(&duplicate).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			klog.Criticalf("abf: couldn't check db for duplicates (list %v): %v", id, err)
			continue
		}

		if err == nil { // we had no problem finding a duplicate, so handle and continue
			duplicate.IntegrationOne += ";" + strID
			duplicate.Variants += ";" + dig.String(&json, "arch", "name")

			if testing { // we got this from the build completed (not in testing) list
				// send id to testing
				go a.sendToTesting(id)
			}

			err = models.DB.Save(&duplicate).Error
			if err != nil {
				klog.Criticalf("abf: couldn't save duplicate modification to %v (list %v): %v", duplicate.ID, id, err)
			}

			pkgs := a.makePkgList(json)
			for _, listpkg := range pkgs {
				listpkg.ListID = duplicate.ID

				err = models.DB.Create(&listpkg).Error
				if err != nil {
					klog.Criticalf("abf: couldn't save new list package to %v (list %v): %v", duplicate.ID, id, err)
				}
			}

			// add a link to it
			newLink := models.ListLink{
				ListID: duplicate.ID,
				Name:   fmt.Sprintf("Build List for %v", dig.String(&json, "arch", "name")),
				URL:    fmt.Sprintf("%s/build_lists/%v", abfURL, id),
			}
			if err := models.DB.Create(&newLink).Error; err != nil {
				klog.Criticalf("abf: couldn't save new list link to %v (list %v): %v", duplicate.ID, id, err)
			}

			// ok, we're done here
			continue
		}

		list, err := a.makeBuildList(json)
		if err != nil {
			klog.Criticalf("abf: couldn't make new list (list %v): %v\n", id, err)
			continue
		}

		if testing {
			// Now send it to testing
			go a.sendToTesting(id)
		}

		if err := models.DB.Create(list).Error; err != nil {
			klog.Criticalf("abf: couldn't create new list in db (list %v): %v", id, err)
			continue
		}
	}

	return nil
}

func (a *ABF) pollBuildsCompleted(platformID string) error {
	// regular usage: use 0 (Build has been completed)
	// below: use 12000 ([testing] build has been published)
	req, err := http.NewRequest("GET", abfURL+"/api/v1/build_lists.json?per_page=100&filter[status]=0&filter[ownership]=index&filter[platform_id]="+platformID, nil)
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

func (a *ABF) pollBuildsInTesting(platformID string) error {
	req, err := http.NewRequest("GET", abfURL+"/api/v1/build_lists.json?per_page=100&filter[status]=12000&filter[ownership]=index&filter[platform_id]="+platformID, nil)
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

// Accept tells ABF to publish the build lists specified in the List.
func (a *ABF) Accept(m *models.List) error {
	for _, v := range strings.Split(m.IntegrationOne, ";") {
		noBrackets := v[1 : len(v)-1]
		id := to.Uint64(noBrackets)
		req, err := http.NewRequest("PUT", abfURL+"/api/v1/build_lists/"+to.String(id)+"/publish.json", nil)
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

// Reject tells ABF to reject the build lists specified in the List.
func (a *ABF) Reject(m *models.List) error {
	for _, v := range strings.Split(m.IntegrationOne, ";") {
		noBrackets := v[1 : len(v)-1]
		id := to.Uint64(noBrackets)
		req, err := http.NewRequest("PUT", abfURL+"/api/v1/build_lists/"+to.String(id)+"/reject_publish.json", nil)
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
	req, err := http.NewRequest("PUT", abfURL+"/api/v1/build_lists/"+to.String(id)+"/publish_into_testing.json", nil)
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
	req, err := http.NewRequest("GET", abfURL+"/api/v1/build_lists/"+to.String(id)+".json", nil)
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

func (a *ABF) makeBuildList(list map[string]interface{}) (*models.List, error) {
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

	handleID := to.String(dig.Uint64(&list, "id"))
	handleProject := dig.String(&list, "project", "fullname")
	platform := dig.String(&list, "save_to_repository", "platform", "name")

	bl := models.List{
		Name:     dig.String(&list, "project", "name"),
		Platform: platform,
		Channel:  dig.String(&list, "save_to_repository", "name"),
		Variants: dig.String(&list, "arch", "name"),

		Artifacts: pkg,
		Links: []models.ListLink{
			{
				Name: models.LinkMain,
				URL:  fmt.Sprintf("%s/build_lists/%v", abfURL, handleID),
			},
			{
				Name: models.LinkChangelog,
				URL:  changelog,
			},
			{
				Name: models.LinkSCM,
				URL:  fmt.Sprintf("%s/%v/commits/%v", abfURL, handleProject, platform),
			},
		},
		Changes:   changelog,
		BuildDate: time.Unix(dig.Int64(&list, "updated_at"), 0),

		StageCurrent: models.ListStageNotStarted,
		StageResult:  models.ListRunning,
		Activity: []models.ListActivity{
			{
				UserID:   models.FindUser(models.UserSystem).ID,
				Activity: "Imported this build list from ABF.",
			},
		},

		IntegrationName: "abf",
		IntegrationOne:  "[" + handleID + "]",
		IntegrationTwo:  dig.String(&list, "commit_hash"),
	}

	return &bl, nil
}

// makeDiff shells out to the cmd line git to get a patch. Unfortunately, this
// may need tweaking given the modifications that OpenMandriva have done to
// their ABF.
func (a *ABF) makeDiff(gitURL, fromHash, toHash string) string {
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

	if strings.Contains(gitURL, "@") {
		gitURL = gitURL[:strings.Index(gitURL, "//")+2] + gitURL[strings.Index(gitURL, "@")+1:]
	}

	urlToUse := gitURL
	if abfGitDiffSSH {
		urlToUse = strings.Replace("git@"+gitURL[strings.Index(gitURL, "//")+2:], "/", ":", 1)
	}

	// reusable bytes output
	var b bytes.Buffer

	gitclonecmd := exec.Command("git", "clone", urlToUse, tmpdir)
	gitclonecmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0") // git 2.3+
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
			klog.Criticalf("abf: calling git to kill failed: %v", err)
		}
		<-gitresult // allow goroutine to exit
		klog.Critical("abf: git process called to kill")
	case err := <-gitresult:
		if err != nil { // git exited with non-zero status
			klog.Criticalf("abf: calling git failed: %v", err)
			return "Repository could not be cloned for diff: " + err.Error()
		}
	}

	if fromHash == "" || fromHash == toHash {
		gitdiffcmd := exec.Command("git", "show", "--format=fuller", "--patch-with-stat", "--summary", toHash)
		gitdiffcmd.Dir = tmpdir

		gitdiff, err := gitdiffcmd.CombinedOutput()
		if err != nil {
			klog.Criticalf("abf: calling git diff failed: %v", err)
			return "No diff available: " + err.Error()
		}

		return fmt.Sprintf("$ git show --format=fuller --patch-with-stat --summary %s\n\n%s", toHash, string(gitdiff))
	}

	gitdiffcmd := exec.Command("git", "diff", "--patch-with-stat", "--summary", fromHash+".."+toHash)
	gitdiffcmd.Dir = tmpdir

	gitdiff, err := gitdiffcmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s", gitdiff)
		return "No diff available: " + err.Error()
	}

	return fmt.Sprintf("$ git diff --patch-with-stat --summary %s\n\n%s", fromHash+".."+toHash, string(gitdiff))
}

func (a *ABF) makePkgList(json map[string]interface{}) []models.ListArtifact {
	var pkgs []interface{}
	dig.Get(&json, &pkgs, "packages")

	var pkg []models.ListArtifact

	pkgrep := func(m map[string]interface{}) models.ListArtifact {
		pkgType := dig.String(&m, "type")
		if strings.HasSuffix(dig.String(&m, "name"), "-debuginfo") {
			pkgType = "debuginfo"
		}

		return models.ListArtifact{
			Name:    dig.String(&m, "name"),
			Type:    pkgType,
			Epoch:   dig.Int64(&m, "epoch"),
			Version: dig.String(&m, "version"),
			Release: dig.String(&m, "release"),
			Variant: dig.String(&json, "arch", "name"),
			URL:     dig.String(&m, "url"),
		}
	}

	for _, v := range pkgs {
		asserted := v.(map[string]interface{})
		pkg = append(pkg, pkgrep(asserted))
	}

	return pkg
}
