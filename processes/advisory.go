package processes

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"menteslibres.net/gosexy/to"

	"github.com/robxu9/kahinah/models"
)

func newAdvisoryProcess(m *models.ListStageProcess, cfg map[string]interface{}) Process {
	a := &Advisory{
		model: m,
	}

	// unmarshal data
	if m.Data != nil {
		json.Unmarshal(m.Data, a)
	}

	// and Dialect
	a.Dialect = cfg["dialect"].(string)

	return a
}

func init() {
	Mapping["advisory"] = newAdvisoryProcess
}

var (
	ErrAdvisoryAlreadyReady = errors.New("process/advisory: already ready for finalising")
	ErrAdvisoryNotReady     = errors.New("process/advisory: not ready for finalising")
	ErrAdvisoryNoAction     = errors.New("process/advisory: no action/value specified")
)

type Advisory struct {
	model *models.ListStageProcess

	Dialect string `json:"-"`

	Type        string
	Summary     string
	Description string

	ExistingAdvisoryID uint

	Ready bool // whether we're good on this
}

func (a *Advisory) Start() error {
	// we don't exactly need to start anything
	return nil
}

func (a *Advisory) Stop() error {
	// we don't exactly need to stop anything
	return nil
}

func (a *Advisory) Reset() error {
	a.Type = ""
	a.Summary = ""
	a.Description = ""
	a.ExistingAdvisoryID = 0
	a.Ready = false

	return a.save()
}

func (a *Advisory) save() error {
	bte, err := json.Marshal(a)
	if err != nil {
		return err
	}
	a.model.Data = bte
	return nil
}

func (a *Advisory) Finalise() error {
	if !a.Ready {
		return ErrAdvisoryNotReady
	}

	// -- create the advisory
	list := a.model.ParentList()

	if a.ExistingAdvisoryID == 0 {
		// and now create it
		adv := &models.Advisory{
			Dialect:     a.Dialect,
			Year:        strconv.Itoa(time.Now().Year()),
			Type:        a.Type,
			Summary:     a.Summary,
			Description: a.Description,
			Lists: []models.List{
				*list,
			},
		}

		if err := models.DB.Create(adv).Error; err != nil {
			return err
		}
	} else { // attach it
		if err := models.DB.Model(&models.Advisory{}).Association("Lists").Append(list).Error; err != nil {
			return err
		}
	}

	return a.save()
}

func (a *Advisory) Status() ProcessStatus {
	if !a.Ready {
		return ProcessRunning
	}
	return ProcessOK
}

func (a *Advisory) APIRequest(user *models.User, action string, value string) (interface{}, error) {
	if a.Ready {
		return nil, ErrAdvisoryAlreadyReady
	}

	switch action {
	case "type":
		a.Type = value
		a.model.ParentList().AddActivity(user, fmt.Sprintf("[Advisory] switched type to %v", value))
	case "summary":
		a.Summary = value
		a.model.ParentList().AddActivity(user, fmt.Sprintf("[Advisory] switched summary to %v", value))
	case "description":
		a.Description = value
		a.model.ParentList().AddActivity(user, fmt.Sprintf("[Advisory] switched description to %v", value))

	case "existingAdv":
		a.ExistingAdvisoryID = uint(to.Uint64(value))
		if a.ExistingAdvisoryID == 0 {
			a.model.ParentList().AddActivity(user, "[Advisory] detached from existing advisory")
		} else {
			a.model.ParentList().AddActivity(user, fmt.Sprintf("[Advisory] attached to existing advisory %v", a.ExistingAdvisoryID))
		}
	case "ready":
		if a.ExistingAdvisoryID != 0 {
			a.Ready = true
		} else {
			if a.Type == "" || a.Summary == "" || a.Description == "" {
				return nil, ErrAdvisoryNotReady
			}
		}
		a.model.ParentList().AddActivity(user, "[Advisory] marked ready")
	default:
		return nil, ErrAdvisoryNoAction
	}

	if err := a.save(); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
	}, nil
}

func (a *Advisory) APIMetadata(user *models.User) interface{} {
	return map[string]interface{}{
		"dialect":     a.Dialect,
		"type":        a.Type,
		"summary":     a.Summary,
		"description": a.Description,
		"existingid":  a.ExistingAdvisoryID,
	}
}
