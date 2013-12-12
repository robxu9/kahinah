package integration

import (
	"log"
)

type Integration interface {
	Ping() error
	PingParams(map[string]string) error
	Url(string) string
	Publish(string) error
	Reject(string) error
}

var integrated map[string]Integration

func init() {
	integrated = make(map[string]Integration)
}

func Integrate(k string, i Integration) {
	integrated[k] = i
}

func Ping() {
	for k, v := range integrated {
		go func() {
			err := v.Ping()
			if err != nil {
				log.Printf("Error pinging integrator %s: %s\n", k, err)
			}
		}()
	}
}

func PingParams(m map[string]string) {
	for k, v := range integrated {
		go func() {
			err := v.PingParams(m)
			if err != nil {
				log.Printf("Error pinging integrator %s with parameters: %s\n", k, err)
			}
		}()
	}
}

func Url(handle, id string) string {
	return integrated[handle].Url(id)
}

func Publish(handle, id string) {
	err := integrated[handle].Publish(id)
	if err != nil {
		log.Printf("Error publishing integrator %s with parameters %s: %s\n", handle, id, err)
	}
}

func Reject(handle, id string) {
	err := integrated[handle].Reject(id)
	if err != nil {
		log.Printf("Error rejecting integrator %s with parameters %s: %s\n", handle, id, err)
	}
}
