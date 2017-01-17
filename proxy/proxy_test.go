package proxy_test

import (
	"testing"
	"github.com/cbergoon/glb/registry"
	"github.com/cbergoon/glb/proxy"
)

var serviceRegistry = registry.DefaultRegistry{
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

func TestNewMultipleHostReverseProxy(t *testing.T) {
	var FALSE = false
	var ZERO = 0
	handlerFunc := proxy.NewMultipleHostReverseProxy(serviceRegistry, &FALSE, &ZERO, &FALSE)
	if handlerFunc == nil {
		t.Error("Expected handler func got ", handlerFunc)
	}
}