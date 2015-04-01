package kahinah

import (
	"testing"
	"time"
)

type TestingConnector struct {
	k    *Kahinah
	init chan struct{}
	pass chan struct{}
	fail chan struct{}
	clse chan struct{}
}

func setupConnector() *TestingConnector {
	return &TestingConnector{
		init: make(chan struct{}, 1),
		pass: make(chan struct{}, 1),
		fail: make(chan struct{}, 1),
		clse: make(chan struct{}, 1),
	}
}

func (t *TestingConnector) MakeNewUpdate(test *testing.T) int64 {
	id, err := t.k.NewUpdate(t.Name(), "robxu9/2014/main", "test-1.0.0/amd64", "joe@example.com", BUGFIX, &UpdateContent{
		From:    "120",
		To:      "126",
		Url:     "http://example.com/",
		BuiltAt: time.Now(),
		Packages: []*UpdatePackage{
			&UpdatePackage{
				Name:    "test",
				Epoch:   0,
				Version: "1.0.0",
				Release: "1.robxu9",
				Arch:    "amd64",
				Type:    "src",
				Url:     "http://example.com/test.tar.xz",
			},
			&UpdatePackage{
				Name:    "test",
				Epoch:   0,
				Version: "1.0.0",
				Release: "1.robxu9",
				Arch:    "amd64",
				Type:    "binary",
				Url:     "http://example.com/test.pkg",
			},
		},
		Changes: []*UpdateChange{
			&UpdateChange{
				ChangeAt: time.Now(),
				For:      "1.0.0-1.robxu9",
				By:       "joe@example.com",
				Details:  "did some stuff",
			},
		},
	}, "1234567890", "info here")

	if err != nil {
		test.Fatal(err)
	}

	return id
}

func (t *TestingConnector) MakeNewUpdate2(test *testing.T) int64 {
	id, err := t.k.NewUpdate(t.Name(), "robxu9/2015/main", "test-1.0.0/amd64", "joe@example.com", BUGFIX, &UpdateContent{
		From:    "120",
		To:      "126",
		Url:     "http://example.com/",
		BuiltAt: time.Now(),
		Packages: []*UpdatePackage{
			&UpdatePackage{
				Name:    "test",
				Epoch:   0,
				Version: "1.0.0",
				Release: "1.robxu9",
				Arch:    "amd64",
				Type:    "src",
				Url:     "http://example.com/test.tar.xz",
			},
			&UpdatePackage{
				Name:    "test",
				Epoch:   0,
				Version: "1.0.0",
				Release: "1.robxu9",
				Arch:    "amd64",
				Type:    "binary",
				Url:     "http://example.com/test.pkg",
			},
		},
		Changes: []*UpdateChange{
			&UpdateChange{
				ChangeAt: time.Now(),
				For:      "1.0.0-1.robxu9",
				By:       "joe@example.com",
				Details:  "did some stuff",
			},
		},
	}, "1234567891", "info here")

	if err != nil {
		test.Fatal(err)
	}

	return id
}

func (t *TestingConnector) Name() string {
	return "com.robxu9.kahinah.connector.test"
}

func (t *TestingConnector) Init(k *Kahinah) error {
	t.k = k
	t.init <- struct{}{}
	return nil
}

func (t *TestingConnector) Pass(u *Update) {
	t.pass <- struct{}{}
}

func (t *TestingConnector) Fail(u *Update) {
	t.fail <- struct{}{}
}

func (t *TestingConnector) Close() {
	t.clse <- struct{}{}
}

func TestConnector(t *testing.T) {
	k := setupTest(t)
	defer k.Close()

	c := setupConnector()
	if k.Attach(c) != nil {
		t.Fatal("failed to attach to connector")
	}

	if len(c.init) != 1 {
		t.Fatal("failed to init connector with method")
	}

	if c.k != k {
		t.Fatal("failed to init connector with k")
	}
}
