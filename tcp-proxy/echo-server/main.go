package main

import (
	"io"
	"log"
	"net"
)

//echo simply echos receieved data
func echo(conn net.Conn) {
	defer conn.Close()

	//create a buffer to store received data
	b := make([]byte, 512)
	for {
		//rcv data via conn.Read into buffer
		s, err := conn.Read(b[0:])
		if err == io.EOF {
			log.Println("Client disconnected")
			break
		}
		if err != nil {
			log.Println("Unexpected error")
			break
		}
		log.Printf("Received %d bytes: %s\n", s, string(b))

		//send data via conn.Write
		log.Println("Writing data")
		if _, err := conn.Write(b[0:s]); err != nil {
			log.Fatalln("unable to write data")
		}
	}
}

func main() {
	//bind to 42424
	listener, err := net.Listen("tcp", ":42424")
	if err != nil {
		log.Fatalln("unable to bind port")
	}
	log.Println("Started listener on 0.0.0.0:42424")
	for {
		//wait for conn, create when established
		conn, err := listener.Accept()
		log.Println("Connection received")
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		go echo(conn)
	}
}
