package controllers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"goji.io/pat"

	"golang.org/x/net/context"

	"github.com/jinzhu/gorm"
	"github.com/robxu9/kahinah/data"
	"github.com/robxu9/kahinah/models"
	"github.com/robxu9/kahinah/processes"
	"menteslibres.net/gosexy/to"
)

var (
	// ErrNoCurrentStage signals that StageCurrent doesn't exist for some reason.
	ErrNoCurrentStage = errors.New("kahinah: couldn't find the current stage (yikes!)")
)

//
// --------------------------------------------------------------------
// LISTS
// --------------------------------------------------------------------
//

// ListsAPIHandler shows the collection of lists, with filters, paginated, JSON'ified.
func ListsAPIHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)
	dataRenderer.Type = data.DataJSON

	platform := r.FormValue("platform")
	channel := r.FormValue("channel")
	status := r.FormValue("status")

	limit := int(to.Int64(r.FormValue("limit")))
	if limit <= 0 {
		limit = 50 // reasonable
	}

	page := int(to.Int64(r.FormValue("page")))
	if page <= 0 {
		page = 1
	}

	baseDB := models.DB.Model(&models.List{})
	if platform != "" {
		baseDB = baseDB.Where("platform = ?", platform)
	}

	if channel != "" {
		baseDB = baseDB.Where("channel = ?", channel)
	}

	if status == models.ListRunning || status == models.ListPending || status == models.ListSuccess || status == models.ListFailed {
		baseDB = baseDB.Where("stage_result = ?", status)
	} else if status != "" {
		panic(ErrBadRequest)
	}

	var cnt int
	if err := baseDB.Count(&cnt).Error; err != nil {
		panic(err)
	}

	totalpages := cnt / limit
	if cnt%limit != 0 {
		totalpages++
	}

	if page > totalpages {
		page = totalpages
	}

	var packages []*models.List
	if err := baseDB.Limit(limit).Offset((page - 1) * limit).Order("updated_at desc").Find(&packages).Error; err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	dataRenderer.Data = map[string]interface{}{
		"lists": packages,
		"pages": map[string]interface{}{
			"prev":    page - 1,
			"current": page,
			"next":    page + 1,
			"total":   totalpages,
		},
	}
}

// ListsHandler shows the collection of lists (HTML).
func ListsHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)

	// filters that will be passed on by Vue.js to the API
	platform := r.FormValue("platform")
	channel := r.FormValue("channel")
	status := r.FormValue("status")
	limit := r.FormValue("limit")
	page := r.FormValue("page")

	dataRenderer.Data = map[string]interface{}{
		"Title":    "Lists",
		"Nav":      2,
		"Platform": platform,
		"Channel":  channel,
		"Status":   status,
		"Limit":    limit,
		"Page":     page,
	}
	dataRenderer.Template = "i/list"
}

//
// --------------------------------------------------------------------
// INDIVIDUAL BUILD
// --------------------------------------------------------------------
//

// BuildGetHandler displays build information for a specific build
func BuildGetHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)

	toRender := map[string]interface{}{}

	id := to.Uint64(pat.Param(ctx, "id"))

	var pkg models.List
	if err := models.DB.Where("id = ?", id).First(&pkg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(ErrNotFound)
		} else {
			panic(err)
		}
	}

	toRender["Title"] = "Build " + to.String(id) + ": " + pkg.Name
	toRender["Nav"] = 2
	toRender["ID"] = to.String(id)

	dataRenderer.Data = toRender
	dataRenderer.Template = "builds/build"
}

type buildGetJSONStage struct {
	Name     string
	Status   map[string]processes.ProcessStatus
	Metadata map[string]interface{}
	Optional map[string]bool
}

type buildGetJSON struct {
	ID       uint
	Platform string
	Channel  string   // maps to repo
	Variants []string // maps to architectures
	Name     string

	Artifacts []models.ListArtifact
	Links     []models.ListLink
	Activity  []models.ListActivity
	Changes   string

	BuildDate time.Time
	Updated   time.Time

	PlatformConfig string
	Stages         []buildGetJSONStage
	CurrentStage   string
	Status         string
	Advisory       uint
}

