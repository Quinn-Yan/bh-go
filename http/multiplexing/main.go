package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

var (
	hostProxy = make(map[string]string)
	proxies   = make(map[string]*httputil.ReverseProxy)
)

func init() {
	hostProxy["attacker1.com"] = "http://127.0.0.1:8081"
	hostProxy["attacker2.com"] = "http://127.0.0.1:8082"

	for k, v := range hostProxy {
		remote, err := url.Parse(v)
		if err != nil {
			log.Fatal("Unable to parse proxy target")
		}
		proxies[k] = httputil.NewSingleHostReverseProxy(remote)
	}
}

func main() {
	r := mux.NewRouter()
	for host, proxy := range proxies {
		r.Host(host).Handler(proxy)
	}

	log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", r))
}
