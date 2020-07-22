package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Status struct {
	Message string
	Status  string
}

func main() {
	res, err := http.Get("http://localhost/serve/test.json")
	if err != nil {
		log.Panicln(err)
	}

	var stat Status
	err = json.NewDecoder(res.Body).Decode(&stat)
	if err != nil {
		log.Panicln(err)
	}
	defer res.Body.Close()
	log.Printf("%s -> %s\n", stat.Status, stat.Message)
}

func basicpost() {
	req, err := http.NewRequest("DELETE", "http://localhost/test.txt", nil)
	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp.Body)
	fmt.Println(resp.Status)
	resp.Body.Close()

	form := url.Values{}
	form.Add("foo", "bar")
	r2, err := http.PostForm("http://localhost:80/", form)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(r2.Status)
}

func basicread() {
	resp, err := http.Get("http://localhost/serve/test.txt")
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(string(body))
	resp.Body.Close()
}