// BuildGetJSONHandler displays build information in JSON for a specific build.
func BuildGetJSONHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	dataRenderer := data.FromContext(ctx)
	id := to.Uint64(pat.Param(ctx, "id"))
	user := models.FindUserNoCreate(Authenticated(r))

	// load the requested build list
	var pkg models.List
	if err := models.DB.Where("id = ?", id).First(&pkg).Related(&pkg.Activity, "Activity").Related(&pkg.Artifacts, "Artifacts").Related(&pkg.Links, "Links").Related(&pkg.Stages, "Stages").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(ErrNotFound)
		} else {
			panic(err)
		}
	}

	// load stage information
	var stageInfo []buildGetJSONStage

	for _, v := range pkg.Stages {
		if err := models.DB.Related(&v.Processes, "Processes").Error; err != nil && err != gorm.ErrRecordNotFound {
			panic(err)
		}

		// get all process info
		status := map[string]processes.ProcessStatus{}
		metadata := map[string]interface{}{}
		optional := map[string]bool{}

		for _, p := range v.Processes {
			process, err := processes.BuildProcess(&p)
			if err != nil {
				panic(err)
			}

			status[p.Name] = process.Status()
			metadata[p.Name] = process.APIMetadata(user)
			optional[p.Name] = p.Optional
		}

		stageInfo = append(stageInfo, buildGetJSONStage{
			Name:     v.Name,
			Status:   status,
			Metadata: metadata,
			Optional: optional,
		})
	}

	// render the data in a nice way

	dataRenderer.Data = &buildGetJSON{
		ID:       pkg.ID,
		Platform: pkg.Platform,
		Channel:  pkg.Channel,
		Variants: strings.Split(pkg.Variants, ";"),
		Name:     pkg.Name,

		Artifacts: pkg.Artifacts,
		Links:     pkg.Links,
		Activity:  pkg.Activity,
		Changes:   pkg.Changes,

		BuildDate: pkg.BuildDate,
		Updated:   pkg.UpdatedAt,

		PlatformConfig: pkg.PlatformGitConfig,
		Stages:         stageInfo,
		CurrentStage:   pkg.StageCurrent,
		Status:         pkg.StageResult,
		Advisory:       pkg.AdvisoryID,
	}
	dataRenderer.Type = data.DataJSON
}

// BuildPostHandler handles post actions that occur to the current active stage.
func BuildPostHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	// check for authentication
	user := models.FindUser(MustAuthenticate(r))

	// setup
	dataRenderer := data.FromContext(ctx)

	// read parameters
	id := to.Uint64(pat.Param(ctx, "id"))
	target := r.FormValue("target") // either activity or process
	name := r.FormValue("name")     // activity (ignored), process - find process
	action := r.FormValue("action") // activity (ignored), process passed on
	value := r.FormValue("value")   // activity (comment), process passed on

	// find the build list
	var pkg models.List
	if err := models.DB.Where("id = ?", id).First(&pkg).Related(&pkg.Stages, "Stages").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(ErrNotFound)
		} else {
			panic(err)
		}
	}

	var result interface{}

	// act based on target
	switch target {
	case "activity":
		if value == "" {
			panic(ErrBadRequest)
		}
		pkg.AddActivity(user, value)
		result = map[string]interface{}{
			"success": true,
		}
	case "process":
		// load up the stage & process
		if err := models.DB.Related(&pkg.Stages, "Stages").Error; err != nil {
			panic(err)
		}
		var currentStage *models.ListStage
		for _, v := range pkg.Stages {
			if v.Name == pkg.StageCurrent {
				currentStage = &v
				break
			}
		}

		if currentStage == nil {
			panic(ErrNoCurrentStage)
		}

		if err := models.DB.Related(&currentStage.Processes).Error; err != nil {
			panic(err)
		}

		var selectedProcess *models.ListStageProcess
		for _, v := range currentStage.Processes {
			if v.Name == name {
				selectedProcess = &v
				break
			}
		}

		if selectedProcess == nil {
			panic(ErrBadRequest)
		}

		// initialise the process
		process, err := processes.BuildProcess(selectedProcess)
		if err != nil {
			panic(err)
		}

		r, err := process.APIRequest(user, action, value)
		if err != nil {
			result = map[string]interface{}{
				"error":   true,
				"message": err.Error(),
				"result":  r,
			}
		} else {
			result = r
		}
	default:
		panic(ErrBadRequest)
	}

	dataRenderer.Type = data.DataJSON
	dataRenderer.Data = result

	//http.Redirect(rw, r, r.URL.String(), http.StatusTemporaryRedirect)
}
