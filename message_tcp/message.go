package message_tcp

import "fmt"

type MessageType uint8

type Client interface {
	Connect(port string)
	Disconnect() error
	SendMessage(m Message) error
	ReadMessage() error
}

type Command string

const (
	Publish     Command = "publish"
	Subscribe   Command = "subscribe"
	Unsubscribe Command = "unsubscribe"
	Topics      Command = "topics"
	NewTopic    Command = "newTopic"
	DeleteTopic Command = "delete"
)

type Message struct {
	Sender  string  `json:"sender"`
	Cmd     Command `json:"command"`
	Topic   Topic   `json:"topic"`
	Payload string  `json:"payload"`
}

func (m Message) String() string {
	return fmt.Sprintf("Sender: %s, Cmd: %s, Topic: %s, Payload: %s", m.Sender, m.Cmd, m.Topic, m.Payload)
}

func (m Message) IsEmpty() bool {
	return m == Message{}
}

type Topic string

const (
	Global    Topic = "global"
	Empty     Topic = "empty"
	AllTopics Topic = "topics"
)

type MessageQueue struct {
	Ch chan *Message
}

func NewMessageQueue() *MessageQueue {
	return &MessageQueue{
		Ch: make(chan *Message),
	}
}

func (mq *MessageQueue) Queue(m *Message) {
	mq.Ch <- m
}

func (mq *MessageQueue) Dequeue() (*Message, bool) {
	msg, ok := <-mq.Ch
	return msg, ok 
}

func (mq *MessageQueue) Close() {
	close(mq.Ch)
}
