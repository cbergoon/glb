package registry

import (
	"errors"
	"sync"
)

var lock sync.RWMutex //Exclusive lock for registry data structure.

var (
	ErrServiceNotFound       = errors.New("registry: target name/version not found")
	ErrServiceNameNotAllowed = errors.New("registry: service name not allowed; non-allowable service names [reload|status]")
)

type Registry interface {
	Add(name, version, endpoint string) //Adds an entry to registry.
	Delete(name, version, endpoint string) //Removes an entry from the registry.
	Lookup(name, version string) ([]string, error) //Retrieves a slice of addresses for specified service/version.
	Validate() error //Ensures no services contain a reserved word.
}

//Registry data structure
type DefaultRegistry map[string]map[string][]string

//Retrieves a slice of addresses for specified service/version. If no address is found
//ErrServiceNotFound is returned.
func (r DefaultRegistry) Lookup(name, version string) ([]string, error) {
	lock.RLock()
	targets, ok := r[name][version]
	lock.RUnlock()
	if !ok || len(targets) == 0 {
		return nil, ErrServiceNotFound
	}
	return targets, nil
}

//Adds an entry to the registry. If the address for an entry exists it will be duplicated.
//If the service does not exist a new map[string][]string will be created and added to represent
//the new service. If necessary registry.Lookup can be used to ensure success.
func (r DefaultRegistry) Add(name, version, endpoint string) {
	lock.Lock()
	defer lock.Unlock()

	service, ok := r[name]
	if !ok {
		service = map[string][]string{}
		r[name] = service
	}
	service[version] = append(service[version], endpoint)
}

//Removes an address entry for a given service/version entry. If necessary the
//registry.Lookup function can be used to ensure success.
func (r DefaultRegistry) Delete(name, version, endpoint string) {
	lock.Lock()
	defer lock.Unlock()

	service, ok := r[name]
	if !ok {
		return
	}
begin:
	for i, svc := range service[version] {
		if svc == endpoint {
			copy(service[version][i:], service[version][i+1:])
			service[version][len(service[version])-1] = ""
			service[version] = service[version][:len(service[version])-1]
			goto begin
		}
	}
}

//Ensures registry adheres to the restrictions set forth by the registry definition. Currently
//ensures glb keywords are not used in the service/version combination. If this constraint is
//violated ErrServiceNameNotAllowed is returned.
func (r DefaultRegistry) Validate() error {
	_, ok := r["reload"]
	if ok {
		return ErrServiceNameNotAllowed
	}
	_, ok = r["status"]
	if ok {
		return ErrServiceNameNotAllowed
	}
	return nil
}
