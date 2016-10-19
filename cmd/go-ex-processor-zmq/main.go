package main

import (
	"bytes"
	"flag"
	"log"

	"github.com/zeromq/gomq"
	"github.com/zeromq/gomq/zmtp"
)

func main() {
	var iaddr string
	flag.StringVar(&iaddr, "iaddr", "tcp://localhost:5555", "input data port")

	var oaddr string
	flag.StringVar(&oaddr, "oaddr", "tcp://localhost:5556", "output data port")

	flag.Parse()

	isck := gomq.NewClient(zmtp.NewSecurityNull())
	log.Printf("dialing [%s]...\n", iaddr)
	err := isck.Connect(iaddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("dialing [%s]...\n", iaddr)

	osck := gomq.NewServer(zmtp.NewSecurityNull())
	log.Printf("dialing [%s]...\n", oaddr)
	_, err = osck.Bind(oaddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("dialing [%s]...\n", oaddr)

	for {
		msg, err := isck.Recv()
		if err != nil {
			log.Fatalf("error recv: %v\n", err)
		}

		log.Printf("recv: %v\n", string(msg))
		omsg := bytes.Repeat(msg, 2)
		err = osck.Send(omsg)
		if err != nil {
			log.Fatalf("error send: %v\n", err)
		}
	}
}
