package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
)

var target = flag.String("target", "127.0.0.1", "target ip address")
var portsflag = flag.String("ports", "80", "port (from-to or comma separated - no spaces)")
var portlist []int

func parsetargets() {
	/*
		accept cidr format
		config for timeout on dial
		instead of ports channel rename for addresses
		in main, make list of ip+ports (for port in list for ip in list)
			randomly sort, send to workers
		status? every n seconds or keypress (channel) print index of list
	*/
}

func parseports() []int {
	var portliststr []string

	if strings.Contains(*portsflag, ",") {
		portliststr = strings.Split(*portsflag, ",")
		for _, port := range portliststr {
			if strings.Contains(port, "-") {
				continue
			}
			j, err := strconv.Atoi(port)
			if err != nil {
				log.Fatal("Error converting to integers")
			}
			portlist = append(portlist, j)
		}
	}

	if strings.Contains(*portsflag, "-") {
		portliststr = strings.Split(*portsflag, ",")
		for _, ranges := range portliststr {
			if strings.Contains(ranges, "-") {
				portliststr = strings.Split(ranges, "-")
				if len(portliststr) != 2 {
					fmt.Println("len portlist", len(portliststr))
					log.Fatal("portsflag error")
				}
				low, _ := strconv.Atoi(portliststr[0])
				high, err := strconv.Atoi(portliststr[1])
				if err != nil {
					log.Fatal("Port format error")
				}
				if low > high {
					log.Fatal("Port range error")
				}

				for i := low; i <= high; i++ {
					portlist = append(portlist, i)
				}
			}
		}

	}

	if len(portlist) == 0 { //hack, assuming prior checks haven't added any ports
		j, err := strconv.Atoi(*portsflag)
		if err != nil {
			log.Fatal("Port format error")
		}
		portlist = append(portlist, j)
	}

	for _, i := range portlist {
		if i < 1 || i > 65535 {
			log.Fatal("Invalid port number")
		}
	}

	return portlist
}

func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", *target, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		fmt.Println(address)
		conn.Close()
		results <- p
	}
}

func main() {
	flag.Parse()
	portlist = parseports()

	ports := make(chan int, 100)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for _, p := range portlist {
			ports <- p
		}
	}()

	for i := 0; i < len(portlist); i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)

	if len(openports) == 0 {
		fmt.Println("No open ports...")
	}

	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d/open\n", port)
	}
}
