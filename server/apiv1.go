package main

import "net/http"

type APIv1Endpoint struct {
	Endpoints []Endpoint
}

func (a *APIv1Endpoint) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

}

func (a *APIv1Endpoint) Handle(rw http.ResponseWriter, r *http.Request, e Endpoint) {

}

type Endpoint interface {
	// Find corresponds to GET /endpoint/:id
	Find(id int64, rw http.ResponseWriter, r *http.Request)
	// FindAll corresponds to GET /endpoint
	FindAll(rw http.ResponseWriter, r *http.Request)
	// Update corresponds to PUT /endpoint/:id
	Update(id int64, rw http.ResponseWriter, r *http.Request)
	// Create corresponds to POST /endpoint
	Create(rw http.ResponseWriter, r *http.Request)
	// Delete corresponds to DELETE /endpoint/:id
	Delete(id int64, rw http.ResponseWriter, r *http.Request)
}

type AdvisoryEndpoint struct {
}

func (a *AdvisoryEndpoint) Find(id int64, rw http.ResponseWriter, r *http.Request) {

}
