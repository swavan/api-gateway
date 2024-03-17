package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type reverseProxy struct {
	proxy    *httputil.ReverseProxy
	url      *url.URL
	endpoint string
}

// func ReverseProxy(url *url.URL) *httputil.ReverseProxy {
// 	return httputil.NewSingleHostReverseProxy(url)
// }

func NewReverseProxy(url *url.URL, endpoint string) *reverseProxy {
	fmt.Println("reading reverse proxy")
	rv := reverseProxy{
		proxy:    httputil.NewSingleHostReverseProxy(url),
		url:      url,
		endpoint: endpoint,
	}
	return &rv
}

func (rv *reverseProxy) Redirect(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("[ TinyRP ] Request received at %s at %s\n", req.URL, time.Now().UTC())
	// Update the headers to allow for SSL redirection

	req.URL.Host = rv.url.Host
	req.URL.Scheme = rv.url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = rv.url.Host

	//trim reverseProxyRoutePrefix
	path := req.URL.Path
	req.URL.Path = strings.TrimLeft(path, rv.endpoint)

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	fmt.Printf("[ TinyRP ] Redirecting request to %s at %s\n", req.URL, time.Now().UTC())

	rv.proxy.ServeHTTP(w, req)
}

func Redirect(proxy *httputil.ReverseProxy, url *url.URL, endpoint string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("[ TinyRP ] Request received at %s at %s\n", req.URL, time.Now().UTC())
		// Update the headers to allow for SSL redirection

		req.URL.Host = url.Host
		req.URL.Scheme = url.Scheme
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = url.Host

		//trim reverseProxyRoutePrefix
		path := req.URL.Path
		req.URL.Path = strings.TrimLeft(path, endpoint)

		// Note that ServeHttp is non blocking and uses a go routine under the hood
		fmt.Printf("[ TinyRP ] Redirecting request to %s at %s\n", req.URL, time.Now().UTC())

		proxy.ServeHTTP(w, req)
	}
}
