package models

import "time"

const (
	AdvisoryBugfix   = "bugfix"
	AdvisorySecurity = "security"
	AdvisoryFeature  = "feature"
)

type Advisory struct {
	// we can't use gorm.Model because of ID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Dialect    string `gorm:"primary_key"`
	Year       string `gorm:"primary_key"`
	AdvisoryID uint   `gorm:"primary_key" sql:"AUTO_INCREMENT"`

	Type        string
	Summary     string
	Description string

	Lists []List // lists
}

func (a *Advisory) TableName() string {
	return DBPrefix + "advisories"
}
