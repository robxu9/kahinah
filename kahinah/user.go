package kahinah

import (
	"errors"
	"time"
)

var (
	// ErrUserExists - a user with that email exists
	ErrUserExists = errors.New("kahinah: user with that email exists")
	// ErrNoSuchUser - no such user
	ErrNoSuchUser = errors.New("kahinah: no such user")
)

// User represents a user with an email address.
type User struct {
	Id int64

	Email string

	Advisories []int64 `sql:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// NewUser creates a new user with an email address.
func (k *Kahinah) NewUser(email string) (int64, error) {
	k.userMutex.Lock()
	defer k.userMutex.Unlock()

	if _, err := k.FindUser(email); err == nil {
		return 0, ErrUserExists
	}

	user := &User{
		Email: email,
	}

	if err := k.db.Save(email).Error; err != nil {
		return 0, err
	}

	return user.Id, nil
}

// FindUser looks for a user with a specified email.
func (k *Kahinah) FindUser(email string) (int64, error) {
	var id int64

	if err := k.db.Where(&User{Email: email}).Pluck("id", &id).Error; err != nil {
		return 0, ErrNoSuchUser
	}

	return id, nil
}

// RetrieveUser retrieves a user with the specified id.
func (k *Kahinah) RetrieveUser(id int64) (*User, error) {
	record := &User{}

	if k.db.First(record, id).RecordNotFound() {
		return nil, ErrNoSuchUser
	}

	// get a list of all advisories
	if err := k.db.Where(&Advisory{UserId: record.Id}).Pluck("id", &record.Advisories).Error; err != nil {
		panic(err)
	}

	return record, nil
}
