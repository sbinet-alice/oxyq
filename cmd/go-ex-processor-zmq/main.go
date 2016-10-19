package main

import (
	"bytes"
	"flag"
	"log"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	var iaddr string
	flag.StringVar(&iaddr, "iaddr", "tcp://localhost:5555", "input data port")

	var oaddr string
	flag.StringVar(&oaddr, "oaddr", "tcp://localhost:5556", "output data port")

	flag.Parse()

	isck, err := zmq.NewSocket(zmq.PULL)
	if err != nil {
		log.Fatal(err)
	}
	defer isck.Close()

	log.Printf("dialing [%s]...\n", iaddr)
	err = isck.Connect(iaddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("dialing [%s]...\n", iaddr)

	osck, err := zmq.NewSocket(zmq.PUSH)
	if err != nil {
		log.Fatal(err)
	}
	defer osck.Close()

	log.Printf("dialing [%s]...\n", oaddr)
	err = osck.Connect(oaddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("dialing [%s]...\n", oaddr)

	for {
		msg, err := isck.RecvBytes(0)
		if err != nil {
			log.Fatalf("error recv: %v\n", err)
		}

		log.Printf("recv: %v\n", string(msg))
		omsg := bytes.Repeat(msg, 2)
		_, err = osck.SendBytes(omsg, 0)
		if err != nil {
			log.Fatalf("error send: %v\n", err)
		}
	}
}
