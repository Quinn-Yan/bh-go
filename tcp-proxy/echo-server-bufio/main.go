package main

import (
	"io"
	"log"
	"net"
)

//echo simply echos receieved data
func echo(conn net.Conn) {
	defer conn.Close()

	if _, err := io.Copy(conn, conn); err != nil {
		log.Fatalln("unable to read/write data")
	}

	/* above io.Copy does the following basically
	reader := bufio.NewReader(conn)
	s, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln("Unable to read data")
	}
	log.Printf("Read %d bytes: %s", len(s), s)

	log.Println("Writing data")
	writer := bufio.NewWriter(conn)
	if _, err = writer.WriteString(s); err != nil {
		log.Fatalln("unable to write data")
	}
	writer.Flush()
	*/
}

func main() {
	//bind to 42424
	listener, err := net.Listen("tcp", "127.0.0.1:42424")
	if err != nil {
		log.Fatalln("unable to bind port")
	}
	log.Println("Started listener on 127.0.0.1:42424")
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
