package abf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/pmylund/go-cache"
	"github.com/robxu9/kahinah/backend"
)

func (c *Connector) sendTesting(id int) {
	resp, err := c.makeRequest("PUT", fmt.Sprintf(EndpointBuildListSendTesting, id), nil)
	if err != nil {
		// log and die
		log.Printf("failed to call send testing endpoint for id %v: %v", id, err)
	}
	defer resp.Body.Close()
}

func (c *Connector) getUserEmail(id int) string {
	if cached, found := c.cache.Get("user/" + strconv.Itoa(id)); found {
		return cached.(string)
	}

	resp, err := c.makeRequest("GET", fmt.Sprintf(EndpointUser, id), nil)
	if err != nil {
		// log and die
		log.Printf("failed to call user endpoint: %v", err)
		return ""
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var user User
	if err := decoder.Decode(&user); err != nil {
		log.Printf("failed to decode to user: %v", err)
		return ""
	}

	c.cache.Set("user/"+strconv.Itoa(id), user.User.Email, cache.DefaultExpiration)

	return user.User.Email
}

func (c *Connector) getChangelog(w *WebBuildList) []byte {
	for _, v := range w.BuildList.Results {
		if v.FileName == ResultChangelog {
			b := &bytes.Buffer{}
			result, err := c.makeRequest("GET", v.URL, nil)
			if err != nil {
				log.Printf("unable to make request to %v for changelog build list %v", v.URL, w.BuildList.ID)
				return []byte("couldn't get changelog :(")
			}
			defer result.Body.Close()

			// I really don't care about the error, even though I should
			io.Copy(b, result.Body)

			return b.Bytes()
		}
	}

	return []byte{}
}

func (c *Connector) handleResponse(resp *http.Response) {
	decoder := json.NewDecoder(resp.Body)

	var buildlists BuildLists
	if err := decoder.Decode(&buildlists); err != nil {
		log.Printf("failed to decode list of builds: %v", err)
		return
	}

	// grab a map of arches
	arches := map[int]string{}
	for _, v := range buildlists.Dictionary.Arches {
		arches[v.ID] = v.Name
	}

	for _, v := range buildlists.BuildLists {
		c.handleID(arches[v.ArchID], v.ID)
	}
}

func (c *Connector) handleID(arch string, id int) {
	// if we already have this, then there's no need to go over it again.
	dbQuery, err := c.k.Run(c.k.DB().Table(kahinah.TableUpdate).Filter(map[string]interface{}{
		"connector": c.Name(),
		"connector_store": map[string]string{
			arch: strconv.Itoa(id),
		},
	}))

	if err != nil {
		log.Printf("failed to query kahinah for existing buildlists: %v", err)
		return
	}

	defer dbQuery.Close()

	if !dbQuery.IsNil() {
		log.Printf("already found build list %v processed, skipping", id)
		return
	}

	// get the two buildlists
	apiResp, err := c.makeRequest("GET", fmt.Sprintf(EndpointAPIBuildList, id), nil)
	if err != nil {
		log.Printf("failed to call api get specific build list %v: %v", id, err)
		return
	}
	defer apiResp.Body.Close()

	webResp, err := c.makeRequest("GET", fmt.Sprintf(EndpointWebBuildList, id), nil)
	if err != nil {
		log.Printf("failed to call get specific build list %v: %v", id, err)
		return
	}
	defer webResp.Body.Close()

	// we store copies of the responses in the DB
	apiCopy := &bytes.Buffer{}
	webCopy := &bytes.Buffer{}

	apiDecoder := json.NewDecoder(io.TeeReader(apiResp.Body, apiCopy))
	webDecoder := json.NewDecoder(io.TeeReader(webResp.Body, webCopy))

	var apiList APIBuildList
	var webList WebBuildList

	if err = apiDecoder.Decode(&apiList); err != nil {
		log.Printf("failed to decode api specific build list %v: %v", id, err)
		return
	}

	if err = webDecoder.Decode(&webList); err != nil {
		log.Printf("failed to decode web specific build list %v: %v", id, err)
		return
	}

	// check requirements (do we care about this arch/platform?)
	if !c.containsPlatform(apiList.BuildList.BuildForPlatform.ID) {
		log.Printf("platform not accepted for build list %v", id)
		return
	}

	if !c.containsArch(apiList.BuildList.Arch.Name) {
		log.Printf("arch not accepted for build list %v", id)
		return
	}

	// since this build list isn't in the db, check if it's waiting
	key := waitingKey{
		Fullname: apiList.BuildList.Project.Fullname,
		Hash:     apiList.BuildList.CommitHash,
	}
	waitingUpdate := c.waiting[key]

	if waitingUpdate != nil {
		// lovely, it is waiting. check if we have the arch, and add if
		// necessary, and then send it off if necessary.
		if _, ok := waitingUpdate.ConnectorStore[arch]; !ok {
			// we don't have this particular arch stored, let's store it now.
			waitingUpdate.ConnectorStore[arch] = strconv.Itoa(id)
			waitingUpdate.ConnectorStore[arch+"_API"] = apiCopy.String()
			waitingUpdate.ConnectorStore[arch+"_WEB"] = webCopy.String()

			// and let's also make sure to store this in the packages list
			waitingUpdate.Packages[arch] = webList.Packages()
		}
	} else {
		// there is no existing update - so let's create it
		waitingUpdate = &kahinah.Update{
			Platform:  apiList.BuildList.BuildForPlatform.Name,
			Name:      apiList.Name(),
			EVR:       apiList.EVR(),
			Submitter: c.getUserEmail(apiList.BuildList.User.ID),
			Type:      apiList.BuildList.UpdateType,
			Connector: c.Name(),
			Diff:      []byte(fmt.Sprintf("please see https://abf.io/%v/diff/%v...%v", apiList.BuildList.Project.Fullname, apiList.BuildList.LastPublishedCommitHash, apiList.BuildList.CommitHash)),
			Changelog: c.getChangelog(&webList),
			ConnectorStore: map[string]string{
				arch:          strconv.Itoa(id),
				arch + "_API": apiCopy.String(),
				arch + "_WEB": webCopy.String(),
			},
			Packages: map[string][]string{
				arch: webList.Packages(),
			},
		}
		c.waiting[key] = waitingUpdate
	}

	// check to see if we've added all arches necessary
	// (if we didn't add to the update, this is a consistency check)
	for _, v := range c.Arches {
		if _, ok := waitingUpdate.ConnectorStore[v]; !ok {
			// no, we're missing this architecture. we're done for now.
			// return.
			return
		}
	}

	// we've added all architectures, process this now.
	c.handlePublish(key)
}

func (c *Connector) handlePublish(w waitingKey) {
	update := c.waiting[w]
	if update == nil {
		panic("c.waiting[w] should not be nil")
	}

	update.Date = time.Now()

	// insert the update, then connect this to an advisory
	// TODO: insert here

	// parse rules and look for corresponding advisory.
	// if no such advisory exists, create one.
	rule := c.k.Rules[0]
	for _, v := range c.k.Rules {
		if v.Matches(update) {
			rule = v
		}
	}

	// we have a rule; now find an advisory that matches that rule name.
	// if we already have this, then there's no need to go over it again.
	dbQuery, err := c.k.Run(c.k.DB().Table(kahinah.TableAdvisory).Filter(map[string]interface{}{
		"pushed":     false,
		"deprecated": false,
		"ruleset":    rule.RuleName,
	}))

	if err != nil {
		log.Printf("failed to query kahinah for un-pushed un-deprecated advisory with ruleset %v: %v", rule.RuleName, err)
		return
	}

	defer dbQuery.Close()
}

func (c *Connector) poll() {
	for _, platform := range c.Platforms {
		for _, status := range []string{StatusBuildCompleted, StatusBuildTestingPublished} {
			resp, err := c.makeRequest("GET", fmt.Sprintf(EndpointBuildLists, status, platform), nil)
			if err != nil {
				log.Printf("failed to get build lists for status %v and platform %v", status, platform)
				continue
			}
			defer resp.Body.Close()

			c.handleResponse(resp)
		}
	}
}

func (c *Connector) poller() {
	for {
		select {
		case <-c.clse:
			c.timer.Stop()
			c.wait.Done()
			return
		case <-c.timer.C:
			c.poll()
			c.timer.Reset(c.CheckEvery)
		}
	}
}
