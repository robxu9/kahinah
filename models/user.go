package models

import "github.com/jinzhu/gorm"

const (
	UserSystem = "[System]"
)

type User struct {
	gorm.Model

	Username    string
	Email       string
	Integration string

	APIKey string `json:"-"`

	Permissions []UserPermission
	Activities  []ListActivity
}

func (u *User) String() string {
	return u.Username
}

func (u *User) Save() {
	DB.Save(u)
}

func (u *User) TableName() string {
	return DBPrefix + "users"
}

type UserPermission struct {
	gorm.Model

	Permission string
	UserID     uint
}

func (u *UserPermission) Save() {
	DB.Save(u)
}

func (u *UserPermission) TableName() string {
	return DBPrefix + "userpermissions"
}

func init() {
	// make sure we have a system user
	FindUser(UserSystem)
}

// FindUserByID attempts to find the user with the given ID; if it doesn't
// exist, it returns nil.
func FindUserByID(id uint) *User {
	var user User

	if err := DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}

	return &user
}

// FindUser finds the user with the given username. If a user doesn't exist,
// it creates one and returns it.
func FindUser(username string) *User {
	var user User

	if err := DB.First(&user, &User{Username: username}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			user = User{
				Username: username,
			}
			if err = DB.Create(&user).Error; err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	return &user
}

func FindUserByAPI(apikey string) *User {
	var user User

	if err := DB.First(&user, &User{APIKey: apikey}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}

	return &user
}

func FindUserNoCreate(username string) *User {
	var user User

	if err := DB.First(&user, &User{Username: username}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}

	return &user
}
