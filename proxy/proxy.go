package proxy

import (
	"errors"
	"fmt"
	"github.com/cbergoon/glb/registry"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var (
	ErrInvalidTarget       = errors.New("proxy: invalid service/version")
	ErrInvalidPath         = errors.New("proxy: invalid path to resource")
	roundRobbinCounter int = 0
)

var ParseTarget = parseTarget
var DialTarget = dialTarget

func parseTarget(target *url.URL) (name, version string, err error) {
	path := target.Path
	if len(path) > 1 && path[0] == '/' {
		path = path[1:]
	}
	tmp := strings.Split(path, "/")
	if len(tmp) < 2 {
		log.Print(ErrInvalidPath)
		return "", "", ErrInvalidPath
	}
	name, version = tmp[0], tmp[1]
	target.Path = "/" + strings.Join(tmp[2:], "/")
	return name, version, nil
}

func dialTarget(network, serviceName, serviceVersion string, reg registry.Registry) (net.Conn, error) {
	endpoints, err := reg.Lookup(serviceName, serviceVersion)
	if err != nil {
		log.Print(err)
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
			log.Printf("proxy: error could not access %s/%s at %s", serviceName, serviceVersion, endpoint)
			endpoints = append(endpoints[:roundRobbinCounter], endpoints[roundRobbinCounter+1:]...)
			continue
		}

		roundRobbinCounter = roundRobbinCounter + 1
		return conn, nil
	}
	e := fmt.Errorf("No endpoint available for %s/%s", serviceName, serviceVersion)
	log.Print(e)
	return nil, e
}

func NewMultipleHostReverseProxy(reg registry.Registry, basic *bool, idleConTimeout *int, disableKeepAlive *bool) http.HandlerFunc {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: func(network, addr string) (net.Conn, error) {
			addr = strings.Split(addr, ":")[0]
			tmp := strings.Split(addr, "/")
			if len(tmp) != 2 {
				log.Print(ErrInvalidTarget)
				return nil, ErrInvalidTarget
			}
			return DialTarget(network, tmp[0], tmp[1], reg)
		},
		TLSHandshakeTimeout: 10 * time.Second,
		IdleConnTimeout:     time.Duration(*idleConTimeout) * time.Second,
		DisableKeepAlives:   *disableKeepAlive,
	}
	return func(w http.ResponseWriter, req *http.Request) {
		var name, version string
		var err error
		if !(*basic) {
			name, version, err = ParseTarget(req.URL)
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
