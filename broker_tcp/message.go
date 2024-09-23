package broker_tcp

import "pad/client_tcp"

// TODO: move the message and the message type to their own module
// TODO: message queue

type MessageType uint8

const (
	Unicast MessageType = iota
	Multicast
	Subscribe
	Unsubscribe
)

type Message struct {
	Sender      clientTcp.ClientTCP
	Receiver    string
	MessageType MessageType
	Payload     string
}

type Topic string

// topics contain a map of the topics and the subsribers
type MessageQueue struct {
	Topics         map[Topic][]clientTcp.ClientTCP
	QueuedMessages []Message
}

func (mq *MessageQueue) Queue(m Message) error {
	return nil
}

func (mq *MessageQueue) Dequeue(m Message) error {
	return nil
}
