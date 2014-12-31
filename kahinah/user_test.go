package kahinah

import "testing"

func TestUser(t *testing.T) {
	k := setupTest(t)
	defer k.Close()

	if _, err := k.FindUser("test@example.com"); err == nil {
		t.Fatal("should not be nil")
	}

	id, err := k.NewUser("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	find, err := k.FindUser("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	if find != id {
		t.Fatalf("id should be same as found")
	}

	user, err := k.RetrieveUser(id)
	if err != nil {
		t.Fatalf("retrieving user should not fail")
	}

	if user.Email != "test@example.com" {
		t.Fatalf("email failed to store")
	}
}
