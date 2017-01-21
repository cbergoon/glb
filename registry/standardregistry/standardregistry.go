package serviceregistry

import (
	"github.com/cbergoon/glb/registry"
	"sort"
	"strings"
	"sync"
)

//var lock sync.RWMutex

//Registry data structure
//type ServiceRegistry map[string]map[string][]registry.Target

type orderedServices []service

func (s orderedServices) Len() int {
	return len(s)
}

func (s orderedServices) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s orderedServices) Less(i, j int) bool {
	less := strings.Compare(s[i].Value, s[j].Value)
	if less < 0 {
		return true
	} else {
		return false
	}
}

type orderedKeys []key

func (s orderedKeys) Len() int {
	return len(s)
}

func (s orderedKeys) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s orderedKeys) Less(i, j int) bool {
	less := strings.Compare(s[i].Value, s[j].Value)
	if less < 0 {
		return true
	} else {
		return false
	}
}

type StandardRegistry struct {
	lock     sync.RWMutex //Exclusive lock for registry data structure.
	Services orderedServices
}

type service struct {
	Value string
	Keys  orderedKeys
}

type key struct {
	Value              string
	RoundRobbinCounter int
	Targets            registry.OrderedTargets
}

func indexOf(length int, predicate func(i int) bool) int {
	for i := 0; i < length; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

//Retrieves a slice of addresses for specified service/version. If no address is found
//ErrServiceNotFound is returned.
func (r *StandardRegistry) Lookup(svc, key string) (registry.OrderedTargets, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	serviceIndex := indexOf(len(r.Services), func(i int) bool { return r.Services[i].Value == svc })
	if serviceIndex == -1 {
		return nil, registry.ErrServiceNotFound
	}
	keyIndex := indexOf(len(r.Services[serviceIndex].Keys), func(i int) bool { return r.Services[serviceIndex].Keys[i].Value == key })
	if keyIndex == -1 {
		return nil, registry.ErrServiceNotFound
	}
	return r.Services[serviceIndex].Keys[keyIndex].Targets, nil
}

//Adds an entry to the registry. If the address for an entry exists it will be duplicated.
//If the service does not exist a new map[string][]string will be created and added to represent
//the new service. If necessary registry.Lookup can be used to ensure success.
func (r *StandardRegistry) Add(svc string, key string, t registry.Target) {
	r.lock.Lock()
	defer r.lock.Unlock()
	serviceIndex := indexOf(len(r.Services), func(i int) bool { return r.Services[i].Value == svc })
	if serviceIndex == -1 {
		r.Services = append(r.Services, service{Value: svc})
		sort.Sort(r.Services)
	}
	keyIndex := indexOf(len(r.Services[serviceIndex].Keys), func(i int) bool { return r.Services[serviceIndex].Keys[i].Value == key })
	if keyIndex == -1 {
		r.Services[serviceIndex].Keys = append(r.Services[serviceIndex].Keys, key{Value: key})
		sort.Sort(r.Services[serviceIndex].Keys)
	}
	targetIndex := indexOf(len(r.Services[serviceIndex].Keys), func(i int) bool { return r.Services[serviceIndex].Keys[keyIndex].Targets[i].Address == t.Address })
	if targetIndex == -1 {
		r.Services[serviceIndex].Keys[keyIndex].Targets = append(r.Services[serviceIndex].Keys[keyIndex].Targets, t)
		sort.Sort(r.Services[serviceIndex].Keys[keyIndex].Targets)
	}
}

//Removes an address entry for a given service/version entry. If necessary the
//registry.Lookup function can be used to ensure success.
func (r *StandardRegistry) Delete(svc string, key string, t registry.Target) {
	r.lock.Lock()
	defer r.lock.Unlock()
	serviceIndex := indexOf(len(r.Services), func(i int) bool { return r.Services[i].Value == svc })
	if serviceIndex != -1 {
		r.Services = append(r.Services[:serviceIndex], r.Services[serviceIndex+1:]...)
	}
	keyIndex := indexOf(len(r.Services[serviceIndex].Keys), func(i int) bool { return r.Services[serviceIndex].Keys[i].Value == key })
	if keyIndex != -1 {
		r.Services[serviceIndex].Keys = append(r.Services[serviceIndex].Keys[:keyIndex], r.Services[serviceIndex].Keys[keyIndex+1:]...)
	}
	targetIndex := indexOf(len(r.Services[serviceIndex].Keys), func(i int) bool { return r.Services[serviceIndex].Keys[keyIndex].Targets[i].Address == t.Address })
	if targetIndex != -1 {
		r.Services[serviceIndex].Keys[keyIndex].Targets = append(r.Services[serviceIndex].Keys[keyIndex].Targets[:targetIndex], r.Services[serviceIndex].Keys[keyIndex].Targets[targetIndex+1:]...)
	}
}

//Ensures registry adheres to the restrictions set forth by the registry definition. Currently
//ensures glb keywords are not used in the service/version combination. If this constraint is
//violated ErrServiceNameNotAllowed is returned.
func (r *StandardRegistry) Validate() error {
	r.lock.RLock()
	defer r.lock.RUnlock()
	var protectedWords []string = []string{"status", "reload"}
	for _, disallowed := range protectedWords {
		serviceIndex := indexOf(len(r.Services), func(i int) bool { return r.Services[i].Value == disallowed })
		keyIndex := indexOf(len(r.Services[serviceIndex].Keys), func(i int) bool { return r.Services[serviceIndex].Keys[i].Value == disallowed })
		if serviceIndex != -1 || keyIndex != -1 {
			return registry.ErrServiceNameNotAllowed
		}
	}
	return nil
}

func (r *StandardRegistry) IncrementFailures(svc string, key string, t registry.Target, amount int) (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	serviceIndex := indexOf(len(r.Services), func(i int) bool { return r.Services[i].Value == svc })
	if serviceIndex == -1 {
		return -1, registry.ErrServiceNotFound
	}
	keyIndex := indexOf(len(r.Services[serviceIndex].Keys), func(i int) bool { return r.Services[serviceIndex].Keys[i].Value == key })
	if keyIndex == -1 {
		return -1, registry.ErrServiceNotFound
	}
	targetIndex := indexOf(len(r.Services[serviceIndex].Keys), func(i int) bool { return r.Services[serviceIndex].Keys[keyIndex].Targets[i].Address == t.Address })
	if targetIndex == -1 {
		return -1, registry.ErrServiceNotFound
	}
	r.Services[serviceIndex].Keys[keyIndex].Targets[targetIndex].Failures += amount
	return r.Services[serviceIndex].Keys[keyIndex].Targets[targetIndex].Failures, nil
}

func (r *StandardRegistry) SetRoundRobbinCounter(svc string, key string, value int) (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	serviceIndex := indexOf(len(r.Services), func(i int) bool { return r.Services[i].Value == svc })
	if serviceIndex == -1 {
		return -1, registry.ErrServiceNotFound
	}
	keyIndex := indexOf(len(r.Services[serviceIndex].Keys), func(i int) bool { return r.Services[serviceIndex].Keys[i].Value == key })
	if keyIndex == -1 {
		return -1, registry.ErrServiceNotFound
	}
	r.Services[serviceIndex].Keys[keyIndex].RoundRobbinCounter = value
	return r.Services[serviceIndex].Keys[keyIndex].RoundRobbinCounter, nil
}

func (r *StandardRegistry) GetRoundRobbinCounter(svc string, key string) (int, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	serviceIndex := indexOf(len(r.Services), func(i int) bool { return r.Services[i].Value == svc })
	if serviceIndex == -1 {
		return -1, registry.ErrServiceNotFound
	}
	keyIndex := indexOf(len(r.Services[serviceIndex].Keys), func(i int) bool { return r.Services[serviceIndex].Keys[i].Value == key })
	if keyIndex == -1 {
		return -1, registry.ErrServiceNotFound
	}
	return r.Services[serviceIndex].Keys[keyIndex].RoundRobbinCounter, nil
}
