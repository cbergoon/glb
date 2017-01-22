package main

import (
	"fmt"
	"log"
	"net/http"

	"os"

	"github.com/cbergoon/glb/config"
	"github.com/cbergoon/glb/proxy"
	"github.com/cbergoon/glb/registry/standardregistry"
)

const (
	CONFIG_FILE = "glb.json"   //File containing configurataion
	CERT_FILE   = "server.crt" //SSL Certificate
	KEY_FILE    = "server.key" //SSL Key
)

var serviceRegistry *serviceregistry.StandardRegistry = &serviceregistry.StandardRegistry{} //Service registry to store service-address mappings.
var BasicProxy bool = false                                                         //Enable single service "default" service/version. Removes requirement of service/version in URL.
var IdleConnTimeoutSeconds int = 1                                                  //Duration the transport should keep connections alive. Zero imposes no limit.
var DisableKeepAlives bool = false                                                  //Do not keep alive, reconnect on each request.

//Starts load balancer, redirect for HTTPS and, service endpoints.
func runLoadBalancer(addr, port, sslPort string) {
	//Redirect to HTTPS
	if sslPort != "" {
		log.Print("HTTPS config specified; starting HTTP redirect server.")
		go http.ListenAndServe(port, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req,
				"https://"+addr+sslPort+req.URL.String(),
				http.StatusMovedPermanently)
		}))
	} else {
		log.Print("HTTP only config specified; not starting HTTP redirect server.")
	}
	//GLB Service Endpoints
	http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%v\n", *serviceRegistry)
	})
	http.HandleFunc("/reload", func(w http.ResponseWriter, req *http.Request) {
		config, err := config.ReadParseConfig(CONFIG_FILE, serviceRegistry)
		if err != nil {
			log.Print(err)
			os.Exit(-1)
		}
		BasicProxy = config.Basic
		IdleConnTimeoutSeconds = config.IdleConnTimeoutSeconds
		DisableKeepAlives = config.DisableKeepAlives
		fmt.Fprintf(w, "%v\n", *serviceRegistry)
	})
	//Proxy Endpoint
	http.HandleFunc("/", proxy.NewLoadBalanceHostReverseProxy(serviceRegistry, &BasicProxy, &IdleConnTimeoutSeconds, &DisableKeepAlives))
	if sslPort != "" {
		log.Print("HTTPS config specified; listen and serve HTTPS")
		log.Print("Using Certificate File: ", CERT_FILE, " and Key File: ", KEY_FILE)
		log.Fatal(http.ListenAndServeTLS(sslPort, CERT_FILE, KEY_FILE, nil))
	} else {
		log.Print("HTTP only config specified; listen and serve HTTP")
		log.Fatal(http.ListenAndServe(sslPort, nil))
	}
}

//Application entry point gets configuration and starts the load balancer.
func main() {
	//Configure
	config, err := config.ReadParseConfig(CONFIG_FILE, serviceRegistry)
	if err != nil {
		log.Print(err)
		os.Exit(-1)
	}
	BasicProxy = config.Basic
	IdleConnTimeoutSeconds = config.IdleConnTimeoutSeconds
	DisableKeepAlives = config.DisableKeepAlives
	//Run
	runLoadBalancer(config.Host.Addr, config.Host.Port, config.Host.SslPort)
}
