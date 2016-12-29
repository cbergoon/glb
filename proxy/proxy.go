package proxy

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
	"github.com/cbergoon/load-balancer/registry"
)

var (
	ErrInvalidService = errors.New("invalid service/version")
	roundRobbinCounter int = 0
)

// ExtractNameVersion is called to lookup the service name / version from
// the requested URL. It should update the URL's Path to reflect the target
// expectation.
var ExtractNameVersion = extractNameVersion

// LoadBalance is the default balancer which will use a random endpoint
// for the given service name/version.
var ResolveValid = resolveValid

// extractNameVersion lookup the target path and extract the name and version.
// It updates the target Path trimming version and name.
// Expected format: `/<name>/<version>/...`
func extractNameVersion(target *url.URL) (name, version string, err error) {
	path := target.Path
	if len(path) > 1 && path[0] == '/' {
		path = path[1:]
	}
	tmp := strings.Split(path, "/")
	if len(tmp) < 2 {
		return "", "", fmt.Errorf("Invalid path")
	}
	name, version = tmp[0], tmp[1]
	target.Path = "/" + strings.Join(tmp[2:], "/")
	return name, version, nil
}

// loadBalance is a basic loadBalancer which randomly
// tries to connect to one of the endpoints and try again
// in case of failure.
func resolveValid(network, serviceName, serviceVersion string, reg registry.Registry) (net.Conn, error) {
	endpoints, err := reg.Lookup(serviceName, serviceVersion)
	if err != nil {
		return nil, err
	}
	for {
		// No more endpoint, stop
		if len(endpoints) == 0 {
			break
		}

		if roundRobbinCounter >= len(endpoints) {
			roundRobbinCounter = 0
		}

		endpoint := endpoints[roundRobbinCounter]

		// Try to connect
		conn, err := net.Dial(network, endpoint)
		if err != nil {
			reg.Failure(serviceName, serviceVersion, endpoint, err)
			// Failure: remove the endpoint from the current list and try again.
			endpoints = append(endpoints[:roundRobbinCounter], endpoints[roundRobbinCounter+1:]...)
			continue
		}
		// Success: return the connection.
		roundRobbinCounter = roundRobbinCounter + 1
		return conn, nil
	}
	// No available endpoint.
	return nil, fmt.Errorf("No endpoint available for %s/%s", serviceName, serviceVersion)
}

// NewMultipleHostReverseProxy creates a reverse proxy handler
// that will randomly select a host from the passed `targets`
func NewMultipleHostReverseProxy(reg registry.Registry) http.HandlerFunc {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: func(network, addr string) (net.Conn, error) {
			addr = strings.Split(addr, ":")[0]
			tmp := strings.Split(addr, "/")
			if len(tmp) != 2 {
				return nil, ErrInvalidService
			}
			return ResolveValid(network, tmp[0], tmp[1], reg)
		},
		TLSHandshakeTimeout: 10 * time.Second,
	}
	return func(w http.ResponseWriter, req *http.Request) {
		name, version, err := ExtractNameVersion(req.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		(&httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = name + "/" + version
			},
			Transport: transport,
		}).ServeHTTP(w, req)
	}
}
