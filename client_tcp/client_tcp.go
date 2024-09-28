package client_tcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"pad/helpers"
	"pad/message_tcp"
)

type ClientTCP struct {
	Name   string `json:"name"`
	conn   net.Conn
	r      *bufio.Reader
	w      *bufio.Writer
	mutex  *sync.Mutex
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewClientTCP(name string) *ClientTCP {
	ctx, cancel := context.WithCancel(context.Background())
	return &ClientTCP{
		Name:   name,
		mutex:  &sync.Mutex{},
		wg:     &sync.WaitGroup{},
		ctx:    ctx,
		cancel: cancel,
	}
}

func (c *ClientTCP) Connect(port string) {
	if port == "" {
		port = "12345"
	}
	conn, err := net.Dial("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to connect client: %s\n", err)
	}

	c.conn = conn
	c.r = bufio.NewReader(conn)
	c.w = bufio.NewWriter(conn)

	msg, err := c.r.ReadString('\n')
	if err != nil {
		log.Printf("Failed to read buffer ON CONNECTION: %v\n", err)
	} else {
		msg = msg[:len(msg)-1]
		helpers.CPrintf(helpers.Blue, "Got server greeting: %v\n", msg)
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		js, err := json.Marshal(c)
		if err != nil {
			log.Printf("Failed to marshall json ON CONNECTION: %v\n", err)
		}
		c.mutex.Lock()
		defer c.mutex.Unlock()

		_, err = c.w.Write(append(js, '\n'))
		if err != nil {
			log.Printf("Failed to write to buffer ON CONNECTION: %v\n", err)
		}

		err = c.w.Flush()
		if err != nil {
			log.Printf("Failed to flush buffer ON CONNECTION: %v\n", err)
		}
	}()
	c.wg.Wait()
}

func (c *ClientTCP) Disconnect() {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			log.Printf("Err while closing connection: %v\n", err)
		}
		c.conn = nil
	}

	c.cancel()

	if c.w != nil {
		err := c.w.Flush()
		if err != nil {
			log.Printf("Err while flushing while closing connection: %v\n", err)
		}
		c.w = nil
	}

	c.r = nil

	c = nil
}

func (c *ClientTCP) SendMessage(m message_tcp.Message) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.w == nil {
		return fmt.Errorf("Reader is nil")
	}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = c.w.Write(append(b, '\n'))
	if err != nil {
		return err
	}

	err = c.w.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientTCP) Publish(p string, t message_tcp.Topic) {
	m := message_tcp.Message{
		Sender:  c.Name,
		Cmd:     message_tcp.Publish,
		Topic:   t,
		Payload: p,
	}

	c.SendMessage(m)
}

func (c *ClientTCP) Subscribe(t message_tcp.Topic) {
	m := message_tcp.Message{
		Sender:  c.Name,
		Cmd:     message_tcp.Subscribe,
		Topic:   t,
		Payload: "",
	}

	c.SendMessage(m)
}
func (c *ClientTCP) Unsubscribe(t message_tcp.Topic) {
	m := message_tcp.Message{
		Sender:  c.Name,
		Cmd:     message_tcp.Subscribe,
		Topic:   t,
		Payload: "",
	}

	c.SendMessage(m)
}

func (c *ClientTCP) NewTopic(t message_tcp.Topic) {
	m := message_tcp.Message{
		Sender:  c.Name,
		Cmd:     message_tcp.NewTopic,
		Topic:   t,
		Payload: "",
	}

	c.SendMessage(m)
}

func (c *ClientTCP) DeleteTopic(t message_tcp.Topic) {
	m := message_tcp.Message{
		Sender:  c.Name,
		Cmd:     message_tcp.DeleteTopic,
		Topic:   t,
		Payload: "",
	}

	c.SendMessage(m)
}

func (c *ClientTCP) Topics() {
	m := message_tcp.Message{
		Sender:  c.Name,
		Cmd:     message_tcp.Topics,
		Topic:   message_tcp.Empty,
		Payload: "",
	}

	c.SendMessage(m)
}

func (c *ClientTCP) ReadMessage() error {
	if c.r == nil {
		return fmt.Errorf("Reader is nil\n")
	}

	s, err := c.r.ReadString('\n')

	if err == net.ErrClosed {
		return err
	} else if err != nil {
		return err
	}

	if s == "" {
		return fmt.Errorf("Received empty message\n")
	}

	m := message_tcp.Message{}
	err = json.Unmarshal([]byte(s), &m)
	if err != nil {
		return err
	}

	fmt.Printf("%s got message: %s\n", c.Name, m)
	return nil
}

// TODO: make it a cli, that  can be connected untill it receives a message
// TODO: or untill it sends one and gets receives a confirmation
// TODO: make a function to send message to another client for testing purposes
// with a time and topic as params

func Run(name string, port string) {
	helpers.AwaitSIGINT()

	fmt.Println("Started Client main")

	wg := &sync.WaitGroup{}
	time.Sleep(500 * time.Millisecond)
	helpers.Wait(wg,
		func() {
			c := NewClientTCP("Twister")
			c.Connect(port)

			go func() {
				for {
					err := c.ReadMessage()
					if err != io.EOF && err != nil {
						fmt.Println(err)
					}
				}

			}()

			helpers.Wait(wg, func() {
				c.NewTopic("dsa")
				time.Sleep(500 * time.Millisecond)
				c.Subscribe("dsa")
				time.Sleep(500 * time.Millisecond)
				c.Topics()
				time.Sleep(500 * time.Millisecond)
				c.Publish("Arent i a genius", "dsa")
			})

		},
		func() {
			c := NewClientTCP("TESTIES")
			c.Connect(port)

			go func() {
				for {
					err := c.ReadMessage()
					if err != io.EOF && err != nil {
						fmt.Println(err)
					}
				}
			}()

			helpers.Wait(wg, func() {
				time.Sleep(500 * time.Millisecond)
				c.Topics()
				time.Sleep(500 * time.Millisecond)
				c.Subscribe("dsa")
				c.Publish("I like being a genius", "dsa")
			})
		})
}
