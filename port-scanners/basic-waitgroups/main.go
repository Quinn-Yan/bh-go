package main

import (
	"fmt"
	"net"
	"sync"
)

var target = "127.0.0.1"
var portcount = 1024

func request(port int, wg *sync.WaitGroup) {
	defer wg.Done()

	address := fmt.Sprintf("%s:%d", target, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return
	}
	fmt.Println("Open: ", address)
	conn.Close()
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i < portcount; i++ {
		wg.Add(1)
		go request(i, &wg)
	}
	wg.Wait()
	fmt.Println("Done")
}
