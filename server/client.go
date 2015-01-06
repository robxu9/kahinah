package main

import (
	"net/http"
	"os"
)

type ClientEndpoint struct {
}

func (c *ClientEndpoint) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	if path == "/" { // redirect to index.html
		path = "/index.html"
	}

	path = "public" + path

	// get info
	stat, err := os.Stat(path)
	if err != nil || stat.IsDir() {
		// then just serve index.html
		http.ServeFile(rw, req, "public/index.html")
		return
	}

	// if the file exists, serve it directly
	http.ServeFile(rw, req, path)
}
