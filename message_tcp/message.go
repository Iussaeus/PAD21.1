package message_tcp

import (
)

type Client interface{
	SendMessage(receiver string, message string)
	ReadMessage()
	Connect()
	Disconnect()
}

type MessageType uint8

const (
	Unicast MessageType = iota
	Multicast
)

// TODO: move the message and the message type to their own module
// TODO: message queue

type Message struct {
	Sender      Client
	Receiver    string
	MessageType MessageType
	Payload     string
}

type Topic string

// topics contain a map of the topics and the subsribers
type MessageQueue struct {
	Topics         map[Topic][]Client
	QueuedMessages []Message
}

func (mq *MessageQueue) Queue(m Message) error {
	return nil
}

func (mq *MessageQueue) Dequeue(m Message) error {
	return nil
}
