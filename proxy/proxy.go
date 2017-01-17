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
	//Todo: Move to registry to allow predictable balancing when multiple
	//Todo: service/version combinations are present.
	roundRobbinCounter int = 0 //Index for current round-robbin address.
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

//Creates a new reverse proxy that represents the configuration specified. This is done by
//creating a new http.Transport object that utilizes configuration passed in the dial
//function defined above. A http.Handler function is returned which will complete the proxy
//loop when invoked.
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
