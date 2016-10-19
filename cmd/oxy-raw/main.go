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

var (
	endDDL = uint32(0xdeadface)
)

func main() {
	fname := os.Args[1]
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buf := make([]byte, cdhSize)
	var evthdr EventHeader
	bread(f, &evthdr)

	fmt.Printf("=== hdr ===\n %#v\n", evthdr)
	if evthdr.Magic != Magic {
		log.Fatalf("error magic number: got=0x%x want=0x%x\n", evthdr.Magic, Magic)
	}
	evthdr.print()

	var sub EventHeader
	bread(f, &sub)
	fmt.Printf("=== sub\n%#v\n", sub)
	sub.print()

	var equip EquipmentHeader
	bread(f, &equip)
	equip.print()

	var cdh CDH
	bread(f, &cdh)
	fmt.Printf("\n=== CDH ===\n%#v\n", cdh)

	bread(f, &equip)
	equip.print()

	buf = buf[:8*word]
	bread(f, buf)
	fmt.Printf("=== ddl-hdr ===\n %v\n", dump(buf))

	buf = buf[:word]
	bread(f, buf)
	fmt.Printf("=== blk-len ===\n %v\n", dump(buf))
	fmt.Printf(">>> %d\n", binary.LittleEndian.Uint32(buf[:]))

	buf = buf[:1*word]
	bread(f, buf)
	fmt.Printf("=== tot-len ===\n %v\n", dump(buf))
	fmt.Printf(">>> %d\n", binary.LittleEndian.Uint32(buf[:]))

	f.Seek(0, 0)
	testReader(f)
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

func testReader(f *os.File) {
	log.Printf("=== test-reader ===\n")
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
	err := r.Err()
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		log.Fatal(err)
	}
}
