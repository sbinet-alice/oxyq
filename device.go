// Copyright 2016 The oxyq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oxyq

import (
	"fmt"
	"log"
	"strings"

	"github.com/sbinet-alice/oxyq/config"
	"github.com/sbinet-alice/oxyq/mq"
	_ "github.com/sbinet-alice/oxyq/mq/nanomsg"
	_ "github.com/sbinet-alice/oxyq/mq/zeromq"
)

type Channel struct {
	cfg config.Channel
	sck mq.Socket
	cmd chan CmdType
	msg chan Msg
}

func (ch *Channel) Name() string {
	return ch.cfg.Name
}

func (ch *Channel) Send(data []byte) (int, error) {
	err := ch.sck.Send(data)
	return len(data), err
}

func (ch *Channel) Recv() ([]byte, error) {
	return ch.sck.Recv()
}

func (ch *Channel) run() {
	for {
		select {
		case msg := <-ch.msg:
			_, err := ch.Send(msg.Data)
			if err != nil {
				log.Fatal(err)
			}
		case ch.msg <- ch.recv():
		case cmd := <-ch.cmd:
			switch cmd {
			case CmdEnd:
				return
			}
		}
	}
}

func (ch *Channel) recv() Msg {
	data, err := ch.Recv()
	return Msg{
		Data: data,
		Err:  err,
	}
}

func newChannel(drv mq.Driver, cfg config.Channel) (Channel, error) {
	ch := Channel{
		cmd: make(chan CmdType),
		cfg: cfg,
	}
	// FIXME(sbinet) support multiple sockets to send/recv to/from
	if len(cfg.Sockets) != 1 {
		panic("oxyq: not implemented")
	}
	sck, err := drv.NewSocket(socketType(cfg.Sockets[0].Type))
	if err != nil {
		return ch, err
	}
	ch.sck = sck
	return ch, nil
}

type device struct {
	name  string
	chans map[string][]Channel
	cmds  chan CmdType
	msgs  map[msgAddr]chan Msg
}

func newDevice(drv mq.Driver, cfg config.Device) (*device, error) {
	log.Printf("--- new device: %v\n", cfg)
	dev := device{
		chans: make(map[string][]Channel),
		cmds:  make(chan CmdType),
		msgs:  make(map[msgAddr]chan Msg),
	}

	for _, opt := range cfg.Channels {
		log.Printf("--- new channel: %v\n", opt)
		ch, err := newChannel(drv, opt)
		if err != nil {
			return nil, err
		}
		ch.msg = make(chan Msg)
		dev.chans[opt.Name] = []Channel{ch}
		dev.msgs[msgAddr{name: opt.Name, id: 0}] = ch.msg
	}
	return &dev, nil
}

func (dev *device) Chan(name string, i int) (chan Msg, error) {
	msg, ok := dev.msgs[msgAddr{name, i}]
	if !ok {
		return nil, fmt.Errorf("oxyq: no such channel (name=%q index=%d)", name, i)
	}
	return msg, nil
}

func (dev *device) Done() chan CmdType {
	return nil
}

func (dev *device) isControler() {}

func (dev *device) run() {
	for n, chans := range dev.chans {
		log.Printf("--- init channels [%s]...\n", n)
		for i, ch := range chans {
			log.Printf("--- init channel[%s][%d]...\n", n, i)
			sck := ch.cfg.Sockets[0]
			switch strings.ToLower(sck.Method) {
			case "bind":
				go func() {
					err := ch.sck.Listen(sck.Address)
					if err != nil {
						log.Fatal(err)
					}
				}()
			case "connect":
				go func() {
					err := ch.sck.Dial(sck.Address)
					if err != nil {
						log.Fatal(err)
					}
				}()
			default:
				log.Fatalf("oxyq: invalid socket method (value=%q)", sck.Method)
			}
		}
	}

	for n, chans := range dev.chans {
		log.Printf("--- start channels [%s]...\n", n)
		for i := range chans {
			go chans[i].run()
		}
	}

}

type Device interface {
	Configure(cfg config.Device) error
	Init(ctrl Controler) error
	Run(ctrl Controler) error
	Pause(ctrl Controler) error
	Reset(ctrl Controler) error
}

type Controler interface {
	Chan(name string, i int) (chan Msg, error)
	Done() chan CmdType

	isControler()
}

type msgAddr struct {
	name string
	id   int
}

type Msg struct {
	Data []byte
	Err  error
}

func Main(dev Device) error {
	cfg, err := config.Parse()
	if err != nil {
		return err
	}

	drvName := cfg.Transport
	if drvName == "" {
		drvName = "zeromq"
	}

	drv, err := mq.Open(drvName)
	if err != nil {
		return err
	}

	devName := cfg.ID
	devCfg, ok := cfg.Options.Device(devName)
	if !ok {
		return fmt.Errorf("oxyq: no such device %q", devName)
	}

	sys, err := newDevice(drv, devCfg)
	if err != nil {
		return err
	}

	err = dev.Configure(devCfg)
	if err != nil {
		return err
	}

	go sys.run()

	err = dev.Init(sys)
	if err != nil {
		return err
	}

	err = dev.Run(sys)
	if err != nil {
		return err
	}

	return nil
}

func deviceConfig(name string, cfg config.Config) (config.Device, error) {
	var (
		dev config.Device
		err error
	)

	return dev, err
}
