// Copyright 2016 The oxyq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oxyq_test

import (
	"fmt"
	"testing"

	"github.com/sbinet-alice/oxyq"
	_ "github.com/sbinet-alice/oxyq/nanomsg"
	_ "github.com/sbinet-alice/oxyq/zeromq"
)

func TestPushPullNN(t *testing.T) {
	const (
		N    = 5
		tmpl = "data-%02d"
	)

	drv, err := oxyq.Open("nanomsg")
	if err != nil {
		t.Fatal(err)
	}
	pull, err := drv.NewSocket(oxyq.Pull)
	if err != nil {
		t.Fatal(err)
	}
	defer pull.Close()

	push, err := drv.NewSocket(oxyq.Push)
	if err != nil {
		t.Fatal(err)
	}
	defer push.Close()

	go func() {
		err := push.Dial("tcp://localhost:5555")
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < N; i++ {
			err = push.Send([]byte(fmt.Sprintf(tmpl, i)))
			if err != nil {
				t.Fatalf("error sending data[%d]: %v\n", i, err)
			}
		}
		err = push.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	err = pull.Listen("tcp://*:5555")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < N; i++ {
		msg, err := pull.Recv()
		if err != nil {
			t.Fatal(err)
		}
		if got, want := string(msg), fmt.Sprintf(tmpl, i); got != want {
			t.Errorf("push-pull[%d]: got=%q want=%q\n", i, got, want)
		}
	}
	err = pull.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPushPullZMQ(t *testing.T) {
	const (
		N    = 5
		tmpl = "data-%02d"
	)

	drv, err := oxyq.Open("zeromq")
	if err != nil {
		t.Fatal(err)
	}
	pull, err := drv.NewSocket(oxyq.Pull)
	if err != nil {
		t.Fatal(err)
	}
	defer pull.Close()

	push, err := drv.NewSocket(oxyq.Push)
	if err != nil {
		t.Fatal(err)
	}
	defer push.Close()

	go func() {
		err := push.Dial("tcp://localhost:5555")
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < N; i++ {
			err = push.Send([]byte(fmt.Sprintf(tmpl, i)))
			if err != nil {
				t.Fatalf("error sending data[%d]: %v\n", i, err)
			}
		}
		err = push.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	err = pull.Listen("tcp://*:5555")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < N; i++ {
		msg, err := pull.Recv()
		if err != nil {
			t.Fatal(err)
		}
		if got, want := string(msg), fmt.Sprintf(tmpl, i); got != want {
			t.Errorf("push-pull[%d]: got=%q want=%q\n", i, got, want)
		}
	}
}
