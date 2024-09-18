package messagetcp

import (
	"pad/client_tcp"
)

type MessageType uint8

const (
	Unicast MessageType = iota
	Multicast
)

// TODO: move the message and the message type to their own module
// TODO: message queue

type Message struct {
	Sender      clientTcp.ClientTCP
	Receiver    string
	MessageType MessageType
	Content     string
}
