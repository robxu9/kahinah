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
	"github.com/robxu9/kahinah/common/conf"
	"github.com/robxu9/kahinah/common/klog"
)

var (
	DBPrefix = conf.GetDefault("database.prefix", "kh_").(string)
	DBType   = conf.GetDefault("database.type", "sqlite3").(string)
	DBName   = conf.GetDefault("database.name", "data.sqlite").(string)
	DBHost   = conf.GetDefault("database.host", "localhost:3306").(string)
	DBUser   = conf.GetDefault("database.user", "root").(string)
	DBPass   = conf.GetDefault("database.pass", "toor").(string)
	DBDebug  = conf.GetDefault("database.debug", false).(bool)

	DB *gorm.DB
)

func init() {
	var err error
	switch DBType {
	case "mysql":
		DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", DBUser, DBPass, DBHost, DBName))
	case "sqlite3":
		DB, err = gorm.Open("sqlite3", "file:"+DBName)
	case "postgres":
		DB, err = gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s", DBUser, DBPass, DBHost, DBName))
	default:
		klog.Fatalf("db: I don't know how to handle %v", DBType)
	}

	DB.LogMode(DBDebug)
	DB.SetLogger(gorm.Logger{LogWriter: klog.Logger})
	if err = DB.DB().Ping(); err != nil {
		klog.Fatalf("db: couldn't ping the database: %v", err)
	}

	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)

	if err = DB.AutoMigrate(&Advisory{}, &List{}, &ListActivity{},
		&ListArtifact{}, &ListLink{}, &ListStage{}, &ListStageProcess{},
		&User{}, &UserPermission{}).Error; err != nil {
		klog.Fatalf("db: failed to automigrate: %v", err)
	}
}
