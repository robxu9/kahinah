package models

import (
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	// ListRunning means that the list is running QA tasks
	ListRunning = "running"
	// ListPending means that the list is pending user intervention
	ListPending = "pending"
	// ListSuccess means that all tasks completed successfully
	ListSuccess = "success"
	// ListFailed means that the QA tasks failed to complete
	ListFailed = "failed"

	LinkMain      = "_mainURL"
	LinkChangelog = "_changelogURL"
	LinkSCM       = "_scmlogURL"

	ListStageNotStarted = "_stage_pending"
	ListStageFinished   = "_stage_finished"

	ArtifactBinary = "binary"
	ArtifactSource = "source"
)

var (
	checkAllListStagesChan  = make(chan int, 1)
	checkAllListStagesMutex = &sync.Mutex{}
)

type List struct {
	gorm.Model

	// -- basic information
	Name     string // project name
	Platform string // targeted platform
	Channel  string // targeted sub-platform/repository
	Variants string // variants/architectures (split by ';')

	// -- contained data
	Artifacts  []ListArtifact
	Links      []ListLink
	Changes    string // provide a field which describes changes, whether textual or diff
	BuildDate  time.Time
	AdvisoryID uint // link to advisory

	// -- current stage and activity
	StageCurrent string // either NotStarted, Finished, or the stage defined in between
	StageResult  string // the status of the stage (running, pending, passed, failed)
	// (this will only be set from pending or running when the whole list has finished or a stage has failed)
	Activity []ListActivity

	// -- current stages (populated during StageNotStarted)
	PlatformGitConfig string      // where we read our configuration from
	Stages            []ListStage // the defined stages

	// -- last but not least, integration stuff (five string slots for usage)
	IntegrationName  string
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

// CheckStage checks the current stage of the list for completion. If so,
// it pushes the list to the next stage and fires off tasks.
// If the present stage is StageNotStarted, it populates the rest of the stages.
func (l *List) CheckStage() {
	// TODO: implement
}

// CheckAllListStages retrieves all stages that are not
func CheckAllListStages() {
	if len(checkAllListStagesChan) >= 1 {
		return // we only want one at a time since this runs so frequently,
		// and we really don't want 9000 goroutines waiting on a mutex.
	}
	checkAllListStagesChan <- 1
	checkAllListStagesMutex.Lock() // double lock b/c len() >= 1 is imperfect
	defer checkAllListStagesMutex.Unlock()
	// TODO: implement
	// Get all stages not currently ListSuccess or ListFailed
	<-checkAllListStagesChan // relieve it
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

func (l *ListActivity) MailActivity() {
	// TODO: implement
}

type ListStage struct {
	gorm.Model

	ListID uint   // link to List
	Name   string // stage name

	// -- processes
	Processes []ListStageProcess
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
