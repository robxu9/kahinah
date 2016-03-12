package controllers

import (
	"net/http"
	"strings"

	"github.com/robxu9/kahinah/data"
	"github.com/robxu9/kahinah/integration"
	"goji.io/pattern"
	"golang.org/x/net/context"
)

// IntegrationHandler handles webhooks to integration handlers in the integration
// package by calling their hook function with the request.
func IntegrationHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	// integration is push only with status code returns. no response output.
	dataRender := data.FromContext(ctx)
	dataRender.Type = data.DataNoRender

	// get the handler name (as the first part)
	handlerName := pattern.Path(ctx)
	impl := integration.Implementations[strings.Split(handlerName, "/")[0]]

	if impl == nil {
		http.NotFound(rw, r)
	}

	impl.Hook(r)

	rw.WriteHeader(200)
}
