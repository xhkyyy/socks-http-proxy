package main

import (
	"fmt"
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
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			rule[line] = true
		}
	}
}

func UrlMatcheHost(host string) bool {
	host = strings.TrimSuffix(host, ":443")
	_, find := rule[host]

	if !find {
		domainArray := getSubDomains(host)
		if domainArray != nil && len(domainArray) > 0 {
			for _, d := range domainArray {
				if _, find = rule[d]; find {
					return find
				}
			}
		}
	}
	return find
}

func getSubDomains(host string) []string {
	fmt.Println("<--->" + host)
	arr := strings.Split(host, ".")

	arrLen := len(arr)

	if arrLen <= 1 {
		return nil
	}

	var domains []string

	for i := 1; i < arrLen-1; i++ {

		domains = append(domains, strings.Join(arr[i:], "."))
	}

	return domains
}
