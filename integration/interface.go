package integration

import (
	"errors"

	"github.com/robxu9/kahinah/models"
)

var (
	ErrDisabled       = errors.New("integration: disabled")
	ErrNotImplemented = errors.New("integration: not implemented")
)

type Integration interface {
	Poll() error            // poll for new updates
	Hook(interface{}) error // receive and process a webhook from the BS

	Accept(*models.BuildList) error // signal to the BS to accept
	Reject(*models.BuildList) error // signal to the BS to reject
}

var handler Integration

func Integrate(i Integration) {
	handler = i
}

func Poll() error {
	return handler.Poll()
}

func Hook(i interface{}) error {
	return handler.Hook(i)
}

func Accept(l *models.BuildList) error {
	return handler.Accept(l)
}

func Reject(l *models.BuildList) error {
	return handler.Reject(l)
}
