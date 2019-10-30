package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func ProxyHTTPS(responseWriter http.ResponseWriter, req *http.Request) {
	var proxyConn, clientConn net.Conn
	var err error
	defer func() {
		if proxyConn != nil {
			log.Println("ProxyHTTPS() close proxyConn")
			_ = proxyConn.Close()
		}

		if clientConn != nil {
			log.Println("ProxyHTTPS() close clientConn")
			_ = clientConn.Close()
		}
	}()

	hij, ok := responseWriter.(http.Hijacker)
	if !ok {
		log.Println("ProxyHTTPS() error")
		return
	}

	clientConn, _, err = hij.Hijack()
	if err != nil {
		log.Println(err)
		return
	}

	proxyConn, err = net.Dial("tcp", req.URL.Host)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = clientConn.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		_, _ = io.Copy(clientConn, proxyConn)
		setDeadLine(proxyConn, clientConn)
	}()

	_, _ = io.Copy(proxyConn, clientConn)
	setDeadLine(proxyConn, clientConn)
}

func ProxyHTTP(responseWriter http.ResponseWriter, req *http.Request) {
	var resp *http.Response
	var err error

	defer func() {

		if err := recover(); err != nil {
			fmt.Println(err)
		}

		if resp != nil && resp.Body != nil {
			log.Println("ProxyHTTP() close resp.Body")
			_ = resp.Body.Close()
		}
	}()

	clientHand := newClient()

	req.RequestURI = ""

	DelHopHeaders(&req.Header)

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		AppendHostToXForwardHeader(&req.Header, &clientIP)
	}

	resp, err = clientHand.Do(req)
	if err != nil {
		log.Println("ProxyHTTP() error:", err)
		return
	}

	DelHopHeaders(&resp.Header)
	h := responseWriter.Header()
	CopyHeader(&h, &resp.Header)
	responseWriter.WriteHeader(resp.StatusCode)

	_, _ = io.Copy(responseWriter, resp.Body)
}

func newClient() *http.Client {
	httpClient := &http.Client{}
	return httpClient
}

func HttpDispatcher(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Port() == HttpsPort {
		ProxyHTTPS(rw, req)
	} else {
		ProxyHTTP(rw, req)
	}
}

func setDeadLine(proxyConn net.Conn, clientConn net.Conn) {
	_ = proxyConn.SetDeadline(time.Now())
	_ = clientConn.SetDeadline(time.Now())
}
