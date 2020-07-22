package main

import (
	"fmt"
	"net"
	"sort"
)

var target = "127.0.0.1"
var portcount = 1024

func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", target, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		fmt.Println("Open: ", address)
		conn.Close()
		results <- p
	}
}

func main() {
	ports := make(chan int, 100)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 1; i <= portcount; i++ {
			ports <- i
		}
	}()

	for i := 0; i < portcount; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}

	fmt.Println("Done")
}
