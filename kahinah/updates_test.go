package kahinah

import "testing"

func TestUpdate(t *testing.T) {
	k := setupTest(t)
	defer k.Close()

	c := setupConnector()
	k.Attach(c)

	id := c.MakeNewUpdate(t)

	if count := k.CountUpdates(); count != 1 {
		t.Fatalf("counting updates failed: should be 1 but is %d", count)
	}

	list, err := k.ListUpdates(0, 10)
	if err != nil {
		t.Fatalf("failed to list updates: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("list of updates is not 1")
	}

	if list[0] != id {
		t.Fatalf("list[0] should be id which is %d but is %d", id, list[0])
	}

	// first retrieve the cache'd update and make sure it's okay
	update, err := k.RetrieveUpdate(id)
	if err != nil {
		t.Fatalf("should not have failed to retrieve update from cache: %v", err)
	}

	if len(update.Content.Packages) != 2 {
		t.Fatalf("packages is not 2")
	}

	if len(update.Content.Changes) != 1 {
		t.Fatalf("changes is not 1")
	}

	if update.Connector != c.Name() {
		t.Fatalf("connector is not the same: %v", update.Connector)
	}

	// now flush the cache
	k.cache.Flush()

	// and try again
	update2, err := k.RetrieveUpdate(id)
	if err != nil {
		t.Fatalf("failed to cold retrieve update: %v", err)
	}

	// it's not the same as the cache, right?
	if update == update2 {
		t.Fatalf("should not be pulling from cache")
	}

	if update2.Id != update.Id {
		t.Fatalf("id is not the same")
	}

	if len(update2.Content.Packages) != 2 {
		t.Fatalf("packages is not 2")
	}

	if len(update2.Content.Changes) != 1 {
		t.Fatalf("changes is not 1")
	}

	if update2.Connector != c.Name() {
		t.Fatalf("connector is not the same: %v", update2.Connector)
	}

	t.Log(update2)

	// FindUpdatesWithConnector
	updates, err := k.FindUpdatesWithConnector(update2.Connector, update2.ConnectorId, update2.ConnectorInfo)
	if err != nil {
		t.Fatalf("findupdateswithconnector failed: %v", err)
	}

	if len(updates) != 1 {
		t.Fatalf("num updates should be 1")
	}

	if updates[0] != update2.Id {
		t.Fatalf("updates is returning wrong id: %v", updates[0])
	}
}
