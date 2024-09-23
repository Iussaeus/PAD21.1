package clientTcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"pad/helpers"
	"pad/message_tcp"
	"sync"
	"time"
)

// TODO: use bufio properly
// TODO: send the marshallized data

var _  Client = (*ClientTCP)(nil)

type ClientTCP struct {
	Name   string
	Topics []message_tcp.Topic
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	mutex  *sync.Mutex
	wg     *sync.WaitGroup
}

func (c *ClientTCP) Connect() {
	conn, err := net.Dial("tcp", "localhost:12345")
	if err != nil {
		log.Printf("Failed to connect client: %s", err)
	}

	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)
	c.mutex = &sync.Mutex{}
	c.wg = &sync.WaitGroup{}

	if err != nil {
		log.Printf("Failed to marshall json ON CONNECTION: %v\n", err)
	}

	msg, err := c.reader.ReadString('\n')
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
		_, err = c.writer.WriteString(string(js))
		if err != nil {
			log.Printf("Failed to write to buffer ON CONNECTION: %v\n", err)
		}

		err = c.writer.Flush()
		if err != nil {
			log.Printf("Failed to flush buffer ON CONNECTION: %v\n", err)
		}
	}()
	c.wg.Wait()
}

func (c *ClientTCP) Disconnect() {
	c.conn.Close()
}

func (c *ClientTCP) SendMessage(receiver string, msg string) {
	response := fmt.Sprintf("%s ,from: %s\n", msg, c.Name)

	_, err := c.writer.WriteString(response)
	if err != nil {
		log.Printf("Writer failed: %v", err)
		return
	}

	err = c.writer.Flush()
	if err != nil {
		log.Printf("writer failed to flush: %v", err)
		return
	}
}

func (c *ClientTCP) ReadMessage() {
	// fmt.Println("Reading")
	s, err := c.reader.ReadString('\n')
	if err != nil {
		log.Printf("Reader failed: %v", err)
	}

	s = s[:len(s)-1]

	if s != "" {
		fmt.Println("Got message", s)
	} else {
		fmt.Println("Don't got message,")
	}

	// fmt.Println("Done Reading")
}

func Run() {
	fmt.Println("Started Client from main")
	top := []message_tcp.Topic{"a", "b", "c"}
	c := ClientTCP{
		Name: "Test",
		Topics: top,
	}

	c.Connect()

	t := time.After(10 * time.Second)
cycle:
	for {
		go c.SendMessage("", "Get the message, please")
		time.Sleep(time.Millisecond * 500)

		go func() {
			for {
				time.Sleep(time.Millisecond * 500)
				c.ReadMessage()
			}
		}()
		select {
		case <-t:
			helpers.CPrintf(helpers.Red, "End of execution")
			break cycle
		default:
			continue
		}
	}

	fmt.Println("end of client from main")
}
