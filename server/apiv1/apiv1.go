package apiv1

import (
	"net/http"

	"github.com/gorilla/mux"
)

type APIv1 struct {
	k *kahinah.Kahinah
	c *common.Config
	r *mux.Router
}

func (a *APIv1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.r.ServeHTTP(w, r)
}

func NewAPIv1(k *kahinah.Kahinah, c *common.Config) *APIv1 {

}

// The standard endpoints:
// find ~> GET /$obj/{id:[0-9]+}
// findall ~> GET /$obj
// update ~> PUT /$obj/{id:[0-9]+}
// create ~> POST /$obj
// delete ~> DELETE /$obj/{id:[0-9]+}
//
// If one is not supported, return 405/http.StatusMethodNotAllowed.
type StandardHandlers interface {
	Find(http.ResponseWriter, *http.Request)
	FindAll(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	Create(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

func v1RegisterStandardHandlers(r *mux.Router, s StandardHandlers) {
	specificRouter := r.Path("/{id:[0-9]+}").Subrouter()
	specificRouter.Methods("GET").HandlerFunc(s.Find)
	specificRouter.Methods("PUT").HandlerFunc(s.Update)
	specificRouter.Methods("DELETE").HandlerFunc(s.Delete)

	r.Path("/").Methods("GET").HandlerFunc(s.FindAll)
	r.Path("/").Methods("POST").HandlerFunc(s.Create)
}

func NewAPIv1Endpoint(k *kahinah.Kahinah, config *common.Config) *mux.Router {
	rend := render.New(render.Options{
		IsDevelopment: config.DevMode,
	})

	r := mux.NewRouter()

	// ~> version
	r.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rend.JSON(rw, http.StatusOK, map[string]interface{}{
			"name":    "kahinah",
			"api":     1,
			"version": VERSION,
			})
		})

	// ~> updates
	// standard router
	v1RegisterStandardHandlers(r.Path("/updates/").Subrouter(), &v1UpdateHandlers{})

	// advisories

	// users

	return r
}

type v1UpdateHandlers struct{
	k *kahinah.Kahinah
	r *render.Render
	c *common.Config
}

func (v *v1UpdateHandlers) Find(rw http.ResponseWriter, r *http.Request) {

}

func (v *v1UpdateHandlers) FindAll(rw http.ResponseWriter, r *http.Request) {

}

func (v *v1UpdateHandlers) Update(rw http.ResponseWriter, r *http.Request) {

}

func (v *v1UpdateHandlers) Create(rw http.ResponseWriter, r *http.Request) {

}

func (v *v1UpdateHandlers) Delete(rw http.ResponseWriter, r *http.Request) {

}
