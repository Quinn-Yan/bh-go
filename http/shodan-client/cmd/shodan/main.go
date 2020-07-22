package main

import (
	"fmt"
	"log"
	"os"

	"github.com/thesubtlety/bh-go/http/shodan-client/shodan"
)

func main() {

	shodanenv := "SHODAN_API_KEY"

	if len(os.Args) != 2 {
		log.Fatalln("Usage: shodan searchterm")
	}
	q := os.Args[1]

	apiKey, ok := os.LookupEnv(shodanenv)
	if !ok {
		log.Fatalf("%s env variable not set", shodanenv)
	}

	c := shodan.New(apiKey)
	info, err := c.APIInfo()
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf(
		"Query Credits: %d\nScan Credits: %d\n\n",
		info.QueryCredits,
		info.ScanCredits)

	hostSearch, err := c.HostSearch(q)
	if err != nil {
		log.Panicln(err)
	}

	if len(hostSearch.Matches) > 0 {
		fmt.Printf("%16s%8s\t%14s\t%48s\n", "IP", "Port", "Country", "Hostname")
		var hostname string
		for _, host := range hostSearch.Matches {
			for _, hnames := range host.Hostname {
				hostname = hnames
				break
			}
			fmt.Printf("%16s%8d\t%14s\t%48s\n", host.IPString, host.Port, host.Location.CountryName, hostname)
		}
	} else {
		fmt.Println("No results found")
	}

}
