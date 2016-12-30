package main

import (
	"fmt"
	"log"
	"net/http"

	"os"

	"github.com/cbergoon/glb/proxy"
	"github.com/cbergoon/glb/registry"
)

var ServiceRegistry *registry.DefaultRegistry = &registry.DefaultRegistry{}
var BasicProxy bool = false
var IdleConnTimeoutSeconds int = 1
var DisableKeepAlives bool = false


func runLoadBalancer(addr, port, sslPort string) {
	//Redirect to HTTPS
	go http.ListenAndServe(port, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req,
			"https://"+addr+sslPort+req.URL.String(),
			http.StatusMovedPermanently)
	}))
	//Service Endpoints
	http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%v\n", *ServiceRegistry)
	})
	http.HandleFunc("/reload", func(w http.ResponseWriter, req *http.Request) {
		config, err := ReadParseConfig()
		if err != nil {
			os.Exit(-1)
		}
		*ServiceRegistry = config.Registry
		BasicProxy = config.Basic
		IdleConnTimeoutSeconds = config.IdleConnTimeoutSeconds
		DisableKeepAlives = config.DisableKeepAlives
		fmt.Fprintf(w, "%v\n", *ServiceRegistry)
	})
	//Proxy Endpoint
	http.HandleFunc("/", proxy.NewMultipleHostReverseProxy(ServiceRegistry, &BasicProxy, &IdleConnTimeoutSeconds, &DisableKeepAlives))

	log.Fatal(http.ListenAndServeTLS(sslPort, "server.crt", "server.key", nil))
}

func main() {
	//Configure
	config, err := ReadParseConfig()
	if err != nil {
		os.Exit(-1)
	}
	*ServiceRegistry = config.Registry
	BasicProxy = config.Basic
	IdleConnTimeoutSeconds = config.IdleConnTimeoutSeconds
	DisableKeepAlives = config.DisableKeepAlives
	//Run
	runLoadBalancer(config.Host.Addr, config.Host.Port, config.Host.SslPort)
}
