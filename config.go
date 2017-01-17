package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/cbergoon/glb/registry"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	ErrFailedToParse    = errors.New("proxy-config: failed to parse configuration")
	ErrFailedToReadFile = errors.New("proxy-config: failed to read configuration file")
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

func ReadParseConfig(configFile string) (ProxyConfig, error) {
	data, err := readConfig(configFile)
	if err != nil {
		log.Print(err)
		os.Exit(-1)
	}
	config, err := parseConfig(bytes.NewReader(data))
	if err != nil {
		return config, ErrFailedToParse
	}
	return config, nil
}

func readConfig(configFile string) ([]byte, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, ErrFailedToReadFile
	}
	return data, nil
}

func parseConfig(jsonStream io.Reader) (ProxyConfig, error) {
	dec := json.NewDecoder(jsonStream)
	var p ProxyConfig
	for {
		if err := dec.Decode(&p); err == io.EOF {
			break
		} else if err != nil {
			return p, ErrFailedToParse
		}

	}
	return p, nil
}
