package main

import (
	"flag"
	"log"

	zmq "gopkg.in/zeromq/goczmq.v1"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "tcp://localhost:5555", "output data port")

	flag.Parse()

	log.Printf("dialing [%s]...\n", addr)
	sck, err := zmq.NewPull(addr)
	if err != nil {
		log.Fatal(err)
	}
	defer sck.Destroy()
	log.Printf("dialing [%s]...\n", addr)

	for {
		msg := []byte("HELLO")
		err = sck.SendFrame(msg, zmq.FlagNone)
		if err != nil {
			log.Fatalf("error send: %v\n", err)
		}
	}
}
