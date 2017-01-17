package registry_test

import (
	"strconv"
	"testing"
	"github.com/cbergoon/glb/registry"
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

var service1 = []string{"localhost:8888"}
var service2 = []string{"localhost:7777", "localhost:6666"}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func TestLookup(t *testing.T) {
	//Test Valid Lookup
	addrs, err := serviceRegistry.Lookup("service1", "v1")
	if err != nil {
		t.Error("Expected nil error got ", err)
		return
	}
	if addrs == nil {
		t.Error("Expected not nil result got nil")
		return
	}
	for _, addr := range addrs {
		if !contains(service1, addr) {
			t.Error("Expected localhost:8888 got ", addr)
		}
	}
	//Test Valid Lookup with Multiple Addresses
	addrsm, err := serviceRegistry.Lookup("service2", "v2")
	if err != nil {
		t.Error("Expected nil error got ", err)
		return
	}
	if addrsm == nil {
		t.Error("Expected not nil result got nil")
		return
	}
	for _, addr := range addrsm {
		if !contains(service2, addr) {
			t.Error("Expected ", service2, " got ", addrsm)
		}
	}
	//Test Valid Lookup Missing Entry
	addrsmn, err := serviceRegistry.Lookup("notexist", "ne1")
	if err == nil {
		t.Error("Expected not nil error got ", err)
		return
	}
	if len(addrsmn) != 0 {
		t.Error("Expected empty slice got ", addrsmn)
		return
	}
	return
}

func TestAdd(t *testing.T) {
	var s = []string{
		"localhost:8080",
		"localhost:80",
		"localhost",
		"10.10.10.10:9090",
		"10.10.10.10:90",
		"10.10.10.10",
		"1.1.1.1:8080",
		"100.200.0.1:90",
		"255.255.255.255",
	}
	//Build Registry
	for i, a := range s {
		//All Addresses to service/version
		serviceRegistry.Add("allname", "allversion", a)
		//One Address per service/version
		serviceRegistry.Add("eachname"+strconv.Itoa(i), "eachversion"+strconv.Itoa(i), a)
		//Multiple version Addresses per service
		serviceRegistry.Add("multi", "multi"+strconv.Itoa(i), a)
	}
	//Test All Addresses to service/version
	allAddrs, err := serviceRegistry.Lookup("allname", "allversion")
	if err != nil {
		t.Error("Expected nil error got ", err)
		return
	}
	for _, addr := range allAddrs {
		if !contains(allAddrs, addr) {
			t.Error("Expected ", s, " got ", allAddrs)
		}
	}
	for i, a := range s {
		addrs, err := serviceRegistry.Lookup("eachname"+strconv.Itoa(i), "eachversion"+strconv.Itoa(i))
		if err != nil {
			t.Error("Expected nil error got ", err)
			return
		}
		if !contains(addrs, a) {
			t.Error("Expected ", a, " got ", addrs)
			return
		}
	}
	//Test Multi Version to service/version
	for i, a := range s {
		multiAddrs, err := serviceRegistry.Lookup("multi", "multi"+strconv.Itoa(i))
		if err != nil {
			t.Error("Expected nil error got ", err)
			return
		}
		if multiAddrs[0] != a {
			t.Error("Expected ", a, " got ", multiAddrs)
			return
		}
	}
	return
}

func TestDelete(t *testing.T) {
	serviceRegistry.Delete("service1", "v1", "localhost:8888")
	_, err := serviceRegistry.Lookup("service1", "v1")
	if err == nil {
		t.Error("Expected ", registry.ErrServiceNotFound, " got ", err)
	}
	serviceRegistry.Delete("service2", "v2", "localhost:7777")
	addrs, err := serviceRegistry.Lookup("service2", "v2")
	if err != nil && len(addrs) == 1{
		t.Error("Expected ", registry.ErrServiceNotFound, " got ", err)
	}
	return
}

func TestValidate(t *testing.T) {
	err := serviceRegistry.Validate()
	if err != nil {
		t.Error("Expected nil error got ", err)
	}
	serviceRegistry.Add("reload", "v1", "localhost:8080")
	err = serviceRegistry.Validate()
	if err == nil {
		t.Error("Expected ", registry.ErrServiceNameNotAllowed, " got ", err)
	}
	serviceRegistry.Delete("reload", "v1", "localhost:8080")
	serviceRegistry.Add("status", "v1", "localhost:8080")
	err = serviceRegistry.Validate()
	if err == nil {
		t.Error("Expected ", registry.ErrServiceNameNotAllowed, " got ", err)
	}
	serviceRegistry.Delete("status", "v1", "localhost:8080")
}
