package abf

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/pmylund/go-cache"
	"github.com/robxu9/kahinah/backend"
)

const (
	BaseAPI = "https://abf.io/api/v1"
	BaseWeb = "https://abf.io"

	StatusBuildCompleted        = "0"
	StatusBuildTestingPublished = "12000"

	ResultChangelog = "changelog.log"
)

var (
	ErrNotSuccess = errors.New("connector: did not return 200/SUCCESS")

	EndpointBuildLists   = fmt.Sprintf("%s/%s", BaseWeb, "/build_lists?per_page=100&filter[status]=%v&filter[ownership]=everything&filter[save_to_platform_id]=%v")
	EndpointAPIBuildList = fmt.Sprintf("%s/%s", BaseAPI, "/build_lists/%v.json")
	EndpointWebBuildList = fmt.Sprintf("%s/%s", BaseWeb, "/build_lists/%v")

	EndpointUser = fmt.Sprintf("%s/%s", BaseAPI, "/users/%v.json")

	EndpointBuildListSendTesting = fmt.Sprintf("%s/%s", BaseAPI, "/build_lists/%v/publish_into_testing.json")
)

// Connector is the ABF connector to Kahinah.
// Name: ru.rosalinux.abf
// Info: the ABF id of the update
type Connector struct {
	// Requirements for packages to add
	// Platforms contains a list of platforms ids to monitor
	Platforms []string
	// architectures
	Arches []string

	// basic information
	User   string
	APIKey string
	// Duration setting
	CheckEvery time.Duration

	// kahinah struct
	k *kahinah.K
	// timer for poller
	timer *time.Timer
	// close channel
	clse chan struct{}
	// waitgroup
	wait *sync.WaitGroup
	// http client
	client *http.Client
	// cache
	cache *cache.Cache

	// list of build requests that made it through and are waiting for their
	// counterparts in other architectures
	// Connector speciifc information in the ConnectorStore:
	// 		ARCH: ID
	//		ARCH_WEB: <JSON>
	//		ARCH_API: <JSON>
	waiting map[waitingKey]*kahinah.Update
}

// waitingKey is the key for the waiting map. exactly what it sounds like.
type waitingKey struct {
	Fullname string // full project name
	Hash     string // after hash
}

func (c *Connector) Name() string {
	return "ru.rosalinux.abf"
}

func (c *Connector) Init(k *kahinah.K) error {
	c.k = k
	c.timer = time.NewTimer(c.CheckEvery)
	c.clse = make(chan struct{})
	c.client = &http.Client{}
	c.cache = cache.New(30*time.Minute, 2*time.Hour)

	c.wait = &sync.WaitGroup{}

	c.wait.Add(1)
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
	c.wait.Wait()
}

// reminder: ALWAYS CALL resp.Body.Close()!
func (c *Connector) makeRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	if body == nil {
		req.Header.Add("Content-Length", "0")
	}
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(c.User, c.APIKey)

	listResp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if listResp.StatusCode != 200 {
		return nil, err
	}

	return listResp, nil
}

func (c *Connector) containsPlatform(i int) bool {
	s := strconv.Itoa(i)
	for _, v := range c.Platforms {
		if s == v {
			return true
		}
	}

	return false
}

func (c *Connector) containsArch(s string) bool {
	for _, v := range c.Platforms {
		if s == v {
			return true
		}
	}

	return false
}
