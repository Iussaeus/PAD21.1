package client

import (
	"net"
)

type ClientTCP struct {
	ip   net.IP
	name string
}

func (c *ClientTCP) Start() {}

func (c *ClientTCP) Connect() {}

func (c *ClientTCP) Disconnect() {}

func (c *ClientTCP) SendMessage(receiver string, msg string) {} 

func main() {
}
