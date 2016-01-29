package integration

import (
	"log"

	"github.com/robxu9/kahinah/models"
)

type Integration interface {
	// Check for new packages to add to the buildsystem
	Ping() error
	// Check for new packages to add to the buildsystem, with parameters
	PingParams(map[string]string) error
	// Retrieve the url of the buildlist commit page
	Commits(*models.BuildList) string
	// Retrieve the url of the corresponding buildlist
	Url(*models.BuildList) string
	// Publish the buildlist
	Publish(*models.BuildList) error
	// Reject the buildlist
	Reject(*models.BuildList) error
}

var handler Integration

func Integrate(i Integration) {
	handler = i
}

func Ping() {
	go func() {
		err := handler.Ping()
		if err != nil {
			log.Printf("Error pinging integrator: %s\n", err)
		}
	}()
}

func PingParams(m map[string]string) {
	go func() {
		err := handler.PingParams(m)
		if err != nil {
			log.Printf("Error pinging integrator with parameters: %s\n", err)
		}
	}()
}

func Commits(m *models.BuildList) string {
	return handler.Commits(m)
}

func Url(m *models.BuildList) string {
	return handler.Url(m)
}

func Publish(m *models.BuildList) {
	err := handler.Publish(m)
	if err != nil {
		log.Printf("Error calling publishing integrator for id %s: %s\n", m.HandleId, err)
	}
}

func Reject(m *models.BuildList) {
	err := handler.Reject(m)
	if err != nil {
		log.Printf("Error calling rejecting integrator for id %s: %s\n", m.HandleId, err)
	}
}
