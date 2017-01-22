package serviceregistry

import (
	"github.com/cbergoon/glb/registry"
	"sync"
)

type StandardRegistry struct {
	lock     sync.RWMutex //Exclusive lock for registry data structure.
	Services map[string]*service
}

type service struct {
	Value string
	Keys  map[string]*key
}

type key struct {
	Value              string
	RoundRobbinCounter int
	Targets            registry.OrderedTargets
}

//Retrieves a slice of targets for specified service/key. If no address is found
//ErrServiceNotFound is returned.
func (r *StandardRegistry) Lookup(svcValue string, keyValue string) (registry.OrderedTargets, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if svcValue == "reload" || svcValue == "status" {
		return nil, registry.ErrServiceNameNotAllowed
	}
	if keyValue == "reload" || keyValue == "status" {
		return nil, registry.ErrServiceNameNotAllowed
	}
	s, ok := r.Services[svcValue]
	if !ok {
		return nil, registry.ErrServiceNotFound
	}
	k, ok := s.Keys[keyValue]
	if !ok {
		return nil, registry.ErrServiceNotFound
	}
	return k.Targets, nil
}

//Adds an entry to the registry. If the address for an entry exists it will be duplicated.
//If the service does not exist a new key and target will be created and added to represent
//the new service. If necessary registry.Lookup can be used to ensure success.
func (r *StandardRegistry) Add(svcValue string, keyValue string, t registry.Target) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if svcValue == "reload" || svcValue == "status" {
		return
	}
	if keyValue == "reload" || keyValue == "status" {
		return
	}
	if r.Services == nil {
		r.Services = make(map[string]*service)
	}
	_, ok := r.Services[svcValue]
	if !ok {
		r.Services[svcValue] = &service{Value: svcValue}
	}
	if r.Services[svcValue].Keys == nil {
		r.Services[svcValue].Keys = make(map[string]*key)
	}
	_, ok = r.Services[svcValue].Keys[keyValue]
	if !ok {
		r.Services[svcValue].Keys[keyValue] = &key{Value: keyValue, RoundRobbinCounter: 0}
		//r.Services[svc].Keys[key].Targets = append(registry.OrderedTargets)
	}
	r.Services[svcValue].Keys[keyValue].Targets = append(r.Services[svcValue].Keys[keyValue].Targets, t)
}

//Removes an address entry for a given service/key entry. If necessary the
//registry.Lookup function can be used to ensure success.
func (r *StandardRegistry) Delete(svcValue string, keyValue string, t registry.Target) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.Services == nil {
		return
	}
	_, ok := r.Services[svcValue]
	if !ok {
		return
	}
	if r.Services[svcValue].Keys == nil {
		return
	}
	_, ok = r.Services[svcValue].Keys[keyValue]
	if !ok {
		return
	}
	if r.Services[svcValue].Keys[keyValue].Targets == nil {
		return
	}
	targetIndex := indexOf(len(r.Services[svcValue].Keys[keyValue].Targets), func(i int) bool { return r.Services[svcValue].Keys[keyValue].Targets[i].Address == t.Address })
	if targetIndex < 0 {
		return
	}
	r.Services[svcValue].Keys[keyValue].Targets = append(r.Services[svcValue].Keys[keyValue].Targets[:targetIndex], r.Services[svcValue].Keys[keyValue].Targets[targetIndex+1:]...)
}
//Increment failure counter on target. If target is not found ErrServiceNotFound is returned.
func (r *StandardRegistry) IncrementFailures(svcValue string, keyValue string, t registry.Target, amount int) (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	_, ok := r.Services[svcValue]
	if !ok {
		return 0, registry.ErrServiceNotFound
	}
	_, ok = r.Services[svcValue].Keys[keyValue]
	if !ok {
		return 0, registry.ErrServiceNotFound
	}
	targetIndex := indexOf(len(r.Services[svcValue].Keys[keyValue].Targets), func(i int) bool { return r.Services[svcValue].Keys[keyValue].Targets[i].Address == t.Address })
	if targetIndex < 0 {
		return 0, registry.ErrServiceNotFound
	}
	r.Services[svcValue].Keys[keyValue].Targets[targetIndex].Failures += amount
	return r.Services[svcValue].Keys[keyValue].Targets[targetIndex].Failures, nil
}

//Set round robbin counter on key. If key is not found ErrServiceNotFound is returned.
func (r *StandardRegistry) SetRoundRobbinCounter(svcValue string, keyValue string, value int) (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	_, ok := r.Services[svcValue]
	if !ok {
		return 0, registry.ErrServiceNotFound
	}
	_, ok = r.Services[svcValue].Keys[keyValue]
	if !ok {
		return 0, registry.ErrServiceNotFound
	}
	r.Services[svcValue].Keys[keyValue].RoundRobbinCounter = value
	return r.Services[svcValue].Keys[keyValue].RoundRobbinCounter, nil
}

//Get round robbin counter on key. If key is not found ErrServiceNotFound is returned.
func (r *StandardRegistry) GetRoundRobbinCounter(svcValue string, keyValue string) (int, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	_, ok := r.Services[svcValue]
	if !ok {
		return 0, registry.ErrServiceNotFound
	}
	_, ok = r.Services[svcValue].Keys[keyValue]
	if !ok {
		return 0, registry.ErrServiceNotFound
	}
	return r.Services[svcValue].Keys[keyValue].RoundRobbinCounter, nil
}

func indexOf(length int, predicate func(i int) bool) int {
	for i := 0; i < length; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
