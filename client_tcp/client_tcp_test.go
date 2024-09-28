package client_tcp

import (
	"pad/broker_tcp"
	"pad/helpers"
	"sync"
	"testing"
)

func TestCreateTopic(t *testing.T) {
	port := "54321"

	b := broker_tcp.NewBrokerTCP()
	b.Init(port)
	go b.Serve()

	wg := &sync.WaitGroup{}
	helpers.Wait(wg,
		func() {
			c := NewClientTCP("Twister")
			c.Connect(port)

			go func() {
				for {
					c.ReadMessage()
				}

			}()

			c.NewTopic("dsa")
			helpers.Wait(wg, func() {
				c.Subscribe("dsa")
				c.Topics()
				c.Publish("Arent i a genius", "dsa")
			})

		},
		func() {
			c := NewClientTCP("TESTIES")
			c.Connect(port)

			go func() {
				for {
					c.ReadMessage()
				}
			}()

			helpers.Wait(wg, func() {
				c.Subscribe("dsa")
				c.Topics()
				c.Publish("I like being a genius", "dsa")
			})
		})
}

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
