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
	Host                   Provider //GLB host descriptor.
	Basic                  bool //Basic mode for "default" service/version.
	DisableKeepAlives      bool //Disable keepalives causing a redial on each request.
	IdleConnTimeoutSeconds int //Timeout idle connections after in seconds; zero means no limit.
	Registry               registry.DefaultRegistry //Registry represented by the configuration.
}

type Provider struct {
	Addr    string //Address requests should bind to.
	Port    string //HTTP port; used for redirect and proxy if SslPort is not specified.
	SslPort string //HTTPS port; used for reverse proxy endpoint when specified.
}

//Reads the json configuration file, parses the contents into the configuration
//object "ProxyConfig" and, returns the resulting configuration structure. Returns
//ErrFailedToParse if the JSON could not be marshaled, exits if the file does not
//exist.
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

//Reads the file specified by the configFile argument and returns the contents as a byte
//array. Returns ErrFailedToReadFile if reading file fails.
func readConfig(configFile string) ([]byte, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, ErrFailedToReadFile
	}
	return data, nil
}

//Parses the configuration using a json decoder returns the resulting ProxyConfig. Returns
//ErrFailedToParse if an error occurs.
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
