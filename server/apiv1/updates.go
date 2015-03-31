package apiv1

import (
	"net/http"

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

}
