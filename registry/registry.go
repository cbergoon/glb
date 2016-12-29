package registry

import (
	"errors"
	"log"
	"sync"
)

var lock sync.RWMutex

var (
	ErrServiceNotFound = errors.New("service name/version not found")
)

type Registry interface {
	Add(name, version, endpoint string)
	Delete(name, version, endpoint string)
	Failure(name, version, endpoint string, err error)
	Lookup(name, version string) ([]string, error)
}

// {
//   "serviceName": {
//     "serviceVersion": [
//       "endpoint1:port",
//       "endpoint2:port"
//     ],
//   },
// }
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

func (r DefaultRegistry) Failure(name, version, endpoint string, err error) {
	log.Printf("Error accessing %s/%s (%s): %s", name, version, endpoint, err)
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
