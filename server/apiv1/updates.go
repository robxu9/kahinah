package apiv1

import (
	"net/http"
	"strconv"

	"github.com/robxu9/kahinah/kahinah"
	"github.com/robxu9/kahinah/server/common"
)

/**
 * @api {get} /updates
 * @apiName ListUpdates
 * @apiDescription List the latest updates coming into Kahinah
 * @apiGroup Updates
 * @apiParam {Number} page=1 Pagination - starts at one
 * @apiParam {String="none","bugfix","security","enhancement","new"} type Filter Update Type
 * @apiParam {String} for Filter Update Target
 * @apiUse Lists
 * @apiSuccess {Object[]} updates     List of Updates
 * @apiSuccess {Number} updates.id    Update ID
 */
func (a *APIv1) Updates(rw http.ResponseWriter, r *http.Request, t *common.UserToken) {
	filter_page := 1
	filter_type := r.URL.Query().Get("type")
	filter_for := r.URL.Query().Get("for")

	if p := r.URL.Query().Get("page"); p != "" {
		if i, err := strconv.Atoi(p); err == nil {
			if i >= 1 {
				filter_page = i
			}
		}
	}

	updateIds, err := a.c.K.ListUpdates(int64(LIST_LIMIT*(filter_page-1)), int64(LIST_LIMIT))
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

	links := a.makeLists(r.URL, filter_page-1, int(a.c.K.CountUpdates()))

	a.r.JSON(rw, http.StatusOK, map[string]interface{}{
		"updates": updates,
		"links":   links,
	})
}
