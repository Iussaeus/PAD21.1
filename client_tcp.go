package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"log"
)

// NOTE: maybe make a helper module for colored prints

const (
	Red     = "\033[31m"
	Green   = "\033[32m"
	Blue    = "\033[34m"
	Yellow  = "\033[33m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Reset   = "\033[0m"
)

type MessageType uint8

const (
	Unicast = iota
	Multicast 
)

// TODO: use bufio properly

type ClientTCP struct {
	Name      string
	tcpServer *net.TCPAddr
	reader    bufio.Reader
	writer    bufio.Writer
}

// TODO: move the message and the message type to their own module
// TODO: message queue

type Message struct {
	Sender ClientTCP
	Receiver string
	MessageType MessageType
	Content string
}

func (c *ClientTCP) Start() {
	tcpServer, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to resolve tcp address: %s", err)
	}
	c.tcpServer = tcpServer
}

// TODO: send the marshallized data

func (c *ClientTCP) SendMessage(receiver string, msg string) error {
	if c.tcpServer == nil {
		return fmt.Errorf("tcp server is nil")
	}

	m, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("Mashall failed: %v", err)
	}

	fmt.Printf("\n\t%sMarshalled JSON: %v%s", Red, string(m), Reset)

	c1 := ClientTCP{}
	err = json.Unmarshal(m, &c1)
	if err != nil {
		log.Fatalf("Unmarshall failed: %v", err)
	}

	fmt.Printf("\n\t%sUnmarshalled JSON: %v%s", Green, c1, Reset)

	conn, err := net.DialTCP("tcp", nil, c.tcpServer)
	if err != nil {
		log.Fatalf("Failed to connect client: %s", err)
	}

	_, err = conn.Write([]byte(fmt.Sprintf("\nTest message: %s\nfrom: %s", msg, c.Name)))

	return nil
}

func main() {
	fmt.Println("Started Client from main")
	c := ClientTCP{
		Name: "Test",
	}

	c.Start()
	c.SendMessage("", "weird times am i right ?")
	fmt.Println("end of Client from main")
}
