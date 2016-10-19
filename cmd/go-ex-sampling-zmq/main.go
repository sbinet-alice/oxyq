package main

import (
	"flag"
	"log"

	"github.com/zeromq/gomq"
	"github.com/zeromq/gomq/zmtp"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "tcp://localhost:5555", "output data port")

	flag.Parse()

	sck := gomq.NewClient(zmtp.NewSecurityNull())
	defer sck.Close()

	log.Printf("dialing [%s]...\n", addr)
	err := sck.Connect(addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("dialing [%s]...\n", addr)

	for {
		msg := []byte("HELLO")
		err = sck.Send(msg)
		if err != nil {
			log.Fatalf("error send: %v\n", err)
		}
	}
}
