package serviceregistry

import (
	"sync"
	"github.com/cbergoon/glb/registry"
)

var lock sync.RWMutex //Exclusive lock for registry data structure.

//Registry data structure
type ServiceRegistry map[string]map[string][]registry.Target

//Retrieves a slice of addresses for specified service/version. If no address is found
//ErrServiceNotFound is returned.
func (r ServiceRegistry) Lookup(name, version string) ([]registry.Target, error) {
	lock.RLock()
	targets, ok := r[name][version]
	lock.RUnlock()
	if !ok || len(targets) == 0 {
		return nil, registry.ErrServiceNotFound
	}
	return targets, nil
}

//Adds an entry to the registry. If the address for an entry exists it will be duplicated.
//If the service does not exist a new map[string][]string will be created and added to represent
//the new service. If necessary registry.Lookup can be used to ensure success.
func (r ServiceRegistry) Add(name, version, t registry.Target) {
	lock.Lock()
	defer lock.Unlock()

	service, ok := r[name]
	if !ok {
		service = map[string][]registry.Target{}
		r[name] = service
	}
	service[version] = append(service[version], t)
}

//Removes an address entry for a given service/version entry. If necessary the
//registry.Lookup function can be used to ensure success.
func (r ServiceRegistry) Delete(name, version, t registry.Target) {
	lock.Lock()
	defer lock.Unlock()

	service, ok := r[name]
	if !ok {
		return
	}
begin:
	for i, tgt := range service[version] {
		if tgt.Address == t.Address {
			copy(service[version][i:], service[version][i+1:])
			service[version][len(service[version])-1].Address = ""
			service[version][len(service[version])-1].RoundRobbinIndex = 0
			service[version][len(service[version])-1].Failures = 0
			service[version] = service[version][:len(service[version])-1]
			goto begin
		}
	}
}

//Ensures registry adheres to the restrictions set forth by the registry definition. Currently
//ensures glb keywords are not used in the service/version combination. If this constraint is
//violated ErrServiceNameNotAllowed is returned.
func (r ServiceRegistry) Validate() error {
	_, ok := r["reload"]
	if ok {
		return registry.ErrServiceNameNotAllowed
	}
	_, ok = r["status"]
	if ok {
		return registry.ErrServiceNameNotAllowed
	}
	return nil
}
