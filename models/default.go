package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var (
	Prefix = beego.AppConfig.String("db_prefix")
	DbType = beego.AppConfig.String("db_type")
	DbName = beego.AppConfig.String("db_name")
	DbHost = beego.AppConfig.String("db_host")
	DbUser = beego.AppConfig.String("db_user")
	DbPass = beego.AppConfig.String("db_pass")
)

func init() {
	orm.Debug, _ = beego.AppConfig.Bool("orm.debug")
	orm.DefaultTimeLoc = time.UTC

	orm.RegisterModelWithPrefix(Prefix, new(BuildList))
	orm.RegisterModelWithPrefix(Prefix, new(BuildListPkg))

	orm.RegisterModelWithPrefix(Prefix, new(Karma))

	orm.RegisterModelWithPrefix(Prefix, new(User))
	orm.RegisterModelWithPrefix(Prefix, new(UserPermission))

	if DbType == "mysql" {
		orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", DbUser, DbPass, DbHost, DbName), 30)
	} else {
		orm.RegisterDataBase("default", "sqlite3", "file:"+DbName, 30)

	}
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		panic(err)
	}
}
