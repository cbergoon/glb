package proxy_test

import (
	"github.com/cbergoon/glb/proxy"
	"github.com/cbergoon/glb/registry/standardregistry"
	"testing"
)

var serviceRegistry = serviceregistry.StandardRegistry{}

func TestNewLoadBalanceHostReverseProxy(t *testing.T) {
	var FALSE = false
	var ZERO = 0
	handlerFunc := proxy.NewLoadBalanceHostReverseProxy(serviceRegistry, &FALSE, &ZERO, &FALSE)
	if handlerFunc == nil {
		t.Error("Expected handler func got ", handlerFunc)
	}
}
