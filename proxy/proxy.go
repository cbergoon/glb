package proxy

import (
	"errors"
	"fmt"
	"github.com/cbergoon/glb/registry"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var (
	ErrInvalidService      = errors.New("invalid service/version")
	roundRobbinCounter int = 0
)

var ExtractNameVersion = extractNameVersion

var ResolveValid = resolveValid

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

func resolveValid(network, serviceName, serviceVersion string, reg registry.Registry) (net.Conn, error) {
	endpoints, err := reg.Lookup(serviceName, serviceVersion)
	if err != nil {
		return nil, err
	}
	for {
		if len(endpoints) == 0 {
			break
		}

		if roundRobbinCounter >= len(endpoints) {
			roundRobbinCounter = 0
		}

		endpoint := endpoints[roundRobbinCounter]

		conn, err := net.Dial(network, endpoint)
		if err != nil {
			reg.Failure(serviceName, serviceVersion, endpoint, err)
			endpoints = append(endpoints[:roundRobbinCounter], endpoints[roundRobbinCounter+1:]...)
			continue
		}

		roundRobbinCounter = roundRobbinCounter + 1
		return conn, nil
	}
	return nil, fmt.Errorf("No endpoint available for %s/%s", serviceName, serviceVersion)
}

func NewMultipleHostReverseProxy(reg registry.Registry, basic *bool, idleConTimeout *int, disableKeepAlive *bool) http.HandlerFunc {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: func(network, addr string) (net.Conn, error) {
			fmt.Println("In Dail")
			addr = strings.Split(addr, ":")[0]
			tmp := strings.Split(addr, "/")
			if len(tmp) != 2 {
				return nil, ErrInvalidService
			}
			return ResolveValid(network, tmp[0], tmp[1], reg)
		},
		TLSHandshakeTimeout: 10 * time.Second,
		IdleConnTimeout:     time.Duration(*idleConTimeout) * time.Second,
		DisableKeepAlives:   *disableKeepAlive,
	}
	return func(w http.ResponseWriter, req *http.Request) {
		var name, version string
		var err error
		if !(*basic) {
			name, version, err = ExtractNameVersion(req.URL)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if *basic {
			name = "default"
			version = "default"
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
