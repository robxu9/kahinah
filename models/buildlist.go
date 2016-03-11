package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	ListRunning  = "running"
	ListPending  = "pending"
	ListAccepted = "accepted"
	ListRejected = "rejected"

	LinkMain      = "_mainURL"
	LinkChangelog = "_changelogURL"
	LinkSCM       = "_scmlogURL"

	ListStageNotStarted = "_stage_pending"
	ListStageFinished   = "_stage_finished"

	ArtifactBinary = "binary"
	ArtifactSource = "source"
)

type List struct {
	gorm.Model

	// -- basic information
	Name     string // project name
	Platform string // targeted platform
	Channel  string // targeted sub-platform/repository
	Variants string // variants/architectures (split by ';')

	// -- contained data
	Artifacts  []*ListArtifact
	Links      []*ListLink
	Changes    string // provide a field which describes changes, whether textual or diff
	BuildDate  time.Time
	AdvisoryID uint // link to advisory

	// -- current stage and activity
	StageCurrent string // either NotStarted, Finished, or the stage defined in between
	Activity     []*ListActivity

	// -- current stages
	Stages []*ListStage

	// -- last but not least, integration stuff (five string slots for usage)
	IntegrationOne   string
	IntegrationTwo   string
	IntegrationThree string
	IntegrationFour  string
	IntegrationFive  string
}

func (l *List) TableName() string {
	return DBPrefix + "lists"
}

func (l *List) AddActivity(u *User, activity string) {
	newActivity := &ListActivity{
		ListID:   l.ID,
		UserID:   u.ID,
		Activity: activity,
	}
	if err := DB.Create(newActivity).Error; err != nil {
		panic(err) // this shouldn't really panic though..
	}
}

type ListArtifact struct {
	gorm.Model

	ListID uint // link to List

	Name    string
	Type    string
	Epoch   int64
	Version string
	Release string
	Variant string

	URL string
}

func (l *ListArtifact) TableName() string {
	return DBPrefix + "listartifacts"
}

type ListLink struct {
	gorm.Model

	ListID uint // link to List

	Name string
	URL  string
}

func (l *ListLink) TableName() string {
	return DBPrefix + "listlinks"
}

type ListActivity struct {
	gorm.Model

	ListID uint // link to List
	UserID uint // link to User

	Activity string // markdown activity comment
}

func (l *ListActivity) TableName() string {
	return DBPrefix + "listactivities"
}

type ListStage struct {
	gorm.Model

	ListID uint   // link to List
	Name   string // stage name

	// -- processes
	Processes []*ListStageProcess
}

func (l *ListStage) TableName() string {
	return DBPrefix + "liststages"
}

type ListStageProcess struct {
	gorm.Model

	ListStageID uint   // link to ListStage
	Name        string // process name
	Optional    bool   // if this stage is okay to fail

	Data []byte // blob data
}

func (l *ListStageProcess) ParentList() *List {
	// find ListStage...
	var listStage ListStage
	if err := DB.First(&listStage, l.ListStageID).Error; err != nil {
		panic(err)
	}
	// then find List itself
	var list List
	if err := DB.First(&list, listStage.ListID).Error; err != nil {
		panic(err)
	}

	return &list
}

func (l *ListStageProcess) TableName() string {
	return DBPrefix + "liststageprocesses"
}
