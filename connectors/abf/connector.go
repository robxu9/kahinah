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
	"github.com/robxu9/kahinah/kahinah"
)

const (
	Endpoint = "https://abf.io/api/v1"

	StatusBuildCompleted        = "0"
	StatusBuildTestingPublished = "12000"

	EndpointBuildLists = "/build_lists.json?per_page=100&filter[status]=%v&filter[ownership]=index&filter[platform_id]=%v"
	EndpointBuildList  = "/build_lists/%v.json"

	EndpointUser = "/users/%v.json"

	EndpointBuildListSendTesting = "/build_lists/%v/publish_into_testing.json"
)

// Connector is the ABF connector to Kahinah.
// Name: ru.rosalinux.abf
// Info: the ABF id of the update
type Connector struct {
	// Platforms contains a list of platforms ids to monitor
	Platforms []string
	// Username
	User string
	// API key
	APIKey string
	// Duration setting
	CheckEvery time.Duration

	// kahinah struct
	k *kahinah.Kahinah
	// timer for poller
	timer *time.Timer
	// close channel
	clse chan struct{}
	// http client
	client *http.Client
	// cache
	cache *cache.Cache
}

func (c *Connector) sendTesting(id int) {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s", Endpoint, fmt.Sprintf(EndpointBuildListSendTesting, id)), nil)
	if err != nil {
		// log and die
		log.Printf("failed to construct request for send testing endpoint for id %v: %v", id, err)
	} else {
		req.SetBasicAuth(c.User, c.APIKey)
		req.Header.Add("Content-Length", "0")

		resp, err := c.client.Do(req)
		if err != nil {
			log.Printf("failed to get response: %v", err)
		} else {
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				log.Printf("status code is not 200: %v", resp.StatusCode)
			}
		}
	}
}

