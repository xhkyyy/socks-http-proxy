package main

import (
	"io/ioutil"
	"os"
	"strings"

	"path/filepath"
)

var rule map[string]bool

func parseRealPath(filePath *string) (fp string) {
	fp = *filePath
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		fp = filepath.Join(dir, fp)
		_, err = os.Stat(fp)
		if os.IsNotExist(err) {
			panic(err)
		}
	}
	return
}

func InitRule(filePath string) {
	filePath = parseRealPath(&filePath)

	rule = make(map[string]bool)
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	for _, line := range strings.Split(string(b), "\n") {
		rule[line] = true
	}
}

func UrlMatcheHost(host string) bool {
	host = strings.TrimSuffix(host, ":443")
	//host = strings.TrimSuffix(host, ":80")
	host = strings.TrimPrefix(host, "www.")
	i := SecondLastIndex(host, '.')
	if i != -1 {
		host = host[i+1:]
	}
	//log.Println("------------------------------------>" + host)
	_, find := rule[host]
	return find
}

func SecondLastIndex(s string, c byte) int {
	count := 0
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			count += 1
			if count == 2 {
				return i
			} else {
				continue
			}
		}
	}
	return -1
}
