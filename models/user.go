package models

import "github.com/astaxie/beego/orm"

const (
	UserSystem = "[System]"
)

type User struct {
	Id          uint64 `orm:"auto;pk"`
	Username    string `orm:"type(text)"`
	Email       string `orm:"type(text)"`
	Integration string `orm:"type(text)"` // abf service user id

	APIKey string `orm:"type(text)" json:"-"`

	Permissions []*UserPermission `orm:"rel(m2m);on_delete(set_null)"`
	Karma       []*Karma          `orm:"reverse(many);on_delete(set_null)"`

	BuildLists []*BuildList `orm:"null;reverse(many);on_delete(set_null)"`
	Advisories []*Advisory  `orm:"null;reverse(many);on_delete(set_null)"`
}

func (u *User) String() string {
	return u.Username
}

func (u *User) Save() {
	o := orm.NewOrm()
	o.Update(u)
}

type UserPermission struct {
	Id         uint64  `orm:"auto;pk"`
	Permission string  `orm:"type(text);unique"`
	Users      []*User `orm:"null;reverse(many);on_delete(set_null)"`
}

func (u *UserPermission) Save() {
	o := orm.NewOrm()
	o.Update(u)
}

func init() {
	// make sure we have a system user
	FindUser(UserSystem)
}

// Finds the user with the given username. If a user doesn't exist,
// it creates one and returns it.
func FindUser(username string) *User {
	o := orm.NewOrm()
	qt := o.QueryTable(new(User))

	var user User
	err := qt.Filter("Username", username).One(&user)
	if err != nil && err != orm.ErrNoRows {
		panic(err)
	} else if err != nil {
		user = User{
			Username: username,
		}
		o.Insert(&user)
	} else {
		o.LoadRelated(&user, "Permissions")
		o.LoadRelated(&user, "Karma")
		o.LoadRelated(&user, "BuildLists")
	}

	return &user
}

func FindUserAPI(apikey string) *User {
	o := orm.NewOrm()
	qt := o.QueryTable(new(User))

	var user User
	err := qt.Filter("APIKey", apikey).One(&user)
	if err != nil && err != orm.ErrNoRows {
		panic(err)
	} else if err != nil {
		// No such User
		return nil
	} else {
		o.LoadRelated(&user, "Permissions")
		o.LoadRelated(&user, "Karma")
		o.LoadRelated(&user, "BuildLists")
	}

	return &user
}

func FindUserNoCreate(username string) *User {
	o := orm.NewOrm()
	qt := o.QueryTable(new(User))

	var user User
	err := qt.Filter("Username", username).One(&user)
	if err != nil && err != orm.ErrNoRows {
		panic(err)
	} else if err != nil {
		return nil
	} else {
		o.LoadRelated(&user, "Permissions")
		o.LoadRelated(&user, "Karma")
		o.LoadRelated(&user, "BuildLists")
	}

	return &user
}

func PermRegister(perm string) {
	o := orm.NewOrm()
	qt := o.QueryTable(new(UserPermission))

	num, err := qt.Filter("Permission", perm).Count()
	if err != nil && err != orm.ErrNoRows {
		panic(err)
	}

	if num > 0 {
		return
	}

	permission := UserPermission{Permission: perm}
	_, err = o.Insert(&permission)
	if err != nil {
		panic(err)
	}
}

func PermGet(perm string) *UserPermission {
	o := orm.NewOrm()
	qt := o.QueryTable(new(UserPermission))

	var p UserPermission
	err := qt.Filter("Permission", perm).One(&p)
	if err != nil && err != orm.ErrNoRows {
		panic(err)
	} else if err != nil {
		return nil
	}

	return &p
}

func PermGetAll() []*UserPermission {
	o := orm.NewOrm()
	qt := o.QueryTable(new(UserPermission))

	var perms []*UserPermission
	_, err := qt.All(&perms)

	if err != nil && err != orm.ErrNoRows {
		panic(err)
	}

	for _, v := range perms {
		o.LoadRelated(v, "Users")
	}

	return perms
}
