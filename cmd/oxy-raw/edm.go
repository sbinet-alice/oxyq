package main

import (
	"fmt"
	"time"
)

const (
	EventIDBytes              = 8
	EventTriggerPatternBytes  = 16
	EventDetectorPatternBytes = 4
	AllAttributeWords         = 3

	Magic = uint32(0xda1e5afe)
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
	TypeAttr        TypeAttr
	LDC             int32
	GDC             int32
	TimestampSec    uint32
	TimestampUsec   uint32
}

func (hdr *EventHeader) Time() time.Time {
	return time.Unix(int64(hdr.TimestampSec), int64(hdr.TimestampUsec)*1000)
}

func (hdr *EventHeader) print() {
	fmt.Printf("evt.Size= %d\n", hdr.Size)
	fmt.Printf("evt.Head= %d\n", hdr.HeadSize)
	fmt.Printf("evt.Vers= 0x%x (maj=%d min=%d)\n", hdr.Version, (hdr.Version&0xffff0000)>>16, (hdr.Version & 0x0000ffff))
	fmt.Printf("evt.Type= %v\n", hdr.Type)
	fmt.Printf("time: %v\n", hdr.Time())
	fmt.Printf("evt.gdc= %d\n", int32(hdr.GDC))
	fmt.Printf("evt.ldc= %d\n", int32(hdr.LDC))
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

type TypeAttr [AllAttributeWords]uint32

const (
	AttrEventSwapped = 66
	AttrEventPaged   = 67
	AttrSuperEvent   = 68
	AttrOrbitBC      = 69
	AttrKeepPages    = 70
	AttrHLTDecision  = 71
)

func (ta TypeAttr) IsSuperEvent() bool {
	//#define TEST_SYSTEM_ATTRIBUTE(m,b)  (((m)[SYS_ATTR_2_W(b)] & ATTR_2_B(b)) != 0)
	//#define ATTR_2_B(b)     (1<<((b)&0x1f))

	return ta[2]&(1<<((AttrSuperEvent)&0x1f)) != 0
}

// type EventID [EventIDBytes >> 2]uint32
type EventID [2]uint32

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
	BlockLength         uint32
	Xfield0             uint32
	Xfield1             uint32
	Xfield2             uint32
	Xfield3             uint32
	EventTriggerPattern [4]uint32
	ROIHigh             uint32
}

type EquipmentHeader struct {
	Size          uint32
	Type          uint32
	ID            uint32
	TypeAttr      [AllAttributeWords]uint32
	BasicElemSize uint32
}

func (equip *EquipmentHeader) print() {
	fmt.Printf("=== equip ===\n")
	fmt.Printf("eq.Size= %d\n", equip.Size)
	fmt.Printf("eq.Type= %d\n", equip.Type)
	fmt.Printf("eq.ID=   %d\n", equip.ID)
	fmt.Printf("eq.Attr= %v\n", equip.TypeAttr)
	fmt.Printf("eq.ElemSize= %d\n", equip.BasicElemSize)
}

type EquipmentDescriptor struct {
	Header EquipmentHeader
	Vector EventVector
}

func (eq *EquipmentDescriptor) print() {
	eq.Header.print()
	fmt.Printf("eq.vector= %#v\n", eq.Vector)
}

type EventVector struct {
	BankID uint16
	Size   uint32
	Offset uint32
}

type Equipment struct {
	Header EquipmentHeader
	Raw    [1]uint16
}

type DataHeader struct {
	Size              uint32 // size of raw data in bytes
	Word              uint32 // bunch crossing, L1 trigger message and format version
	EventID           uint32 // orbit number
	SubDets           uint32 // block attributes and participating sub-detectors
	StatusMiniEventID uint32 // status & error bits + mini event ID
	TriggerClassLow   uint32 // low bits of trigger class
	ROILo             uint32 // low bits of ROI data
	ROIHi             uint32 // high bits of ROI data
}

type DataHeaderV3 struct {
	Size              uint32 // size of raw data in bytes
	Word              uint32 // bunch crossing, L1 trigger message and format version
	EventID           uint32 // orbit number
	SubDets           uint32 // block attributes and participating sub-detectors
	StatusMiniEventID uint32 // status & error bits + mini event ID
	Trigger           [4]uint32
	ROI               uint32
}
