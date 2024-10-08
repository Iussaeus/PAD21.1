package broker_tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"slices"
	"sync"

	"pad/broker"
	"pad/helpers"
	"pad/message_tcp"
)

// TODO: periodic ping pong to see if a client is still connected, if not remove the connection
// TODO: remove clients from slice when they disconnect

var _ broker.Broker = (*BrokerTCP)(nil)

type ClientData struct {
	Name string
}

type ClientBuf struct {
	name   string
	reader *bufio.Reader
	writer *bufio.Writer
	mq     *MessageQueue
	conn   net.Conn
	clChan chan struct{}
}

type MessageQueue struct {
	ch chan Message
}

func (mq *MessageQueue) Queue(m Message) {
	mq.ch <- m
}

func (mq *MessageQueue) Dequeue() Message {
	return <-mq.ch
}

func NewMessageQueue() *MessageQueue {
	return &MessageQueue{
		make(chan Message),
	}
}

type Message struct {
	buf *ClientBuf
	m   message_tcp.Message
}

func NewClientBuf(conn net.Conn) *ClientBuf {
	return &ClientBuf{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
		mq:     NewMessageQueue(),
		clChan: make(chan struct{}),
	}
}

func (cb *ClientBuf) String() string {
	return fmt.Sprintf("%s", cb.name)
}

type BrokerTCP struct {
	listen  net.Listener
	clients []*ClientBuf
	topics  map[message_tcp.Topic][]*ClientBuf
	mq      *MessageQueue
	mu      *sync.RWMutex
	wg      *sync.WaitGroup
}

func (b *BrokerTCP) String() string {
	return fmt.Sprintf("Clients:%d, Topics:%s", len(b.clients), b.topics)
}

func NewBrokerTCP() *BrokerTCP {
	return &BrokerTCP{
		mu:      &sync.RWMutex{},
		wg:      &sync.WaitGroup{},
		mq:      NewMessageQueue(),
		clients: []*ClientBuf{},
		topics:  make(map[message_tcp.Topic][]*ClientBuf),
	}
}

func (b *BrokerTCP) Init(port string) {
	if port == "" {
		port = "12345"
	}

	listen, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatalf("Broker failed to listen: %v", err)
	}

	b.listen = listen
	b.topics[message_tcp.Global] = []*ClientBuf{}
}

func (b *BrokerTCP) Accept() (*ClientBuf, error) {
	conn, err := b.listen.Accept()
	if err != nil {
		return nil, err
	}

	cb := NewClientBuf(conn)

	if _, err = cb.writer.WriteString("Greetings\n"); err != nil {
		return nil, fmt.Errorf("Failed to write to buffer ON CONNECTION: %v\n", err)
	}

	if err = cb.writer.Flush(); err != nil {
		return nil, fmt.Errorf("Failed to flush buffer ON CONNECTION: %v\n", err)
	}

	b.wg.Add(1)
	go func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		defer b.wg.Done()

		if m, err := cb.reader.ReadString('\n'); err != nil {
			log.Printf("Failed to read buffer ON CONNECTION: %v", err)
		} else {
			m = m[:len(m)-1]
			helpers.CPrintf(helpers.Blue, "Got client data ON CONNECTION: %v\n", m)

			c := &ClientData{}
			json.Unmarshal([]byte(m), c)

			cb.name = c.Name

			b.clients = slices.DeleteFunc(b.clients, func(c *ClientBuf) bool { return c.name == cb.name })

			b.clients = append(b.clients, cb)
			b.topics[message_tcp.Global] = append(b.topics[message_tcp.Global], cb)

			helpers.CPrintf(helpers.Blue, "Unmarshallized client data ON CONNECTION: %+v\n", *c)
		}
	}()
	b.wg.Wait()

	return cb, nil
}

func (b *BrokerTCP) Serve() {
	for {
		helpers.CPrintf(helpers.Cyan, "Serving")

		cb, err := b.Accept()
		if err != nil {
			log.Printf("Broker Failed to accept conn: %v\n", err)
			continue
		}

		go func(cb *ClientBuf) {
			for {
				if err := b.ReadMessage(cb); err == io.EOF {
					b.DeleteClientBuf(cb)
					helpers.CPrintf(helpers.Green, "Broker after delete: %s\n", b)
					return
				} else if err != nil {
					helpers.CPrintf(helpers.Red, "Broker: error while reading: %v\n", err)
				}
			}
		}(cb)

		go func() {
			for {
				if err := b.SendMessages(); err != nil {
					helpers.CPrintf(helpers.Red, "Broker: error while sending: %v\n", err)
				}
				helpers.CPrintf(helpers.Yellow, "Broker: %s\n", b)
			}
		}()

		helpers.CPrintf(helpers.Cyan, "Done serving")
	}
}

