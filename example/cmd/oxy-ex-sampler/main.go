// Copyright 2016 The oxyq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/sbinet-alice/oxyq"
	"github.com/sbinet-alice/oxyq/config"
	_ "github.com/sbinet-alice/oxyq/zeromq"
)

type Device struct {
	cfg   config.Device
	datac chan oxyq.Msg
}

func (dev *Device) Configure(cfg config.Device) error {
	dev.cfg = cfg
	return nil
}

func (dev *Device) Init(ctrl oxyq.Controler) error {
	datac, err := ctrl.Chan("data1", 0)
	if err != nil {
		return err
	}

	dev.datac = datac
	return nil
}

func (dev *Device) Run(ctrl oxyq.Controler) error {
	for {
		select {
		case dev.datac <- oxyq.Msg{Data: []byte("HELLO")}:
		case <-ctrl.Done():
			return nil
		}
	}
	return nil
}

func (dev *Device) Pause(ctrl oxyq.Controler) error {
	return nil
}

func (dev *Device) Reset(ctrl oxyq.Controler) error {
	return nil
}

func main() {
	err := oxyq.Main(&Device{})
	if err != nil {
		log.Fatal(err)
	}
}
