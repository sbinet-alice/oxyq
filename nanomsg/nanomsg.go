// Copyright 2016 The oxyq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nanomsg

import (
	"fmt"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/bus"
	"github.com/go-mangos/mangos/protocol/pair"
	"github.com/go-mangos/mangos/protocol/pub"
	"github.com/go-mangos/mangos/protocol/pull"
	"github.com/go-mangos/mangos/protocol/push"
	"github.com/go-mangos/mangos/protocol/rep"
	"github.com/go-mangos/mangos/protocol/req"
	"github.com/go-mangos/mangos/protocol/sub"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/sbinet-alice/oxyq"
)

type driver struct{}

func (driver) NewSocket(typ oxyq.Type) (oxyq.Socket, error) {
	var sck mangos.Socket
	var err error

	switch typ {
	case oxyq.Sub, oxyq.XSub:
		sck, err = sub.NewSocket()
	case oxyq.Pub, oxyq.XPub:
		sck, err = pub.NewSocket()
	case oxyq.Push:
		sck, err = push.NewSocket()
	case oxyq.Pull:
		sck, err = pull.NewSocket()
	case oxyq.Req, oxyq.Dealer:
		sck, err = req.NewSocket()
	case oxyq.Rep, oxyq.Router:
		sck, err = rep.NewSocket()
	case oxyq.Pair:
		sck, err = pair.NewSocket()
	case oxyq.Bus:
		sck, err = bus.NewSocket()
	default:
		return nil, fmt.Errorf("oxyq/nanomsg: invalid socket type %v (%d)", typ, int(typ))
	}

	if err != nil {
		return nil, err
	}

	sck.AddTransport(ipc.NewTransport())
	sck.AddTransport(tcp.NewTransport())
	return sck, err
}

func init() {
	oxyq.Register("nanomsg", driver{})
}
