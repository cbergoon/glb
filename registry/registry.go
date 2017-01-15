package registry

import (
	"errors"
	"sync"
)

var lock sync.RWMutex

var (
	ErrServiceNotFound = errors.New("registry: target name/version not found")
	ErrServiceNameNotAllowed = errors.New("registry: service name not allowed; non-allowable service names [reload|status]")
)

type Registry interface {
	Add(name, version, endpoint string)
	Delete(name, version, endpoint string)
	Lookup(name, version string) ([]string, error)
	Validate() error
}

type DefaultRegistry map[string]map[string][]string

func (r DefaultRegistry) Lookup(name, version string) ([]string, error) {
	lock.RLock()
	targets, ok := r[name][version]
	lock.RUnlock()
	if !ok {
		return nil, ErrServiceNotFound
	}
	return targets, nil
}

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
