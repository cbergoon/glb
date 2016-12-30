package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/cbergoon/glb/registry"
	"io"
	"io/ioutil"
	"os"
)

type ProxyConfig struct {
	Host                   Provider
	Basic                  bool
	DisableKeepAlives      bool
	IdleConnTimeoutSeconds int
	Registry               registry.DefaultRegistry
}

type Provider struct {
	Addr    string
	Port    string
	SslPort string
}

const (
	ConfigFile string = "glb.conf"
)

func ReadParseConfig() (ProxyConfig, error) {
	data, err := readConfig()
	if err != nil {
		os.Exit(-1)
	}
	config, err := parseConfig(bytes.NewReader(data))
	return config, err
}

func readConfig() ([]byte, error) {
	data, err := ioutil.ReadFile(ConfigFile)
	return data, err
}

func parseConfig(jsonStream io.Reader) (ProxyConfig, error) {
	dec := json.NewDecoder(jsonStream)
	var p ProxyConfig
	for {
		if err := dec.Decode(&p); err == io.EOF {
			break
		} else if err != nil {
			return p, errors.New("proxy-config: failed to parse configuration")
		}

	}
	return p, nil
}
