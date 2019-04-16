package main

import (
	"flag"
	// "fmt"
	"golang.org/x/net/proxy"
	"log"
	"net/http"
	"strings"
)

const (
	socks5HandleCacheSize = 100
	HttpsPort             = "443"
)

var ConfInfo *Conf
var socks5HandleCache = make(chan proxy.Dialer, socks5HandleCacheSize)

type ProxyServer struct {
}

func (p *ProxyServer) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	done := make(chan bool)
	// t := time.Now()
	go func() {
		if UrlMatchDomain(req.Host) {
			Socks5Dispatcher(wr, req)
		} else {
			HttpDispatcher(wr, req)
		}
		done <- true
	}()
	<-done
	// fmt.Println(req.URL.String() + "------------->" + time.Since(t).String())
}

func initConf(addr *string) {
	ConfInfo = LoadConf(*addr)
	InitRule(ConfInfo.DomainFile)
}

func main() {
	var addr = flag.String("f", "", "config file")
	flag.Parse()

	if *addr == "" || len(strings.TrimSpace(*addr)) <= 0 {
		panic("No config file found")
	}

	initConf(addr)

	for i := 0; i < socks5HandleCacheSize; i++ {
		dialerTmp, _ := proxy.SOCKS5("tcp", ConfInfo.Socks5Addr, nil, proxy.Direct)
		socks5HandleCache <- dialerTmp
	}

	handler := &ProxyServer{}
	if err := http.ListenAndServe(ConfInfo.Addr, handler); err != nil {
		log.Println(err)
	}
}
