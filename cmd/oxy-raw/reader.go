package main

import (
	"encoding/binary"
	"io"
)

type reader struct {
	r   io.ReadSeeker
	buf []byte
	cur int

	hdr *EventHeader
}

func NewReader(r io.ReadSeeker) *reader {
	return &reader{
		r:   r,
		cur: 0,
	}
}

func (r *reader) ReadEvent() error {
	return binary.Read(r.r, binary.LittleEndian, r.hdr)
}
