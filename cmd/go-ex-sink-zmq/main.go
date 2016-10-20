package main

import (
	"flag"
	"log"

	zmq "gopkg.in/zeromq/goczmq.v1"
)

func main() {
	var addr string
	flag.StringVar(&addr, "iaddr", "tcp://localhost:5555", "input data port")

	flag.Parse()

	log.Printf("dialing [%s]...\n", addr)
	sck, err := zmq.NewPull(addr)
	if err != nil {
		log.Fatal(err)
	}
	defer sck.Destroy()
	log.Printf("dialing [%s]...\n", addr)

	for {
		msg, err := sck.RecvMessage()
		if err != nil {
			log.Fatalf("error recv: %v\n", err)
		}

		log.Printf("recv: %v\n", string(msg[0]))
	}
}
