package main

import (
	"bufio"
	"os"
	"testing"
)

const (
	fileNameErr   = "file-error.json"
	fileNameNoErr = "file-no-error.json"
	fileErr       = "{\r\n  \"Basic\": false\r\n  \"DisableKeepAlives\": false,\r\n  \"IdleConnTimeoutSeconds\": 10,\r\n  \"Host\": {\r\n    \"Addr\": \"localhost\",\r\n    \"Port\": \":9090\",\r\n    \"SslPort\": \":8443\"\r\n  },\r\n  \"Registry\": {\r\n    \"s1\": {\r\n      \"v1\": [\r\n        \"localhost:8080\",\r\n        \"localhost:8081\"\r\n      ]\r\n    }\r\n  }\r\n}"
	fileNoErr     = "{\r\n  \"Basic\": false,\r\n  \"DisableKeepAlives\": false,\r\n  \"IdleConnTimeoutSeconds\": 10,\r\n  \"Host\": {\r\n    \"Addr\": \"localhost\",\r\n    \"Port\": \":9090\",\r\n    \"SslPort\": \":8443\"\r\n  },\r\n  \"Registry\": {\r\n    \"s1\": {\r\n      \"v1\": [\r\n        \"localhost:8080\",\r\n        \"localhost:8081\"\r\n      ]\r\n    }\r\n  }\r\n}"
)

func buildFiles() error {
	//Write Files
	fe, err := os.Create(fileNameErr)
	if err != nil {
		return err
	}
	fne, err := os.Create(fileNameNoErr)
	if err != nil {
		return err
	}
	wfe := bufio.NewWriter(fe)
	_, err = wfe.WriteString(fileErr)
	if err != nil {
		return err
	}
	wfne := bufio.NewWriter(fne)
	_, err = wfne.WriteString(fileNoErr)
	if err != nil {
		return err
	}
	wfe.Flush()
	wfne.Flush()
	return nil
}

func cleanUp() error {
	err := os.Remove(fileNameErr)
	err = os.Remove(fileNameNoErr)
	return err
}

func TestReadParseConfig(t *testing.T) {
	err := buildFiles()
	if err != nil {
		t.Error("Could not build files got ", err)
	}
	_, err = ReadParseConfig(fileNameErr)
	if err == nil {
		t.Error("Expected not nil error got ", err)
		return
	}
	proxyConfig, err := ReadParseConfig(fileNameNoErr)
	if err != nil {
		t.Error("Expected nil error got ", err)
		return
	}
	if proxyConfig.Basic != false {
		t.Error("Proxy config not built properly got ", proxyConfig)
		return
	}
	if proxyConfig.DisableKeepAlives != false {
		t.Error("Proxy config not built properly got ", proxyConfig)
		return
	}
	if proxyConfig.IdleConnTimeoutSeconds != 10 {
		t.Error("Proxy config not built properly got ", proxyConfig)
		return
	}
	if proxyConfig.Host.Addr != "localhost" {
		t.Error("Proxy config not built properly got ", proxyConfig)
		return
	}
	if proxyConfig.Host.Port != ":9090" {
		t.Error("Proxy config not built properly got ", proxyConfig)
		return
	}
	if proxyConfig.Host.SslPort != ":8443" {
		t.Error("Proxy config not built properly got ", proxyConfig)
		return
	}
	addrs, err := proxyConfig.Registry.Lookup("s1", "v1")
	if err != nil {
		t.Error("Expected not nil error got ", err)
	}
	if len(addrs) != 2 {
		t.Error("Proxy config not built properly got ", proxyConfig)
		return
	}
	err = cleanUp()
	if err != nil {
		t.Error("Could not remove files got ", err)
	}
}
