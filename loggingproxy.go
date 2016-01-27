// Command loggingproxy is a http proxy server witht he ability to log requests and responses
// flowing through the proxy server for in debugging etc. Since it is in pure Go, it
// is easy to cross-compile it to almost any modern platform and OS.
package main

import (
	"flag"
	"fmt"
	"github.com/artyom/autoflags"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type config struct {
	Addr     string `flag:"addr,the address on which to serve the proxy"`
	Scheme   string `flag:"scheme,the scheme for the remote server"`
	Host     string `flag:"host,the host of the remote proxy server"`
	BasePath string `flag:"basepath,the base path of the remote proxy server"`
	LogReq   bool   `flag:"logreq,log requests to stdout"`
	LogResp  bool   `flag:"logresp,log responses to stdout"`
}

var cfg config
var proxy *httputil.ReverseProxy

func main() {

	cfg.Addr = ":8080"
	cfg.Scheme = "http"
	cfg.Host = "localhost:80"
	cfg.BasePath = ""
	cfg.LogReq = true
	cfg.LogResp = false

	autoflags.Define(&cfg)
	flag.Parse()
	proxy = httputil.NewSingleHostReverseProxy(&url.URL{Scheme: cfg.Scheme, Host: cfg.Host, Path: cfg.BasePath})
	http.ListenAndServe(cfg.Addr, http.HandlerFunc(logServe))
}

func logServe(w http.ResponseWriter, r *http.Request) {
	if cfg.LogReq {
		fmt.Printf("%s %s: %s", time.Now().Format("2006-01-02T15:04:05.999"), r.RemoteAddr, r.URL.String())
	}
	proxy.ServeHTTP(w, r)
}
