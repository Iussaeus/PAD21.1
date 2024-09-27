package client_tcp

import (
	"pad/broker_tcp"
	"pad/helpers"
	"testing"
)

func TestStart(t *testing.T) {
	c := NewClientTCP("Test1")
	b := broker_tcp.NewBrokerTCP()

	p := "54321"

	b.Init(p)

	go b.Serve()
	c.Connect(p)
	r := []string{}

	helpers.Wait(c.wg,
		func() {
			if err := c.ReadMessage(); err != nil {
				r = append(r, err.Error())
			}
		},
		func() {
			if err := c.ReadMessage(); err != nil {
				r = append(r, err.Error())
			}
		},
		func() {
			c.Disconnect()
		},
		func() {
			if err := c.ReadMessage(); err != nil {
				r = append(r, err.Error())
			}
		},
		func() {
			if err := c.ReadMessage(); err != nil {
				r = append(r, err.Error())
			}
		},
		func() {
			if err := c.ReadMessage(); err != nil {
				r = append(r, err.Error())
			}
		},
	)

	helpers.Assert(len(r) != 0, "client reading from closed conn", t)
}
