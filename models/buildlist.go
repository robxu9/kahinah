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
	checkAllListStagesMutex = &sync.RWMutex{}
	checkAllListStagesMap   = map[uint]bool{}
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

// CheckAllListStages retrieves all stages that are not on the "success" or
// "failed" state. For each list, it checks the id (to make sure another
// goroutine is not already checking it), then fires off a goroutine to
// handle the check.
func CheckAllListStages() {
	// first, get the list of all stages that need checking.
	var lists []List
	if err := DB.Where("stage_result in (?)", []string{ListRunning, ListPending}).Find(&lists).Error; err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	// now go through each list, check if it's not already being checked,
	// and fire it off if it's not.
	checkAllListStagesMutex.RLock()
	defer checkAllListStagesMutex.RUnlock()

	for _, v := range lists {
		if !checkAllListStagesMap[v.ID] {
			go func() {
				checkAllListStagesMutex.Lock()

				// check if there's already an existing goroutine checking
				if checkAllListStagesMap[v.ID] {
					checkAllListStagesMutex.Unlock()
					return
				}

				checkAllListStagesMap[v.ID] = true
				checkAllListStagesMutex.Unlock()

				defer func() {
					checkAllListStagesMutex.Lock()
					defer checkAllListStagesMutex.Unlock()
					delete(checkAllListStagesMap, v.ID)
				}()

				v.CheckStage()
			}()
		}
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
