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
	ErrInvalidTarget = errors.New("proxy: invalid service/version")
	ErrInvalidPath   = errors.New("proxy: invalid path to resource")
)

var ParseTarget = parseTarget
var DialTarget = dialTarget

//Extracts the service name and version from the URL provided. Returns ErrInvalidPath if
//the service and/or version are missing from the URL provided.
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

//Establishes the network connection to appropriate address that is determined by the service and version.
//Executes a look up with the registry based on the parameters service and version. If a connection is not
//able to be established to any of the available addresses, an error is returned detailing the failure.
//Note that this function may or may not be called at deterministic intervals depending on the configuration,
//request volume and, load balancer settings.
func dialTarget(network, serviceName, serviceKey string, reg registry.Registry) (net.Conn, error) {
	localRoundRobbin, err := reg.GetRoundRobbinCounter(serviceName, serviceKey)
	if localRoundRobbin < 0 || err != nil {
		log.Print(err)
		return nil, err
	}
	endpoints, err := reg.Lookup(serviceName, serviceKey)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	for {
		if len(endpoints) == 0 {
			break
		}

		if localRoundRobbin >= len(endpoints) {
			localRoundRobbin = 0
			reg.SetRoundRobbinCounter(serviceName, serviceKey, 0)
		}

		endpoint := endpoints[localRoundRobbin].Address

		conn, err := net.Dial(network, endpoint)
		if err != nil {
			log.Printf("proxy: error could not access %s/%s at %s", serviceName, serviceKey, endpoint)
			endpoints = append(endpoints[:localRoundRobbin], endpoints[localRoundRobbin+1:]...)
			continue
		}

		localRoundRobbin = localRoundRobbin + 1
		reg.SetRoundRobbinCounter(serviceName, serviceKey, localRoundRobbin)
		return conn, nil
	}
	e := fmt.Errorf("proxy: error no endpoint available for %s/%s", serviceName, serviceKey)
	log.Print(e)
	return nil, e
}

//Creates a new reverse proxy that represents the configuration specified. This is done by
//creating a new http.Transport object that utilizes configuration passed in the dial
//function defined above. A http.Handler function is returned which will complete the proxy
//loop when invoked.
func NewLoadBalanceHostReverseProxy(reg registry.Registry, basic *bool, idleConTimeout *int, disableKeepAlive *bool) http.HandlerFunc {
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
		var name, key string
		var err error
		if !(*basic) {
			name, key, err = ParseTarget(req.URL)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if *basic {
			name = "default"
			key = "default"
		}
		(&httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = name + "/" + key
			},
			Transport: transport,
		}).ServeHTTP(w, req)
	}
}
