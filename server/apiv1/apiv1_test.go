package apiv1

import (
	"bytes"
	"net"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/robxu9/kahinah/kahinah"
	"github.com/robxu9/kahinah/server/common"
	"github.com/ryanuber/go-glob"
)

type MatchType int

const (
	// MatchNone does no matching on output
	MatchNone MatchType = iota
	// MatchEqual does string equal matching on output
	MatchEqual
	// MatchRegex does string regex matching on output
	MatchRegex
	// MatchGlob does string glob matching on output
	MatchGlob
	// MatchShell does string shell matching on output
	MatchShell
)

type APIv1Test struct {
	t        *testing.T
	handler  http.Handler
	c        *common.Common
	testUser int64
}

func setupAPIv1(t *testing.T) *APIv1Test {
	conf := common.DefaultConfig(0)
	conf.DebugMode = testing.Verbose()
	conf.Database.Dialect = "sqlite3"
	conf.Database.Params = ":memory:"

	c := common.Open(conf)
	handler := New(c)

	// create a user for tokenizing later
	testUser, err := c.K.NewUser("test@example.com")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return &APIv1Test{
		t:        t,
		handler:  handler,
		c:        c,
		testUser: testUser,
	}
}

func (a *APIv1Test) NetTest(methodIn, endpointIn, bodyIn string, authIn bool, statusOut int, bodyOut string, matchOut MatchType) string {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 0,
	})
	if err != nil {
		a.t.Fatalf("failed to setup listener: %v", err)
	}
	defer listener.Close()

	go http.Serve(listener, a.handler)

	request, err := http.NewRequest(methodIn, "http://"+listener.Addr().String()+"/"+endpointIn, strings.NewReader(bodyIn))
	if err != nil {
		a.t.Fatalf("failed to setup request: %v", err)
	}

	if authIn {
		// create a usertoken for that user
		_, tokenstr, err := a.c.GenerateUserToken(a.testUser, "test token", true)
		if err != nil {
			a.t.Fatalf("failed to setup token for auth'ed request: %v", err)
		}

		// use that usertoken
		request.Header.Set("Authorization", "Bearer "+tokenstr)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		a.t.Fatalf("failed to do request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != statusOut {
		a.t.Fatalf("statuscode response is not the same: want %v, got %v", statusOut, response.StatusCode)
	}

	bodyBuffer := &bytes.Buffer{}
	bodyBuffer.ReadFrom(response.Body)

	switch matchOut {
	case MatchEqual:
		if bodyBuffer.String() != bodyOut {
			a.t.Fatalf("body response is not the same: want %v, got %v", bodyOut, bodyBuffer.String())
		}
	case MatchRegex:
		matched, err := regexp.MatchString(bodyOut, bodyBuffer.String())
		if err != nil {
			a.t.Fatalf("failed to compile regexp: %v", err)
		}
		if !matched {
			a.t.Fatalf("body response failed to match regex: want %v, got %v", bodyOut, bodyBuffer.String())
		}
	case MatchGlob:
		if !glob.Glob(bodyOut, bodyBuffer.String()) {
			a.t.Fatalf("body response failed to match glob: want %v, got %v", bodyOut, bodyBuffer.String())
		}
	case MatchShell:
		matched, err := filepath.Match(bodyOut, bodyBuffer.String())
		if err != nil {
			a.t.Fatalf("failed to compile shell glob: %v", err)
		}
		if !matched {
			a.t.Fatalf("body response failed to match shell: want %v, got %v", bodyOut, bodyBuffer.String())
		}
	}

	return bodyBuffer.String()
}

func (a *APIv1Test) MakeUpdate() int64 { // similar to connector_test.go from kahinah
	id, err := a.c.K.NewUpdate("test", "robxu9/2014/main", "test-1.0.0/amd64", "test@example.com", kahinah.BUGFIX, &kahinah.UpdateContent{
		From:    "120",
		To:      "126",
		Url:     "http://example.com/",
		BuiltAt: time.Unix(1398902400, 0),
		Packages: []*kahinah.UpdatePackage{
			&kahinah.UpdatePackage{
				Name:    "test",
				Epoch:   0,
				Version: "1.0.0",
				Release: "1.robxu9",
				Arch:    "amd64",
				Type:    "src",
				Url:     "http://example.com/test.tar.xz",
			},
			&kahinah.UpdatePackage{
				Name:    "test",
				Epoch:   0,
				Version: "1.0.0",
				Release: "1.robxu9",
				Arch:    "amd64",
				Type:    "binary",
				Url:     "http://example.com/test.pkg",
			},
		},
		Changes: []*kahinah.UpdateChange{
			&kahinah.UpdateChange{
				ChangeAt: time.Unix(1398902400, 0),
				For:      "1.0.0-1.robxu9",
				By:       "test@example.com",
				Details:  "did some stuff",
			},
		},
	}, "1234567890", "info here")

	if err != nil {
		a.t.Fatal(err)
	}

	return id
}

func (a *APIv1Test) Close() {
	a.c.Close()
}
