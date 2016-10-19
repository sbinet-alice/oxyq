package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type reader struct {
	r   io.ReadSeeker
	ddl int // nb of bytes for DDL

	err error
	buf *bytes.Reader

	hdr EventHeader
	evt EventHeader
	sub EventHeader
	eqp EquipmentHeader
}

func NewReader(r io.ReadSeeker) *reader {
	return &reader{
		r:   r,
		buf: bytes.NewReader(nil),
	}
}

func (r *reader) Err() error {
	return r.err
}

func (r *reader) Scan() bool {
	if r.err != nil {
		return false
	}
	r.reset()

	pos, err := r.r.Seek(0, io.SeekCurrent)
	if err != nil {
		r.err = err
		return false
	}
	err = binary.Read(r.r, binary.LittleEndian, &r.hdr)
	if err != nil {
		r.err = err
		return false
	}
	_, err = r.r.Seek(pos, io.SeekStart)
	if err != nil {
		r.err = err
		return false
	}
	buf := make([]byte, r.hdr.Size)
	n, err := r.r.Read(buf)
	if err != nil {
		r.err = err
		return false
	}
	if n != int(r.hdr.Size) {
		r.err = io.ErrShortWrite
		return false
	}
	r.buf = bytes.NewReader(buf)
	err = binary.Read(r.buf, binary.LittleEndian, &r.evt)
	if err != nil {
		r.err = err
		return false
	}

	return true
}

// ReadHeader reads an event header at the current position
func (r *reader) ReadHeader() error {
	// check whether there are sub-events
	if r.evt.Size <= r.evt.HeadSize {
		fmt.Printf("evt-size <= evt-hdr-size (%v <= %v)\n", r.evt.Size, r.evt.HeadSize)
		return io.EOF
	}

	for {
		// skip payload if event not selected
		if r.ddl > 0 {
			_, err := r.buf.Seek(int64(r.ddl), io.SeekCurrent)
			if err != nil {
				r.err = err
				return err
			}
		}

		// read first or next sub-evt

		switch {
		case !r.evt.TypeAttr.IsSuperEvent():
			r.sub = r.evt
		default:
			err := binary.Read(r.buf, binary.LittleEndian, &r.sub)
			if err != nil {
				r.err = err
				return err
			}
		}
		// check magic word
		if r.sub.Magic != Magic {
			r.err = fmt.Errorf("invalid event magic number: %v\n", r.sub.Magic)
			return r.err
		}

		// continue if no data in sub-event
		if r.sub.Size == r.sub.HeadSize {
			continue
		}

		// get the first or next equipment
		err := binary.Read(r.buf, binary.LittleEndian, &r.eqp)
		if err != nil {
			r.err = err
			return err
		}

		if r.eqp.Size <= 0 {
			// go to next sub-event if no payload
			continue
		}

		// FIXME(sbinet) check require-header
		if false {
			var (
				sz   uint32
				vers uint32
			)
			pos, _ := r.buf.Seek(0, io.SeekCurrent)
			bread(r.buf, &sz)
			bread(r.buf, &vers)
			r.buf.Seek(pos, io.SeekStart)
			switch (vers & 0xffff0000) >> 16 {
			case 2:
				var ehdr DataHeader
				fmt.Printf("raw-data-header 2\n")
				bread(r.buf, &ehdr)
				fmt.Printf("%#v\n", ehdr)
			case 3:
				fmt.Printf("raw-data-header 3\n")
				var ehdr DataHeaderV3
				bread(r.buf, &ehdr)
				fmt.Printf("%#v\n", ehdr)
			default:
				log.Fatalf("invalid data-header version: %d\n", vers)
			}
		}

		break
	}
	return nil
}

func (r *reader) reset() {
	r.ddl = 0
	r.buf = nil
}
