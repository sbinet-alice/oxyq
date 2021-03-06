// Copyright 2016 The oxyq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mq provides interfaces for message-queue sockets.
package mq

import (
	"fmt"
	"strings"
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

// SocketType describes the type of a socket (PUB, SUP, PUSH, PULL, ...)
type SocketType int

const (
	Invalid SocketType = iota
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

func SocketTypeFrom(name string) SocketType {
	switch strings.ToLower(name) {
	case "sub":
		return Sub
	case "pub":
		return Pub
	case "xpub":
		return XPub
	case "xsub":
		return XSub
	case "push":
		return Push
	case "pull":
		return Pull
	case "req":
		return Req
	case "rep":
		return Rep
	case "dealer":
		return Dealer
	case "router":
		return Router
	case "pair":
		return Pair
	case "bus":
		return Bus
	}
	panic(fmt.Errorf("oxyq: invalid socket type name (value=%q)", name))
}

var drivers struct {
	sync.RWMutex
	db map[string]Driver
}

// Register registers a new OxyQ driver plugin
func Register(name string, drv Driver) {
	drivers.Lock()
	defer drivers.Unlock()
	if _, dup := drivers.db[name]; dup {
		panic(fmt.Errorf("oxyq: driver with name %q already registered", name))
	}
	drivers.db[name] = drv
}

// Open returns a previously registered driver plugin
//
// e.g.
//  zmq, err := oxyq.Open("zeromq")
//  nn,  err := oxyq.Open("nanomsg")
func Open(name string) (Driver, error) {
	drivers.RLock()
	defer drivers.RUnlock()
	drv, ok := drivers.db[name]
	if !ok {
		return nil, fmt.Errorf("oxyq: no such driver %q", name)
	}
	return drv, nil
}

// Driver is an OxyQ plugin to create FairMQ-compatible message queue communications
type Driver interface {
	NewSocket(typ SocketType) (Socket, error)
	Name() string
}

func init() {
	drivers.Lock()
	defer drivers.Unlock()
	drivers.db = make(map[string]Driver)
}
