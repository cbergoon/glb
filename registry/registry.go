package registry

import (
	"errors"
)

var (
	ErrServiceNotFound       = errors.New("registry: target name/version not found")
	ErrServiceNameNotAllowed = errors.New("registry: service name not allowed; non-allowable service names [reload|status]")
)

type Registry interface {
	Add(svc string, key string, t Target)           //Adds an entry to registry.
	Delete(svc string, key string, t Target)        //Removes an entry from the registry.
	Lookup(svc, key string) (OrderedTargets, error) //Retrieves a slice of addresses for specified service/version.
	Validate() error                                //Ensures no services contain a reserved word.
	IncrementFailures(svc string, key string, t Target, amount int) (int, error)
	SetRoundRobbinCounter(svc string, key string, value int) (int, error)
	GetRoundRobbinCounter(svc string, key string) (int, error)
}
