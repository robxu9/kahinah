package apiv1

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/robxu9/kahinah/kahinah"
	"github.com/robxu9/kahinah/server/common"
)

// Updates lists the latest updates coming into kahinah.
/**
 * @api {get} /updates List Updates
 * @apiName ListUpdates
 * @apiDescription List the latest updates coming into Kahinah
 * @apiGroup Updates
 * @apiParam {Number} page=1 Pagination - starts at one
 * @apiParam {String="none","bugfix","security","enhancement","new"} type Filter Update Type
 * @apiParam {String} for Filter Update Target
 * @apiUse Lists
 * @apiSuccess (200) {Object[]} updates     List of Updates
 * @apiSuccess (200) {Number} updates.id    Update ID
 */
func (a *APIv1) Updates(rw http.ResponseWriter, r *http.Request, t *common.UserToken) {
	filterPage := 1
	filterType := r.URL.Query().Get("type")
	filterFor := r.URL.Query().Get("for")

	if p := r.URL.Query().Get("page"); p != "" {
		if i, err := strconv.Atoi(p); err == nil {
			if i >= 1 {
				filterPage = i
			}
		}
	}

	actualType := kahinah.NONE

	switch filterType {
	case "bugfix":
		actualType = kahinah.BUGFIX
	case "security":
		actualType = kahinah.SECURITY
	case "enhancement":
		actualType = kahinah.ENHANCEMENT
	case "new":
		actualType = kahinah.NEW
	}

	updateIds, err := a.c.K.ListUpdates(int64(LIST_LIMIT*(filterPage-1)), int64(LIST_LIMIT), actualType, filterFor)
	if err != nil {
		panic(err)
	}

	updates := []*kahinah.Update{}
	for _, v := range updateIds {
		update, err := a.c.K.RetrieveUpdate(v)
		if err != nil {
			panic(err)
		}
		updates = append(updates, update)
	}

	links := a.makeLists(r.URL, filterPage-1, int(a.c.K.CountUpdates()))

	a.r.JSON(rw, http.StatusOK, map[string]interface{}{
		"updates": updates,
		"links":   links,
	})
}

// UpdateTargets lists the distinct update targets for updates in Kahinah
/**
 * @api {get} /updates/targets List Update Targets
 * @apiName ListUpdateTargets
 * @apiDescription List distinct update targets for updates in Kahinah
 * @apiGroup Updates
 * @apiSuccess (200) {String[]} targets     List of Targets
 */
func (a *APIv1) UpdateTargets(rw http.ResponseWriter, r *http.Request, t *common.UserToken) {
	targets, err := a.c.K.ListUpdateTargets()
	if err != nil {
		panic(err)
	}

	a.r.JSON(rw, http.StatusOK, map[string]interface{}{
		"targets": targets,
	})
}

// Update displays a specific update with the specified ID.
/**
 * @api {get} /updates/:id Get Update
 * @apiName GetUpdate
 * @apiDescription Retrieve a specific update
 * @apiGroup Updates
 * @apiSuccess (200) {Object} update Update
 * @apiError (404) UpdateNotFound No update with the id specified was found.
 */
func (a *APIv1) Update(rw http.ResponseWriter, r *http.Request, t *common.UserToken) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		panic(err)
	}

	update, err := a.c.K.RetrieveUpdate(id)
	if err != nil {
		a.makeError(rw, http.StatusNotFound)
		return
	}

	a.r.JSON(rw, http.StatusOK, map[string]interface{}{
		"update": update,
	})
}
