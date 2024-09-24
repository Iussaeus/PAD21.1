package message_tcp

type MessageType uint8

type Client interface {
	SendMessage(msg Message)
	ReadMessage()
	Connect()
	Disconnect()
}

const (
	Unicast MessageType = iota
	Multicast
	Subscribtion
	Publication
)

type Message struct {
	Sender   string      `json:"sender"`
	Type     MessageType `json:"type"`
	Topic    Topic       `json:"topic"`
	Cmd      Command     `json:"command"`
	Receiver string      `json:"receiver"`
	Payload  string      `json:"payload"`
}

type Command string

const (
	Publish     = "publish"
	Subscribe   = "subscribe"
	Unsubscribe = "unsubscribe"
)
func (m *Message) IsEmpty() bool {
	return *m == Message{}
}

type Topic string

const Global Topic = "Global"

type MessageQueue struct {
	queuedMessages chan Message
}

func NewMessageQueue() *MessageQueue {
	return &MessageQueue{
		queuedMessages: make(chan Message),
	}
}

func (mq *MessageQueue) Queue(m Message) {
	mq.queuedMessages <- m
}

func (mq *MessageQueue) Dequeue() Message {
	return <-mq.queuedMessages
}
