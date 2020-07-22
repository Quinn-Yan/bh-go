package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type router struct {
}

type logger struct {
	Inner http.Handler
}

func (l *logger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("start")
	l.Inner.ServeHTTP(w, req)
	log.Println("finish")
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/a":
		fmt.Fprintf(w, "Executing /a")
	case "/b":
		fmt.Fprintf(w, "Executing /b")
	case "/c":
		fmt.Fprintf(w, "Executing /c")
	default:
		http.Error(w, "404 Not Found", 404)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello %s\n", r.URL.Query().Get("name"))
	fmt.Fprintf(w, "Hello\n")
}

func main() {
	//default
	//http.HandleFunc("/hello", hello)
	//http.ListenAndServe("127.0.0.1:8000", nil)

	//custom router
	//var r router
	//http.ListenAndServe("127.0.0.1:8000", &r)

	//custom middleware
	//f := http.HandlerFunc(hello)
	//l := logger{Inner: f}
	//http.ListenAndServe("127.0.0.1:8000", &l)

	//using gorilla/mux
	r := mux.NewRouter()
	r.HandleFunc("/foo", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "hi foo\n")
	}).Methods("GET")
	r.HandleFunc("/users/{user:[a-z]+}", func(w http.ResponseWriter, req *http.Request) {
		user := mux.Vars(req)["user"]
		fmt.Fprintf(w, "hi %s\n", user)
	}).Methods("GET")
	http.ListenAndServe("127.0.0.1:8000", r)

}
