package main

import (
	"fmt"
	"log"
	"os"
)

type FooReader struct{}

func (fooReader *FooReader) Read(b []byte) (int, error) {
	fmt.Print("in > ")
	return os.Stdin.Read(b)
}

type FooWriter struct{}

func (FooWriter *FooWriter) Write(b []byte) (int, error) {
	fmt.Print("out > ")
	return os.Stdout.Write(b)
}

func main() {
	var (
		reader FooReader
		writer FooWriter
	)

	/*
		//this accomplishes the same as the below manual read, write
			if _, err := io.Copy(&writer, &reader); err != nil {
			log.Fatalln("unable ot read/write data")
		}
	*/

	input := make([]byte, 4096)
	s, err := reader.Read(input)
	if err != nil {
		log.Fatalln("unable ot read data")
	}
	fmt.Printf("Read %d bytes from stdin\n", s)

	s, err = writer.Write(input)
	if err != nil {
		log.Fatalln("unabel to write data")
	}
	fmt.Printf("Wrote %d bytes from stdin\n", s)
}