func (c *Connector) getUserEmail(id int) string {
	if cached, found := c.cache.Get("user/" + strconv.Itoa(id)); found {
		return cached.(string)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", Endpoint, fmt.Sprintf(EndpointUser, id)), nil)
	if err != nil {
		// log and die
		log.Printf("failed to construct request for user endpoint: %v", err)
		return ""
	}

	req.SetBasicAuth(c.User, c.APIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("failed to get response: %v", err)
		return ""
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("status code is not 200: %v", resp.StatusCode)
		return ""
	}

	decoder := json.NewDecoder(resp.Body)

	var user User

	if err := decoder.Decode(&user); err != nil {
		log.Printf("failed to decode to user: %v", err)
		return ""
	}

	c.cache.Set("user/"+strconv.Itoa(id), user.User.Email, cache.DefaultExpiration)

	return user.User.Email
}

func (c *Connector) handleResponse(resp *http.Response, send_testing bool) {
	decoder := json.NewDecoder(resp.Body)

	var buildlists BuildLists
	if err := decoder.Decode(&buildlists); err != nil {
		log.Printf("failed to decode list of builds: %v", err)
		return
	}

	for _, v := range buildlists.BuildLists {
		// let's not poll twice... check if we already have this entry!
		if ids, _ := c.k.FindUpdatesWithConnector(c.Name(), strconv.Itoa(v.ID), ""); len(ids) != 0 {
			// FIXME do we have to care about the error? (I don't think so?)
			continue
		}

		// get the buildlist otherwise
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", Endpoint, fmt.Sprintf(EndpointBuildList, v.ID)), nil)
		if err != nil {
			log.Printf("failed to construct request to get specific build list: %v", err)
			continue
		}

		req.SetBasicAuth(c.User, c.APIKey)

		listResp, err := c.client.Do(req)
		if err != nil {
			log.Printf("failed to get response from getting specific build list: %v", err)
			continue
		}

		if listResp.StatusCode != 200 {
			log.Printf("status code is not 200: %v", resp.StatusCode)
			continue
		}

		// connectorInfo is the json we read so that we can store a copy in db
		connectorInfo := &bytes.Buffer{}

		listDecoder := json.NewDecoder(io.TeeReader(listResp.Body, connectorInfo))

		var b BuildList
		if err := listDecoder.Decode(&b); err != nil {
			log.Printf("failed to decode specific buildlist: %v", err)
			continue
		}

		updateName := b.Name()
		submitter := c.getUserEmail(b.BuildList.User.ID)
		packages := b.Packages()
		changes := b.Changes()

		// Create the update
		id, err := c.k.NewUpdate(
			c.Name(), // connector
			fmt.Sprintf("%v/%v/%v", b.BuildList.SaveToRepository.Platform.Name, b.BuildList.SaveToRepository.Name, b.BuildList.Arch.Name), // target
			updateName, // name
			submitter,  // submitter
			b.Type(),   // updatetype
			&kahinah.UpdateContent{ // updatecontent
				From:     b.BuildList.LastPublishedCommitHash,
				To:       b.BuildList.CommitHash,
				Url:      b.BuildList.URL,
				BuiltAt:  time.Unix(b.BuildList.CreatedAt, 0),
				Packages: packages,
				Changes:  changes,
			},
			strconv.Itoa(b.BuildList.ID), // connectorid
			connectorInfo.String(),       // connector info
		)

		if err != nil {
			log.Printf("failed to create new update: %v", err)
		} else {
			log.Printf("successfully created update %v from buildlist %v", id, b.BuildList.ID)

			if send_testing {
				c.sendTesting(b.BuildList.ID)
			}
		}

	}
}

func (c *Connector) poll() {
	for _, v := range c.Platforms {
		// construct the url for build complete
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", Endpoint, fmt.Sprintf(EndpointBuildLists, StatusBuildCompleted, v)), nil)
		if err != nil {
			// log and die
			log.Printf("failed to construct request for build complete endpoint for platform %v: %v", v, err)
		} else {
			req.SetBasicAuth(c.User, c.APIKey)

			resp, err := c.client.Do(req)
			if err != nil {
				log.Printf("failed to get response: %v", err)
			} else {
				defer resp.Body.Close()
				if resp.StatusCode != 200 {
					log.Printf("status code is not 200: %v", resp.StatusCode)
				} else {
					c.handleResponse(resp, true)
				}
			}
		}

		// construct the url for testing published
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/%s", Endpoint, fmt.Sprintf(EndpointBuildLists, StatusBuildTestingPublished, v)), nil)
		if err != nil {
			// log and die
			log.Printf("failed to construct request for testing published endpoint for platform %v: %v", v, err)
		} else {
			req.SetBasicAuth(c.User, c.APIKey)

			resp, err := c.client.Do(req)
			if err != nil {
				log.Printf("failed to get response: %v", err)
			} else {
				defer resp.Body.Close()
				if resp.StatusCode != 200 {
					log.Printf("status code is not 200: %v", resp.StatusCode)
				} else {
					c.handleResponse(resp, false)
				}
			}
		}
	}
}

func (c *Connector) poller() {
	for {
		select {
		case <-c.clse:
			c.timer.Stop()
			return
		case <-c.timer.C:
			c.poll()
			c.timer.Reset(c.CheckEvery)
		}
	}
}

func (c *Connector) Name() string {
	return "ru.rosalinux.abf"
}

func (c *Connector) Init(k *kahinah.Kahinah) error {
	c.k = k
	c.timer = time.NewTimer(c.CheckEvery)
	c.clse = make(chan struct{})
	c.client = &http.Client{}
	c.cache = cache.New(30*time.Minute, 2*time.Hour)

	go c.poller()
	return nil
}

func (c *Connector) Pass(u *kahinah.Update) {

}

func (c *Connector) Fail(u *kahinah.Update) {

}

func (c *Connector) Close() {
	c.clse <- struct{}{}
	close(c.clse)
}
