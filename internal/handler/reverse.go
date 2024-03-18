package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type reverseProxy struct {
	proxy    *httputil.ReverseProxy
	url      *url.URL
	endpoint string
}

func NewReverseProxy(url *url.URL, endpoint string) *reverseProxy {
	fmt.Println("reading reverse proxy")
	rv := reverseProxy{
		proxy:    httputil.NewSingleHostReverseProxy(url),
		url:      url,
		endpoint: endpoint,
	}
	return &rv
}

func (rv *reverseProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.URL.Host = rv.url.Host
	req.URL.Scheme = rv.url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = rv.url.Host

	path := req.URL.Path
	req.URL.Path = strings.TrimLeft(path, rv.endpoint)

	rv.proxy.ServeHTTP(w, req)
}
