package main

import (
	"io"
	"log"
	"net"
)

func handle(src net.Conn) {
	dst, err := net.Dial("tcp", "www.reddit.com:443")
	if err != nil {
		log.Fatalln("unable to connect to dest host")
	}
	defer dst.Close()

	//run in goroutein to prevent io.copy from blocking
	go func() {
		if _, err := io.Copy(dst, src); err != nil {
			log.Fatalln(err)
		}
	}()

	//copy dest output back to src
	if _, err := io.Copy(src, dst); err != nil {
		log.Fatalln(err)
	}
}

// if you're testing locally, you take care of the http stuff
// e.g curl -v https://localhost:42424 -k -H "www.reddit.com"
func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:42424")
	if err != nil {
		log.Fatalln("unable to bind to port")
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("unable to accept connection")
		}
		go handle(conn)
	}
}
