package main

import (
	"bytes"
	"flag"
	"log"

	zmq "gopkg.in/zeromq/goczmq.v1"
)

func main() {
	var iaddr string
	flag.StringVar(&iaddr, "iaddr", "tcp://localhost:5555", "input data port")

	var oaddr string
	flag.StringVar(&oaddr, "oaddr", "tcp://localhost:5556", "output data port")

	flag.Parse()

	log.Printf("dialing [%s]...\n", iaddr)
	isck, err := zmq.NewPull(iaddr)
	if err != nil {
		log.Fatal(err)
	}
	defer isck.Destroy()
	log.Printf("dialing [%s]...\n", iaddr)

	log.Printf("dialing [%s]...\n", oaddr)
	osck, err := zmq.NewPush(oaddr)
	if err != nil {
		log.Fatal(err)
	}
	defer osck.Destroy()
	log.Printf("dialing [%s]...\n", oaddr)

	for {
		msg, err := isck.RecvMessage()
		if err != nil {
			log.Fatalf("error recv: %v\n", err)
		}

		log.Printf("recv: %v\n", string(msg[0]))
		omsg := bytes.Repeat(msg[0], 2)
		err = osck.SendFrame(omsg, zmq.FlagNone)
		if err != nil {
			log.Fatalf("error send: %v\n", err)
		}
	}
}
