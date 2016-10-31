// Copyright 2016 The oxyq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	word    = 4
	cdhSize = 10 * word
)

func main() {
	fname := os.Args[1]
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := NewReader(f)
	for i := 0; r.Scan(); i++ {
		if r.Err() != nil {
			break
		}
		fmt.Printf("==========================\n")
		fmt.Printf(">>> evt #%d...\n", i)
		r.hdr.print()
		err := r.ReadHeader()
		if err != nil {
			log.Fatal(err)
		}
		r.evt.print()
		r.eqp.print()
	}
	err = r.Err()
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		log.Fatal(err)
	}
}

func bread(r io.Reader, data interface{}) {
	err := binary.Read(r, binary.LittleEndian, data)
	if err != nil {
		log.Fatalf("error reading data: %v\n", err)
	}
}

func dump(data []byte) string {
	var o []string
	for i, v := range data {
		trail := ""
		if i%4 == 3 {
			trail = "\n"
		}
		o = append(o, fmt.Sprintf("%02x%s", v, trail))
	}
	return strings.Join(o, " ")
}
