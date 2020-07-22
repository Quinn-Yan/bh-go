package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/miekg/dns"
)

func main() {
	var (
		flDomain      = flag.String("domain", "", "domain to perform guessing against")
		flWordlist    = flag.String("wordlist", "", "wordlist to use for guessing")
		flWorkerCount = flag.Int("c", 100, "number of workers")
		flServerAddr  = flag.String("server", "1.1.1.1:53", "dns server to resolve against")
	)
	flag.Parse()

	if *flDomain == "" || *flWordlist == "" {
		fmt.Println("-domain and -wordlist required")
		os.Exit(1)
	}

	var results []result
	fqdns := make(chan string, *flWorkerCount) //buffered channel lets us hold >1 message before blocking sender
	output := make(chan []result)
	status := make(chan empty)

	fh, err := os.Open(*flWordlist)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)

	//start workers so nothing blocks
	for i := 0; i < *flWorkerCount; i++ {
		go worker(status, fqdns, output, *flServerAddr)
	}

	//for guess in file, send fqdn to channel
	for scanner.Scan() {
		fqdns <- fmt.Sprintf("%s.%s", scanner.Text(), *flDomain)
	}

	go func() {
		for r := range output {
			results = append(results, r...) //slice append to slice use ...
		}
		var e empty
		status <- e
	}()

	close(fqdns)                          //work already sent, close first
	for i := 0; i < *flWorkerCount; i++ { //receive on status chanel once per worker
		<-status
	}
	close(output)
	<-status

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 4, ' ', 0)
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\n", r.Hostname, r.IPAddress)
	}
	w.Flush()
}

type result struct {
	IPAddress string
	Hostname  string
}

func lookup(fqdn, serverAddr string) []result {
	var results []result
	var cfqdn = fqdn //not modifying original
	for {            //keep following the cname chain until no more exist
		cnames, err := lookupCNAME(cfqdn, serverAddr)
		if err == nil && len(cnames) > 0 {
			cfqdn = cnames[0]
			continue //process next cname
		}
		ips, err := lookupA(cfqdn, serverAddr)
		if err != nil {
			break //no a records
		}
		for _, ip := range ips {
			results = append(results, result{IPAddress: ip, Hostname: fqdn})
		}
		break //processed all, break
	}
	return results
}

func lookupA(fqdn, serverAddr string) ([]string, error) {
	var msg dns.Msg
	var ips []string
	msg.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)
	in, err := dns.Exchange(&msg, serverAddr)
	if err != nil {
		return ips, err
	}
	if len(in.Answer) < 1 {
		return ips, errors.New("no answer")
	}
	for _, answer := range in.Answer {
		if a, ok := answer.(*dns.A); ok {
			ips = append(ips, a.A.String())
		}
	}
	return ips, nil
}

func lookupCNAME(fqdn, serverAddr string) ([]string, error) {
	var msg dns.Msg
	var fqdns []string
	msg.SetQuestion(dns.Fqdn(fqdn), dns.TypeCNAME)
	in, err := dns.Exchange(&msg, serverAddr)
	if err != nil {
		return fqdns, err
	}
	if len(in.Answer) < 1 {
		return fqdns, errors.New("no answer")
	}
	for _, answer := range in.Answer {
		if c, ok := answer.(*dns.CNAME); ok {
			fqdns = append(fqdns, c.Target)
		}
	}
	return fqdns, nil
}

type empty struct{} //just tracks when worker finished

func worker(status chan empty, fqdns chan string, output chan []result, serverAddr string) {
	for fqdn := range fqdns {
		results := lookup(fqdn, serverAddr)
		if len(results) > 0 {
			output <- results
		}
	}
	var e empty
	status <- e //necessary to prevent race condition, exiting before results sent
}