func (b *BrokerTCP) SendMessages() error {
	msg := b.mq.Dequeue()

	helpers.CPrintf(helpers.Magenta, "Got msg %s\n", msg)
	switch msg.m.Cmd {
	case message_tcp.Publish:
		return b.Publish(msg, msg.m.Topic)

	case message_tcp.NewTopic:
		return b.AddTopic(msg.m.Topic)

	case message_tcp.DeleteTopic:
		return b.DeleteTopic(msg.m.Topic)

	case message_tcp.Subscribe:
		return b.Subscribe(msg.buf, msg.m.Topic)

	case message_tcp.Unsubscribe:
		return b.Unsubscribe(msg.buf, msg.m.Topic)

	case message_tcp.Topics:
		return b.Topics(msg.buf)

	case "":
		return nil

	default:
		return fmt.Errorf("Unkown cmd\n")
	}
}

func (b *BrokerTCP) AddTopic(t message_tcp.Topic) error {
	if _, ok := b.topics[t]; ok {
		return fmt.Errorf("Topic \"%s\" already exists\n", t)
	}

	// helpers.CPrintf(helpers.Magenta, "Added Topic")
	b.topics[t] = []*ClientBuf{}

	return nil
}

func (b *BrokerTCP) DeleteTopic(t message_tcp.Topic) error {
	if _, ok := b.topics[t]; ok {
		// helpers.CPrintf(helpers.Magenta, "Deleting Topic")
		b.mu.Lock()
		delete(b.topics, t)
		b.mu.Unlock()

		return nil
	}

	return fmt.Errorf("Topic \"%s\" doesn't exist\n", t)
}

func (b *BrokerTCP) Subscribe(cb *ClientBuf, t message_tcp.Topic) error {
	if c, ok := b.topics[t]; ok {
		for _, client := range c {
			if client == cb {
				return fmt.Errorf("Client %s is already subscribed to %s", cb.name, t)
			}
		}
		// helpers.CPrintf(helpers.Magenta, "Subscribing")
		b.mu.Lock()
		b.topics[t] = append(c, cb)
		b.mu.Unlock()
		return nil
	}

	return fmt.Errorf("Topic \"%s\" doesn't exist\n", t)
}

func (b *BrokerTCP) Unsubscribe(cb *ClientBuf, t message_tcp.Topic) error {
	if clients, ok := b.topics[t]; ok {
		// helpers.CPrintf(helpers.Magenta, "Unsubscribing")
		for i, c := range clients {
			if c == cb {
				b.mu.Lock()
				b.topics[t] = append(clients[:i], clients[i+1:]...)
				b.mu.Unlock()
				return nil
			}
		}
		return fmt.Errorf("Client \"%s\" doesn't exist\n", cb.name)
	}

	return fmt.Errorf("Topic \"%s\" doesn't exist\n", t)
}

func (b *BrokerTCP) Publish(m Message, t message_tcp.Topic) error {
	if clients, ok := b.topics[t]; ok {
		// helpers.CPrintf(helpers.Magenta, "Publishing")
		for _, client := range clients {
			if  m.buf != client {
				// helpers.CPrintf(helpers.Magenta, "Done Publishing")
				b.SendMessage(client, m.m)
			}
		}
		return nil
	}
	return fmt.Errorf("Topic \"%s\" doesn't exist\n", t)
}

func (b *BrokerTCP) Topics(s *ClientBuf) error {
	topics := ""

	b.mu.Lock()
	for topic := range b.topics {
		topics += string(topic) + " "
	}
	b.mu.Unlock()

	m := message_tcp.Message{
		Sender:  "Broker",
		Topic:   message_tcp.AllTopics,
		Payload: topics,
	}

	// helpers.CPrintf(helpers.Magenta, "Sending Topics")
	b.SendMessage(s, m)
	return nil
}

func (b *BrokerTCP) SendMessage(cb *ClientBuf, m message_tcp.Message) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	d, err := json.Marshal(m)

	_, err = cb.writer.WriteString(string(d) + "\n")
	if err != nil {
		return fmt.Errorf("Broker failed to write: %v\n", err)
	}

	err = cb.writer.Flush()
	if err != nil {
		return fmt.Errorf("Writer failed to flush: %v\n", err)
	}

	return nil
}

// TODO: check the cmd, type, topic of the message and act accordingly

func (b *BrokerTCP) ReadMessage(cb *ClientBuf) error {
	s, err := cb.reader.ReadString('\n')
	if err != nil {
		return err
	}

	var msg message_tcp.Message

	err = json.Unmarshal([]byte(s), &msg)
	if err != nil {
		return err
	}

	msgBuf := &Message{m: msg, buf: cb}

	fmt.Printf("Server got msg: %v\n", msg)

	b.mu.Lock()
	b.mq.Queue(*msgBuf)
	b.mu.Unlock()

	return nil
}

func (b *BrokerTCP) DeleteClientBuf(cb *ClientBuf) {
	cb.conn.Close()

	for t, c := range b.topics {
		for i, client := range c {
			if cb == client {
				b.mu.Lock()
				b.topics[t] = append(c[:i], c[i+1:]...)
				b.mu.Unlock()
			}
		}
	}
}

func (b *BrokerTCP) Close() error {
	return fmt.Errorf("Error while closing the broker: %v\n", b.listen.Close())
}

func Run() {
	helpers.AwaitSIGINT()

	fmt.Println("Started Broker")

	b := NewBrokerTCP()
	b.Init("")
	b.Serve()
}
