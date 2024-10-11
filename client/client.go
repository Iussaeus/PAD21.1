package client

import "pad/message_tcp"

type Client interface {
	SendMessage(msg message_tcp.Message)
	ReadMessage()
	Connect()
	Disconnect()
}
