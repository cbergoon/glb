package proxy_test

import (
	"github.com/cbergoon/glb/proxy"
	"github.com/cbergoon/glb/registry"
	"testing"
)

var serviceRegistry = registry.ServiceRegistry{
	"service1": {
		"v1": {
			"localhost:8888",
		},
	},
	"service2": {
		"v2": {
			"localhost:7777",
			"localhost:6666",
		},
	},
}

func TestNewLoadBalanceHostReverseProxy(t *testing.T) {
	var FALSE = false
	var ZERO = 0
	handlerFunc := proxy.NewLoadBalanceHostReverseProxy(serviceRegistry, &FALSE, &ZERO, &FALSE)
	if handlerFunc == nil {
		t.Error("Expected handler func got ", handlerFunc)
	}
}
