// Copyright 2016 The oxyq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package oxyq provides a basic framework to run FairMQ-like tasks.
package oxyq

import (
	"fmt"
	"strings"

	"github.com/sbinet-alice/oxyq/mq"
)

// CmdType describes commands to be sent to a device, via a channel.
type CmdType byte

const (
	CmdInitDevice CmdType = iota
	CmdInitTask
	CmdRun
	CmdPause
	CmdStop
	CmdResetTask
	CmdResetDevice
	CmdEnd
	CmdError
)

func (cmd CmdType) String() string {
	switch cmd {
	case CmdInitDevice:
		return "INIT_DEVICE"
	case CmdInitTask:
		return "INIT_TASK"
	case CmdRun:
		return "RUN"
	case CmdPause:
		return "PAUSE"
	case CmdStop:
		return "STOP"
	case CmdResetTask:
		return "RESET_TASK"
	case CmdResetDevice:
		return "RESET_DEVICE"
	case CmdEnd:
		return "END"
	case CmdError:
		return "ERROR_FOUND"
	}
	panic(fmt.Errorf("oxyq: invalid CmdType value (command=%d)", int(cmd)))
}

func socketType(name string) mq.SocketType {
	switch strings.ToLower(name) {
	case "sub":
		return mq.Sub
	case "pub":
		return mq.Pub
	case "xpub":
		return mq.XPub
	case "xsub":
		return mq.XSub
	case "push":
		return mq.Push
	case "pull":
		return mq.Pull
	case "req":
		return mq.Req
	case "rep":
		return mq.Rep
	case "dealer":
		return mq.Dealer
	case "router":
		return mq.Router
	case "pair":
		return mq.Pair
	case "bus":
		return mq.Bus
	}
	panic(fmt.Errorf("oxyq: invalid socket type name (value=%q)", name))
}
