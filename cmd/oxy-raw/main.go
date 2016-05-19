package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

const (
	word    = 4
	cdhSize = 4 * 10
)

var (
	endDDL = uint32(0xdeadface)
	Magic  = uint32(0xda1e5afe)
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
	fmt.Printf("evt.Size= %d\n", evthdr.Size)
	fmt.Printf("evt.Head= %d\n", evthdr.HeadSize)
	fmt.Printf("evt.Vers= 0x%x (maj=%d min=%d)\n", evthdr.Version, (evthdr.Version&0xffff0000)>>16, (evthdr.Version & 0x0000ffff))
	fmt.Printf("evt.Type= %v\n", evthdr.Type)
	fmt.Printf("evt.Run=  %d\n", evthdr.Run)
	fmt.Printf("evt.ID=   %d (xid=%d, orbit=%d, period=%d)\n", evthdr.ID, evthdr.ID.BunchCrossing(), evthdr.ID.Orbit(), evthdr.ID.Period())

	fmt.Printf("time: %v %v => %v\n", int64(evthdr.TimestampSec), int64(evthdr.TimestampUsec),
		time.Unix(int64(evthdr.TimestampSec), int64(evthdr.TimestampUsec)*1000),
	)

	fmt.Printf("evt.gdc= %d\n", int32(evthdr.GDC))
	fmt.Printf("evt.ldc= %d\n", int32(evthdr.LDC))

	var sub EventHeader
	bread(f, &sub)
	fmt.Printf("=== sub\n%#v\n", sub)

	buf = buf[:7*word]
	bread(f, buf)
	fmt.Printf("=== equ ====\n %v\n", dump(buf))

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
}

func bread(r io.Reader, data interface{}) {
	err := binary.Read(r, binary.LittleEndian, data)
	if err != nil {
		log.Fatalf("error reading data: %v\n", err)
	}
}

const (
	EventIDBytes              = 8
	EventTriggerPatternBytes  = 16
	EventDetectorPatternBytes = 4
	AllAttributeWords         = 3
)

type Event struct {
	Header EventHeader
	Raw    [1]uint16
}

type EventHeader struct {
	Size            uint32
	Magic           uint32
	HeadSize        uint32
	Version         uint32
	Type            EventType
	Run             uint32
	ID              EventID
	TriggerPattern  [EventTriggerPatternBytes >> 2]uint32
	DetectorPattern [EventDetectorPatternBytes >> 2]uint32
	TypeAttr        [AllAttributeWords]uint32
	LDC             uint32
	GDC             uint32
	TimestampSec    uint32
	TimestampUsec   uint32
}

type EventType uint32

//go:generate stringer -type EventType
const (
	StartOfRun                   EventType = 1
	EndOfRun                     EventType = 2
	StartOfRunFiles              EventType = 3
	EndOfRunFiles                EventType = 4
	StartOfBurst                 EventType = 5
	EndOfBurst                   EventType = 6
	PhysicsEvent                 EventType = 7
	CalibrationEvent             EventType = 8
	EventFormatError             EventType = 9
	StartOfData                  EventType = 10
	EndOfData                    EventType = 11
	SystemSoftwareTriggerEvent   EventType = 12
	DetectorSoftwareTriggerEvent EventType = 13
	SyncEvent                    EventType = 14
)

type EventID [EventIDBytes >> 2]uint32

func (id EventID) BunchCrossing() uint32 {
	return id[1] & 0x00000fff
}

func (id EventID) Orbit() uint32 {
	return ((id[0] << 20) & 0xf00000) | ((id[1] >> 12) & 0xfffff)
}

func (id EventID) Period() uint32 {
	return (id[0] >> 4) & 0x0fffffff
}

type CDH struct {
	BlockLength      uint32
	Version          uint8
	L1TriggerMessage [10]byte
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
