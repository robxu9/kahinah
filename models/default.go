package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var (
	Prefix = beego.AppConfig.String("db_prefix")
	Name   = beego.AppConfig.String("db_name")
)

func init() {
	orm.Debug, _ = beego.AppConfig.Bool("orm.debug")
	orm.DefaultTimeLoc = time.UTC

	orm.RegisterModelWithPrefix(Prefix, new(BuildList))
	orm.RegisterModelWithPrefix(Prefix, new(Karma))

	orm.RegisterDataBase("default", "sqlite3", "file:"+Name, 30)
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		panic(err)
	}
}
