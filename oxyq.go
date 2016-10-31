// Copyright 2016 The oxyq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package oxyq provides a basic framework to run FairMQ-like tasks.
package oxyq

import (
	"fmt"
	"sync"
)

// Socket is the main access handle that clients use to access the OxyQ system.
type Socket interface {
	// Close closes the open Socket
	Close() error

	// Send puts the message on the outbound send queue.
	// Send blocks until the message can be queued or the send deadline expires.
	Send(data []byte) error

	// Recv receives a complete message.
	Recv() ([]byte, error)

	// Listen connects alocal endpoint to the Socket.
	Listen(addr string) error

	// Dial connects a remote endpoint to the Socket.
	Dial(addr string) error
}

type Device struct {
	chans map[string][]Channel
}

type Channel struct {
	sck  Socket
	cmd  chan CmdType
	name string
}

type CmdType byte

type Type int

const (
	Invalid Type = iota
	Sub
	Pub
	XSub
	XPub
	Push
	Pull
	Req
	Rep
	Dealer
	Router
	Pair
	Bus
)

var drivers struct {
	sync.RWMutex
	db map[string]Driver
}

func Register(name string, drv Driver) {
	drivers.Lock()
	defer drivers.Unlock()
	if _, dup := drivers.db[name]; dup {
		panic(fmt.Errorf("oxyq: driver with name %q already registered", name))
	}
	drivers.db[name] = drv
}

func Open(name string) (Driver, error) {
	drivers.RLock()
	defer drivers.RUnlock()
	drv, ok := drivers.db[name]
	if !ok {
		return nil, fmt.Errorf("oxyq: no such driver %q", name)
	}
	return drv, nil
}

type Driver interface {
	NewSocket(typ Type) (Socket, error)
}

func init() {
	drivers.Lock()
	defer drivers.Unlock()
	drivers.db = make(map[string]Driver)
}
