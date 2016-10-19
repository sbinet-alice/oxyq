package main

import (
	"flag"
	"log"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	var addr string
	flag.StringVar(&addr, "iaddr", "tcp://localhost:5555", "input data port")

	flag.Parse()

	sck, err := zmq.NewSocket(zmq.PULL)
	if err != nil {
		log.Fatal(err)
	}
	defer sck.Close()

	log.Printf("dialing [%s]...\n", addr)
	err = sck.Connect(addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("dialing [%s]...\n", addr)

	for {
		msg, err := sck.RecvBytes(0)
		if err != nil {
			log.Fatalf("error recv: %v\n", err)
		}

		log.Printf("recv: %v\n", string(msg))
	}
}
