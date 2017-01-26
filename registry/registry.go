package registry

import (
	"errors"
)

var (
	ErrServiceNotFound       = errors.New("registry: target name/key not found")
	ErrServiceNameNotAllowed = errors.New("registry: service name not allowed; non-allowable service names [reload|status]")
)

type Registry interface {
	Add(svcValue string, keyValue string, t Target)                                        //Adds an entry to registry.
	Delete(svcValue string, keyValue string, t Target)                                     //Removes an entry from the registry.
	Lookup(svcValue string, keyValue string) (OrderedTargets, error)                       //Retrieves a slice of addresses for specified service/version.
	IncrementFailures(svcValue string, keyValue string, t Target, amount int) (int, error) //Increments failures counter on target.
	SetRoundRobbinCounter(svcValue string, keyValue string, value int) (int, error)        //Sets round robbin counter on key.
	GetRoundRobbinCounter(svcValue string, keyValue string) (int, error)                   //Gets round robbin counter on key.
}
