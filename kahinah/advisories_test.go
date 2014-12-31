package kahinah

import "testing"

func TestAdvisory(t *testing.T) {
	k := setupTest(t)
	defer k.Close()

	c := setupConnector()
	k.Attach(c)

	id := c.MakeNewUpdate(t)

	user, err := k.NewUser("test@example.com")
	if err != nil {
		t.Fatalf("failed to create new user: %v", err)
	}

	advisoryId, err := k.NewAdvisory(user, []int64{id}, []string{"bugzillaURL"}, "this is a advisory on stupidity", "ROBXU9")
	if err != nil {
		t.Fatalf("failed to create new advisory: %v", err)
	}

	if count := k.CountAdvisories(); count != 1 {
		t.Fatalf("counting advisories failed: should be 1 but is %d", count)
	}

	list, err := k.ListAdvisories(0, 10)
	if err != nil {
		t.Fatalf("failed to list advisories: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("list of advisories is not 1")
	}

	if list[0] != advisoryId {
		t.Fatalf("list[0] should be id which is %d but is %d", id, list[0])
	}

	// first retrieve the cache'd update and make sure it's okay
	advisory, err := k.RetrieveAdvisory(advisoryId)
	if err != nil {
		t.Fatalf("should not have failed to retrieve advisory from cache: %v", err)
	}

	if advisory.Version != 0 {
		t.Fatalf("advisory version is not 0: %v", advisory.Version)
	}

	if advisory.Status != OPEN {
		t.Fatalf("advisory is not open: %v", advisory.Status)
	}

	if advisory.AdvisoryFamily != "ROBXU9" {
		t.Fatalf("advisory family is not the same: %v", advisory.AdvisoryFamily)
	}

	// now flush the cache
	k.cache.Flush()

	// and try again
	advisory2, err := k.RetrieveAdvisory(advisoryId)
	if err != nil {
		t.Fatalf("should not have failed to retrieve advisory from cache: %v", err)
	}

	// it's not the same as the cache, right?
	if advisory == advisory2 {
		t.Fatalf("should not be pulling from cache")
	}

	if advisory2.Version != 0 {
		t.Fatalf("advisory version is not 0: %v", advisory.Version)
	}

	if advisory2.Status != OPEN {
		t.Fatalf("advisory is not open: %v", advisory.Status)
	}

	if advisory2.AdvisoryFamily != "ROBXU9" {
		t.Fatalf("advisory family is not the same: %v", advisory.AdvisoryFamily)
	}

}
