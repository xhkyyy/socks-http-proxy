package main

import (
	"context"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net"
	"net/http"
)

func Socks5ProxyHTTPS(rw http.ResponseWriter, req *http.Request) {
	var dialer proxy.Dialer
	var proxyConn, clientConn net.Conn
	var err error

	defer func() {
		if err := recover(); err != nil {
			fmt.Print("Socks5ProxyHTTPS() panic:")
			fmt.Println(err)

			if dialer != nil {
				socks5HandleCache <- dialer
			}
		}

		if proxyConn != nil {
			log.Println("Socks5ProxyHTTPS() close proxyConn")
			proxyConn.Close()
		}

		if clientConn != nil {
			log.Println("Socks5ProxyHTTPS() close clientConn")
			clientConn.Close()
		}
	}()
	hij, ok := rw.(http.Hijacker)
	if !ok {
		log.Println("error")
		return
	}

	clientConn, _, err = hij.Hijack()
	if err != nil {
		log.Println(err)
		return
	}

	dialer = <-socks5HandleCache
	proxyConn, err = dialer.Dial("tcp", req.URL.Host)
	socks5HandleCache <- dialer
	dialer = nil
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
		io.Copy(clientConn, proxyConn)
	}()

	io.Copy(proxyConn, clientConn)
}

func Socks5ProxyHTTP(wr http.ResponseWriter, req *http.Request) {
	var resp *http.Response
	defer func() {
		if err := recover(); err != nil {
			fmt.Print("Socks5ProxyHTTP panic")
			fmt.Println(err)
		}

		if resp != nil && resp.Body != nil {
			log.Println("Socks5ProxyHTTP() close resp.Body")
			resp.Body.Close()
		}
	}()
	clientHand := newSocks5Client()

	req.RequestURI = ""

	DelHopHeaders(&req.Header)

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		AppendHostToXForwardHeader(&req.Header, &clientIP)
	}

	resp, err := clientHand.Do(req)
	if err != nil {
		log.Println(err)
	}

	DelHopHeaders(&resp.Header)
	h := wr.Header()
	CopyHeader(&h, &resp.Header)
	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
}

func newSocks5Client() *http.Client {
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				dialer := <-socks5HandleCache
				defer func() {
					if err := recover(); err != nil {
						log.Println(err)
					}
					socks5HandleCache <- dialer
				}()
				conn, err := dialer.Dial(network, addr)
				return conn, err
			},
		},
	}
	return httpClient
}

func Socks5Dispatcher(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Port() == HTTPSPORT {
		Socks5ProxyHTTPS(rw, req)
	} else {
		Socks5ProxyHTTP(rw, req)
	}
}
