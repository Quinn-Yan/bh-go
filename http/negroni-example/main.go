package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type badAuth struct {
	Username string
	Password string
}

func (b *badAuth) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	//username := r.URL.Query().Get("username")
	//password := r.URL.Query().Get("password")
	username, password, _ := r.BasicAuth()
	if username != b.Username || password != b.Password {
		http.Error(w, "Unauthorized", 401)
		return
	}
	ctx := context.WithValue(r.Context(), "username", username)
	r = r.WithContext(ctx)
	next(w, r)
}

func hello(w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("username").(string)
	fmt.Fprintf(w, "Hi %s\n", username)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", hello).Methods("GET", "POST")
	n := negroni.Classic()
	n.Use(&badAuth{
		Username: "admin",
		Password: "password",
	})
	n.UseHandler(r)
	http.ListenAndServe("127.0.0.1:8000", n)
}
