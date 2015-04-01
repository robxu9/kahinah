package apiv1

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/robxu9/kahinah/server/common"
	"github.com/unrolled/render"
)

const (
	LIST_LIMIT = 50 // only do 50 per page
)

type APIv1 struct {
	c *common.Common
	r *render.Render
}

func New(c *common.Common) http.Handler {
	api := &APIv1{
		c: c,
		r: render.New(render.Options{}),
	}

	r := mux.NewRouter()

	// Routes:
	// /auth: authentication functions
	// /updates: list of updates
	r.HandleFunc("/updates/targets", c.UserWrapHandler(api.UpdateTargets))
	r.HandleFunc("/updates", c.UserWrapHandler(api.Updates))
	// /advisories: list of advisories

	// /: return apiv1 server version
	r.HandleFunc("/", c.UserWrapHandler(api.Root))

	return r
}

/**
 * @apiDefine Lists
 * @apiSuccess {Object} pages				Page List
 * @apiSuccess {String} pages.next	Next Page of List
 * @apiSuccess {String} pages.prev	Previous Page of List
 * @apiSuccess {String} pages.head	First Page of List
 * @apiSuccess {String} pages.tail	Last Page of List
 * @apiSuccess {Number} pages.total	Total Number of Pages of List
 */
// makeLists takes in the original url, the current page, and the total
// number of entries, then using LIST_LIMIT, it returns the links map
func (a *APIv1) makeLists(originalURL *url.URL, page, totalEntries int) map[string]string {
	// get limit from LIST_LIMIT
	result := map[string]string{}

	// clone url object
	copyURL, err := url.Parse(originalURL.String())
	if err != nil {
		panic(err)
	}

	// next
	if (page+1)*LIST_LIMIT <= totalEntries {
		copyURL.Query().Set("page", strconv.Itoa(page+1))
		result["next"] = copyURL.String()
	}

	// prev
	if (page-1) >= 1 && (page-1)*LIST_LIMIT <= totalEntries {
		copyURL.Query().Set("page", strconv.Itoa(page-1))
		result["prev"] = copyURL.String()
	}

	// head
	copyURL.Query().Set("page", "1")
	result["head"] = copyURL.String()

	// tail
	copyURL.Query().Set("page", strconv.Itoa(totalEntries/LIST_LIMIT+1))
	result["tail"] = copyURL.String()

	// total # of pages
	result["total"] = strconv.Itoa(totalEntries/LIST_LIMIT + 1)

	return result
}

/**
 * @apiDefine Authentication
 * @apiParam {String} token API Token via Query (should only be used for server to server)
 * @apiHeader {String} Authentication Authentication Token (Should be Bearer <token>)
 * @apiError Unauthorized The authentication token failed.
 */

/**
 * @api {get} / Kahinah Server Version
 * @apiName Version
 * @apiDescription Retrieves Kahinah's server version
 * @apiGroup Internals
 * @apiSuccess {Number} version Version of the Server
 */
func (a *APIv1) Root(rw http.ResponseWriter, r *http.Request, t *common.UserToken) {
	a.r.JSON(rw, http.StatusOK, map[string]interface{}{
		"version": a.c.C.Version,
	})
}
