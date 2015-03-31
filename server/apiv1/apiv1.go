package apiv1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/robxu9/kahinah/server/common"
	"github.com/unrolled/render"
)

type APIv1 struct {
	c *common.Common
	r *render.Render
}

func New(c *common.Common) http.Handler {
	api := &APIv1{
		c: c,
		r: render.New(render.Options{
			IndentJSON: true,
			IndentXML:  true,
		}),
	}

	r := mux.NewRouter()

	// Routes:
	// /auth: authentication functions
	// /updates: list of updates
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
 */

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
 * @apiGroup Server Internals
 * @apiSuccess {Number} version Version of the Server
 */
func (a *APIv1) Root(rw http.ResponseWriter, r *http.Request, t *common.UserToken) {
	a.r.JSON(rw, http.StatusOK, map[string]interface{}{
		"version": a.c.C.Version,
	})
}
