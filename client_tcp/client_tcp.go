package clientTcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// TODO: use bufio properly

type ClientTCP struct {
	Name      string
	tcpServer *net.TCPAddr
	reader    bufio.Reader
	writer    bufio.Writer
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

	conn, err := net.DialTCP("tcp", nil, c.tcpServer)
	if err != nil {
		log.Fatalf("Failed to connect client: %s", err)
	}

	_, err = conn.Write([]byte(fmt.Sprintf("\nTest message: %s\nfrom: %s", msg, c.Name)))

	return nil
}

func Run() {
	fmt.Println("Started Client from main")
	c := ClientTCP{
		Name: "Test",
	}

	c.Start()
	c.SendMessage("", "weird times am i right ?")
	fmt.Println("end of Client from main")
}
