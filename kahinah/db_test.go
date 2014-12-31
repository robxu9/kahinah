package kahinah

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTest(t *testing.T) *Kahinah {
	// mimic with sqlite3 db
	k, err := open("sqlite3", ":memory:", testing.Verbose())
	if err != nil {
		t.Fatal(err)
	}

	return k
}

func TestSetup(t *testing.T) {
	k := setupTest(t)
	defer k.Close()
	// okay, well... that's it really
}
