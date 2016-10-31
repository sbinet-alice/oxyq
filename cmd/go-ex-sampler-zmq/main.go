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
	flag.StringVar(&addr, "addr", "tcp://localhost:5555", "output data port")

	flag.Parse()

	sck := gomq.NewPush(zmtp.NewSecurityNull())
	defer sck.Close()

	log.Printf("dialing [%s]...\n", addr)
	_, err := sck.Bind(addr)
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
