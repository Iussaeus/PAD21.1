package client_tcp

import (
	"net"
)

type Client struct {
	ip   net.IP
	name string
}

func (c *Client) Start() {}

func (c *Client) Connect() {}

func (c *Client) Disconnect() {}

func (c *Client) SendMessage(receiver string, msg string) {} 

func main() {
}
