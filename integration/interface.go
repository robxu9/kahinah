package integration

import (
	"net/http"
	"sync"

	"github.com/robxu9/kahinah/models"
)

var (
	Implementations = map[string]Integration{}
	pollAllLock     = &sync.Mutex{}
)

type Integration interface {
	Poll() error          // poll for new updates
	Hook(r *http.Request) // receive and process a webhook from the BS

	Accept(*models.List) error // signal the system to accept
	Reject(*models.List) error // signal the system to reject
}

func PollAll() map[string]error {
	pollAllLock.Lock()
	defer pollAllLock.Unlock()

	implErrors := map[string]error{}

	for name, impl := range Implementations {
		if err := impl.Poll(); err != nil {
			implErrors[name] = err
		}
	}

	return implErrors
}
