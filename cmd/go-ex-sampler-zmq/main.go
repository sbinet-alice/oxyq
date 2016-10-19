package main

import (
	"flag"
	"log"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "tcp://*:5555", "output data port")

	flag.Parse()

	sck, err := zmq.NewSocket(zmq.PUSH)
	if err != nil {
		log.Fatal(err)
	}
	defer sck.Close()

	log.Printf("dialing [%s]...\n", addr)
	err = sck.Bind(addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("dialing [%s]... [done]\n", addr)

	for {
		msg := []byte("HELLO")
		_, err = sck.SendBytes(msg, 0)
		if err != nil {
			log.Fatalf("error send: %v\n", err)
		}
	}
}
