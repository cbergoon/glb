package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/cbergoon/glb/registry"
	"io"
	"io/ioutil"
	"os"
	"log"
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

const (
	ConfigFile string = "glb.conf"
)

func ReadParseConfig() (ProxyConfig, error) {
	data, err := readConfig()
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

func readConfig() ([]byte, error) {
	data, err := ioutil.ReadFile(ConfigFile)
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

