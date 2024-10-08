package client_tcp

import (
	// "io"
	"pad/broker_tcp"
	"pad/helpers"
	// "sync"
	"testing"
	// "time"
)

// func TestCreateTopic(t *testing.T) {
// 	port := "54321"
//
// 	b := broker_tcp.NewBrokerTCP()
// 	b.Init(port)
// 	timer := time.NewTimer(3 * time.Second)
// 	go func(){
// 		b.Serve()
// 	}()
//
// 	wg := &sync.WaitGroup{}
// 	helpers.Wait(wg,
// 		func() {
// 			c := NewClientTCP("Twister")
// 			c.Connect(port)
//
// 			go func() {
// 				for {
// 					if err := c.ReadMessage(); err == io.EOF {
// 						return
// 					}
// 				}
//
// 			}()
//
// 			c.NewTopic("dsa")
// 			helpers.Wait(wg, func() {
// 				c.NewTopic("dsa")
// 				time.Sleep(500 * time.Millisecond)
// 				c.Subscribe("dsa")
// 				time.Sleep(500 * time.Millisecond)
// 				c.Topics()
// 				time.Sleep(500 * time.Millisecond)
// 				c.Publish("Arent i a genius", "dsa")
// 			})
//
// 		},
// 		func() {
// 			c := NewClientTCP("TESTIES")
// 			c.Connect(port)
//
// 			go func() {
// 				for {
// 					if err := c.ReadMessage(); err == io.EOF {
// 						return
// 					}
// 				}
// 			}()
//
// 			helpers.Wait(wg, func() {
// 				time.Sleep(500 * time.Millisecond)
// 				c.Topics()
// 				time.Sleep(500 * time.Millisecond)
// 				c.Subscribe("dsa")
// 				c.Publish("I like being a genius", "dsa")
// 			})
// 		})
// 		select {
// 		case <-timer.C:
// 			return
// 		}
// }

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
