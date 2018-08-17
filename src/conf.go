package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Conf struct {
	DomainFile string `json:"domain_file"`
	Addr       string `json:"addr"`
	Socks5Addr string `json:"socks5_addr"`
}

func LoadConf(filePath string) *Conf {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var confInfo Conf
	json.Unmarshal(b, &confInfo)
	return &confInfo
}
