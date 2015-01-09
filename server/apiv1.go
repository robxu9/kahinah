package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// The standard endpoints:
// find ~> GET /$obj/{id:[0-9]+}
// findall ~> GET /$obj
// update ~> PUT /$obj/{id:[0-9]+}
// create ~> POST /$obj
// delete ~> DELETE /updates/{id:[0-9]+}
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

func NewAPIv1Endpoint() *mux.Router {
	r := mux.NewRouter()

	// version
	r.HandleFunc("/", v1HomeHandler)

	// updates

	// find ~> GET /updates/{id:[0-9]+}
	v1RegisterStandardHandlers(r.Path("/updates/").Subrouter(), &UpdateStdHandlers{})

	// findall ~> GET /updates

	// update ~> PUT /updates/{id:[0-9]+}

	// create ~> POST /updates

	// delete ~> DELETE /updates/{id:[0-9]+}

	// advisories

	// users

	return r
}

func v1HomeHandler(rw http.ResponseWriter, r *http.Request) {
	rend.JSON(rw, http.StatusOK, map[string]interface{}{
		"name":    "kahinah",
		"api":     1,
		"version": VERSION,
	})
}

type UpdateStdHandlers struct {
}

func (u *UpdateStdHandlers) Find(rw http.ResponseWriter, r *http.Request) {

}

func (u *UpdateStdHandlers) FindAll(rw http.ResponseWriter, r *http.Request) {

}

func (u *UpdateStdHandlers) Update(rw http.ResponseWriter, r *http.Request) {

}

func (u *UpdateStdHandlers) Create(rw http.ResponseWriter, r *http.Request) {

}

func (u *UpdateStdHandlers) Delete(rw http.ResponseWriter, r *http.Request) {

}
