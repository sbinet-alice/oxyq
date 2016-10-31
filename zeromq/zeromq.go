// Copyright 2016 The oxyq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zeromq

import (
	"fmt"

	"github.com/sbinet-alice/oxyq"
	"github.com/zeromq/gomq"
	"github.com/zeromq/gomq/zmtp"
)

type socket struct {
	sck gomq.ZeroMQSocket
}

func (s *socket) Close() error {
	s.sck.Close()
	return nil
}

func (s *socket) Send(data []byte) error {
	return s.sck.Send(data)
}

func (s *socket) Recv() ([]byte, error) {
	return s.sck.Recv()
}

func (s *socket) Listen(addr string) error {
	_, err := gomq.BindServer(s.sck.(gomq.Server), addr)
	return err
}

func (s *socket) Dial(addr string) error {
	return gomq.ConnectClient(s.sck.(gomq.Client), addr)
}

type driver struct{}

func (driver) NewSocket(typ oxyq.Type) (oxyq.Socket, error) {
	var (
		sck gomq.ZeroMQSocket
		err error
		m   = zmtp.NewSecurityNull()
	)

	switch typ {
	case oxyq.Sub, oxyq.XSub:
		panic("oxyq/zeromq: oxyq.Sub not yet implemented")

	case oxyq.Pub, oxyq.XPub:
		panic("oxyq/zeromq: oxyq.Pub not yet implemented")

	case oxyq.Push:
		sck = gomq.NewPush(m)

	case oxyq.Pull:
		sck = gomq.NewPull(m)

	case oxyq.Req, oxyq.Dealer:
		sck = gomq.NewClient(m)

	case oxyq.Rep, oxyq.Router:
		sck = gomq.NewServer(m)

	case oxyq.Pair:
		panic("oxyq/zeromq: oxyq.Pair not yet implemented")

	case oxyq.Bus:
		panic("oxyq/zeromq: oxyq.Bus not yet implemented")

	default:
		return nil, fmt.Errorf("oxyq/nanomsg: invalid socket type %v (%d)", typ, int(typ))
	}

	if err != nil {
		return nil, err
	}

	return &socket{sck: sck}, err
}

func init() {
	oxyq.Register("zeromq", driver{})
}
