package kahinah

import (
	"errors"
	"time"
)

var (
	ErrNoSuchUser = errors.New("kahinah: no such user")
)

type User struct {
	Id int64

	Name  string
	Email string

	Advisories []*Advisory `sql:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type UserIP struct {
	Id int64

	UserId     int64
	IP         string
	AccessedAt time.Time
}

type UserToken struct {
	Id int64

	UserId  int64
	Name    string
	Token   string
	Revoked bool

	CreatedAt time.Time
}

func (k *Kahinah) RegisterUser(name, email string) (int64, error) {
	return 0, nil
}

func (k *Kahinah) FindUser(email string) (int64, error) {
	return 0, nil
}

func (k *Kahinah) RetrieveUser(id int64) (*User, error) {
	return nil, nil
}
