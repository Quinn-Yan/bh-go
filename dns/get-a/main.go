package main

import (
	"fmt"
	"os"

	"github.com/miekg/dns"
)

func main() {
	var msg dns.Msg
	fqdn := dns.Fqdn(os.Args[1])
	msg.SetQuestion(fqdn, dns.TypeA)
	in, err := dns.Exchange(&msg, "1.1.1.1:53")
	if err != nil {
		panic(err)
	}
	if len(in.Answer) < 1 {
		fmt.Println("No records")
		return
	}
	for _, answer := range in.Answer {
		//type assertion (a holds type, bool success/fail)
		if a, ok := answer.(*dns.A); ok {
			fmt.Println(a.A)
		}
	}
}
