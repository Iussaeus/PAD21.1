package client_tcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"pad/helpers"
	"pad/message_tcp"
	"sync"
	"time"
)

var _ message_tcp.Client = (*ClientTCP)(nil)

type ClientTCP struct {
	Name   string              `json:"name"`
	conn   net.Conn
	r      *bufio.Reader
	w      *bufio.Writer
	mutex  *sync.Mutex
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewClientTCP(name string, topics []message_tcp.Topic) *ClientTCP {
	ctx, cancel := context.WithCancel(context.Background())
	return &ClientTCP{
		Name:   name,
		mutex:  &sync.Mutex{},
		wg:     &sync.WaitGroup{},
		ctx:    ctx,
		cancel: cancel,
	}
}

func (c *ClientTCP) Connect() {
	conn, err := net.Dial("tcp", "localhost:12345")
	if err != nil {
		log.Printf("Failed to connect client: %s", err)
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
			log.Printf("Err while closing connection: %v", err)
		}
		c.conn = nil
	}

	c.cancel()

	if c.w != nil {
		err := c.w.Flush()
		if err != nil {
			log.Printf("Err while flushing while closing connection: %v", err)
		}
		c.w = nil
	}

	c.r = nil

	c = nil
}

func (c *ClientTCP) SendString(m string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, err := c.w.WriteString(m + "\n")
	if err != nil {
		log.Printf("Writer failed: %v", err)
		return
	}

	err = c.w.Flush()
	if err != nil {
		log.Printf("writer failed to flush: %v", err)
		return
	}
}

func (c *ClientTCP) SendMessage(m message_tcp.Message) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.w == nil {
		log.Println("Reader is nil")
		return
	}

	js, err := json.Marshal(m)
	if err != nil {
		log.Printf("Failed to marshall: %v", err)
	}

	_, err = c.w.Write(append(js, '\n'))
	if err != nil {
		log.Printf("Writer failed: %v", err)
		return
	}

	err = c.w.Flush()
	if err != nil {
		log.Printf("writer failed to flush: %v", err)
		return
	}
}

func (c *ClientTCP) ReadMessage() {
	// fmt.Println("Reading")
	if c.r == nil {
		log.Println("Reader is nil")
		return
	}
	c.wg.Add(1)
	defer c.wg.Done()

	s, err := c.r.ReadString('\n')
	if err != nil {
		log.Printf("Reader failed: %v", err)
		return
	}

	if s != "" {
		fmt.Println("Got message", s)
	}
	// fmt.Println("Done Reading")
}

func Run(name string) {
	fmt.Println("Started Client from main")
	top := []message_tcp.Topic{"a", "b", "c"}
	c := NewClientTCP(name, top)
	c.Connect()

	receiver := ""

	if name == "Test1" {
		receiver = "Test2"
	} else {
		receiver = "Test1"
	}

	t := time.After(3 * time.Second)

	msg := message_tcp.Message{
		Sender:   c.Name,
		Type:     message_tcp.Publication,
		Topic:    top[0],
		Cmd:      message_tcp.Publish,
		Receiver: receiver,
		Payload:  "I hate you",
	}

	go func() {
		for {
			c.ReadMessage()
			select {
			case <-c.ctx.Done():
				return
			default:
				c.wg.Wait()
			}
		}
	}()

	go c.SendMessage(msg)
cycle:
	for {

		select {
		case <-t:
			helpers.CPrintf(helpers.Red, "End of execution")
			c.Disconnect()
			c.wg.Wait()
			break cycle
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}

	fmt.Println("end of client from main")
}
