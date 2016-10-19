package main

import (
	"flag"
	"log"

	"github.com/zeromq/gomq"
	"github.com/zeromq/gomq/zmtp"
)

func main() {
	var addr string
	flag.StringVar(&addr, "iaddr", "tcp://localhost:5555", "input data port")

	flag.Parse()

	sck := gomq.NewClient(zmtp.NewSecurityNull())
	log.Printf("dialing [%s]...\n", addr)
	err := sck.Connect(addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("dialing [%s]...\n", addr)

	for {
		msg, err := sck.Recv()
		if err != nil {
			log.Fatalf("error recv: %v\n", err)
		}

		log.Printf("recv: %v\n", string(msg))
	}
}
