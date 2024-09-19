package clientTcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

// TODO: use bufio properly
// TODO: send the marshallized data

type ClientTCP struct {
	Name   string
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func (c *ClientTCP) Connect() {
	conn, err := net.Dial("tcp", "localhost:12345")
	if err != nil {
		log.Printf("Failed to connect client: %s", err)
	}

	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)

	// msg, err := c.reader.ReadString('\n')
	// if err != nil {
	// 	log.Printf("Reader failed: %v", err)
	// }
	// fmt.Printf("reader got msg: %v\n", msg)
}

func (c *ClientTCP) Disconnect() {
	c.conn.Close()
}

func (c *ClientTCP) SendMessage(receiver string, msg string) {
	go func() {
		response := fmt.Sprintf("Test message: %s ,from: %s\n", msg, c.Name)

		_, err := c.writer.WriteString(response)
		if err != nil {
			log.Printf("Writer failed: %v", err)
			return
		}
		fmt.Println("after write")

		err = c.writer.Flush()
		if err != nil {
			log.Printf("writer failde to flush: %v", err)
			return
		}
		fmt.Println("after flush")
	}()

	go func() {
		s, err := c.reader.ReadString('\n')
		if err != nil {
			log.Printf("Reader failed: %v", err)
		}
		fmt.Println(s)
		fmt.Println("after read")
		time.Sleep(time.Second * 5)
	}()

	select {
	}
}

func Run() {
	fmt.Println("Started Client from main")
	c := ClientTCP{
		Name: "Test",
	}

	c.Connect()

	c.SendMessage("", "weird times am i right ?")

	fmt.Println("end of Client from main")
}
