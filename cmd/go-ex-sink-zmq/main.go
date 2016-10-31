// Copyright 2016 The oxyq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

	sck := gomq.NewPull(zmtp.NewSecurityNull())
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
