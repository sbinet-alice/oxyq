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
	cfg    config.Device
	idatac chan oxyq.Msg
	odatac chan oxyq.Msg
}

func (dev *Device) Configure(cfg config.Device) error {
	dev.cfg = cfg
	return nil
}

func (dev *Device) Init(ctrl oxyq.Controler) error {
	idatac, err := ctrl.Chan("data1", 0)
	if err != nil {
		return err
	}

	odatac, err := ctrl.Chan("data2", 0)
	if err != nil {
		return err
	}

	dev.idatac = idatac
	dev.odatac = odatac
	return nil
}

func (dev *Device) Run(ctrl oxyq.Controler) error {
	for {
		select {
		case data := <-dev.idatac:
			log.Printf("received: %q\n", string(data.Data))
			out := append([]byte(nil), data.Data...)
			out = append(out, []byte(" (modified by "+dev.cfg.ID+")")...)
			dev.odatac <- oxyq.Msg{Data: out}
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
