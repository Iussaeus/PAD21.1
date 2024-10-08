package client_tcp

import (
	"bufio"
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
	Name     string `json:"name"`
	conn     net.Conn
	messages []message_tcp.Message
	r        *bufio.Reader
	w        *bufio.Writer
	mutex    *sync.Mutex
	wg       *sync.WaitGroup
}

func NewClientTCP(name string) *ClientTCP {
	return &ClientTCP{
		Name:  name,
		mutex: &sync.Mutex{},
		wg:    &sync.WaitGroup{},
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
	if c.w != nil {
		err := c.w.Flush()
		if err != nil {
			log.Printf("Err while flushing while closing connection: %v\n", err)
		}
		c.w = nil
	}

	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			log.Printf("Err while closing connection: %v\n", err)
		}
		c.conn = nil
	}

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

	if err != nil {
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

	c.messages = append(c.messages, m)

	fmt.Printf("%s got message: %s\n", c.Name, m)
	return nil
}

// TODO: make it a cli, that  can be connected untill it receives a message
// TODO: or untill it sends one and gets receives a confirmation
// TODO: make a function to send message to another client for testing purposes
// with a time and topic as params

func InitNewClientFunc(name string, port string, wch chan struct{}, f func(c *ClientTCP)) {
	c := NewClientTCP(name)
	c.Connect(port)

	go func() {
		for {
			select {
			default:
				err := c.ReadMessage()
				if err == io.EOF && err != nil {
				}
			case <-wch:
				wch <- struct{}{}
				return
			}
		}
	}()

	helpers.Wait(c.wg, func() {
		f(c)
		wch <- struct{}{}
	})
}

func (c *ClientTCP) InitFunc(port string, wch chan struct{}, f func(c *ClientTCP)) {
	c.Connect(port)

	go func() {
		for {
			select {
			default:
				err := c.ReadMessage()
				if err == io.EOF && err != nil {
				}
			case <-wch:
				wch <- struct{}{}
				return
			}
		}
	}()

	helpers.Wait(c.wg, func() {
		f(c)
		wch <- struct{}{}
	})
}

func Run(name string, port string) {
	fmt.Println("Started Client main")

	ch := make(chan struct{})

	go InitNewClientFunc("SaneaNeLoh", "", ch,
		func(c *ClientTCP) {
			c.NewTopic("dsa")
			time.Sleep(300 * time.Millisecond)
			c.Subscribe("dsa")
			time.Sleep(500 * time.Millisecond)
			c.Topics()
			time.Sleep(600 * time.Millisecond)
			c.Publish("Hell", "dsa")
			time.Sleep(600 * time.Millisecond)
		})

	go func() {
		c := NewClientTCP("Test")
		c.InitFunc("", ch,
			func(c *ClientTCP) {
				time.Sleep(500 * time.Millisecond)
				c.Subscribe("dsa")
				time.Sleep(500 * time.Millisecond)
				c.Topics()
				time.Sleep(600 * time.Millisecond)
				c.Publish("1", "dsa")
				time.Sleep(600 * time.Millisecond)
			})
	}()

	go func() {
		c1 := NewClientTCP("Test2")
		c1.InitFunc("", ch,
			func(c *ClientTCP) {
				time.Sleep(500 * time.Millisecond)
				c.Subscribe("dsa")
				time.Sleep(500 * time.Millisecond)
				c.Topics()
				time.Sleep(600 * time.Millisecond)
				c.Publish("2", "dsa")
				time.Sleep(600 * time.Millisecond)
			})
	}()

	go func() {
		c2 := NewClientTCP("Test3")
		c2.InitFunc("", ch,
			func(c *ClientTCP) {
				time.Sleep(500 * time.Millisecond)
				c.Subscribe("dsa")
				time.Sleep(500 * time.Millisecond)
				c.Topics()
				time.Sleep(600 * time.Millisecond)
				c.Publish("3", "dsa")
				time.Sleep(600 * time.Millisecond)
			})
	}()

	go func() {
		c3 := NewClientTCP("Test4")
		c3.InitFunc("", ch,
			func(c *ClientTCP) {
				time.Sleep(500 * time.Millisecond)
				c.Subscribe("dsa")
				time.Sleep(500 * time.Millisecond)
				c.Topics()
				time.Sleep(600 * time.Millisecond)
				c.Publish("4", "dsa")
				time.Sleep(600 * time.Millisecond)
			})
	}()

	go func() {
		c4 := NewClientTCP("Test5")
		c4.InitFunc("", ch,
			func(c *ClientTCP) {
				time.Sleep(500 * time.Millisecond)
				c.Subscribe("dsa")
				time.Sleep(500 * time.Millisecond)
				c.Topics()
				time.Sleep(600 * time.Millisecond)
				c.Publish("5", "dsa")
				time.Sleep(600 * time.Millisecond)
			})
	}()

	<-ch
}
