package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robxu9/kahinah/conf"
)

var (
	DBPrefix = conf.Config.GetDefault("database.prefix", "kh_").(string)
	DBType   = conf.Config.GetDefault("database.type", "sqlite3").(string)
	DBName   = conf.Config.GetDefault("database.name", "data.sqlite").(string)
	DBHost   = conf.Config.GetDefault("database.host", "localhost:3306").(string)
	DBUser   = conf.Config.GetDefault("database.user", "root").(string)
	DBPass   = conf.Config.GetDefault("database.pass", "toor").(string)
	DBDebug  = conf.Config.GetDefault("database.debug", false).(bool)
)

func init() {
	orm.Debug = DBDebug
	orm.DefaultTimeLoc = time.Local
	//orm.DefaultTimeLoc = time.UTC

	orm.RegisterModelWithPrefix(DBPrefix, new(BuildList))
	orm.RegisterModelWithPrefix(DBPrefix, new(BuildListPkg))
	orm.RegisterModelWithPrefix(DBPrefix, new(BuildListLink))
	orm.RegisterModelWithPrefix(DBPrefix, new(Karma))

	orm.RegisterModelWithPrefix(DBPrefix, new(Advisory))

	orm.RegisterModelWithPrefix(DBPrefix, new(User))
	orm.RegisterModelWithPrefix(DBPrefix, new(UserPermission))

	switch DBType {
	case "mysql":
		orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", DBUser, DBPass, DBHost, DBName), 30)
	case "sqlite3":
		orm.RegisterDataBase("default", "sqlite3", "file:"+DBName, 30)
	case "postgres":
		orm.RegisterDataBase("default", "postgres", fmt.Sprintf("postgres://%s:%s@%s/%s", DBUser, DBPass, DBHost, DBName))
	}

	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		panic(err)
	}
}
