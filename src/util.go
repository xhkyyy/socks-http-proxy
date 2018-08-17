package main

import (
	"net/http"
	"strings"
)

const (
	xffForHeader = "X-Forwarded-For"
)

var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func CopyHeader(dst, src *http.Header) {
	for k, vv := range *src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func DelHopHeaders(header *http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func AppendHostToXForwardHeader(header *http.Header, host *string) {
	var build strings.Builder
	if prior, ok := (*header)[xffForHeader]; ok {
		for _, v := range prior {
			build.WriteString(v)
			build.WriteString(", ")
		}
	}
	build.WriteString(*host)
	header.Set(xffForHeader, build.String())
}
