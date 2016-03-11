package models

import (
	"fmt"

	// mysql functionality
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	// postgresql funcationality
	_ "github.com/lib/pq"
	// sqlite3 functionality
	_ "github.com/mattn/go-sqlite3"
	"github.com/robxu9/kahinah/conf"
	"github.com/robxu9/kahinah/log"
)

var (
	DBPrefix = conf.Config.GetDefault("database.prefix", "kh_").(string)
	DBType   = conf.Config.GetDefault("database.type", "sqlite3").(string)
	DBName   = conf.Config.GetDefault("database.name", "data.sqlite").(string)
	DBHost   = conf.Config.GetDefault("database.host", "localhost:3306").(string)
	DBUser   = conf.Config.GetDefault("database.user", "root").(string)
	DBPass   = conf.Config.GetDefault("database.pass", "toor").(string)
	DBDebug  = conf.Config.GetDefault("database.debug", false).(bool)

	DB *gorm.DB
)

func init() {
	var db gorm.DB
	var err error
	switch DBType {
	case "mysql":
		db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", DBUser, DBPass, DBHost, DBName), 30)
	case "sqlite3":
		db, err = gorm.Open("sqlite3", "file:"+DBName, 30)
	case "postgres":
		db, err = gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s", DBUser, DBPass, DBHost, DBName))
	default:
		log.Logger.Fatalf("db: I don't know how to handle %v", DBType)
	}
	DB = &db

	DB.LogMode(DBDebug)
	DB.SetLogger(gorm.Logger{log.Logger})
	if err = DB.DB().Ping(); err != nil {
		log.Logger.Fatalf("db: couldn't ping the database: %v", err)
	}

	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)

	if err = DB.AutoMigrate(&Advisory{}, &List{}, &ListActivity{},
		&ListArtifact{}, &ListLink{}, &ListStage{}, &ListStageProcess{},
		&User{}, &UserPermission{}).Error; err != nil {
		log.Logger.Fatalf("db: failed to automigrate: %v", err)
	}
}
