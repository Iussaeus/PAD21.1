package broker_tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"slices"
	"sync"

	"pad/broker"
	"pad/helpers"
	"pad/message_tcp"
)

// TODO: make it deploy worthy

// TODO: persistent data storage
// TODO: remove clients from slice when they disconnect

var _ broker.Broker = (*BrokerTCP)(nil)

// mq is a client queue
type ClientData struct {
	Name string
}

type ClientBuf struct {
	name   string
	reader *bufio.Reader
	writer *bufio.Writer
	mq     *message_tcp.MessageQueue
	conn   net.Conn
}

func NewClientBuf(conn net.Conn) *ClientBuf {
	return &ClientBuf{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
		mq:     message_tcp.NewMessageQueue(),
	}
}

func (cb *ClientBuf) String() string {
	return fmt.Sprintf("ClientBuf: %s", cb.name)
}

type BrokerTCP struct {
	listen  net.Listener
	clients []*ClientBuf
	topics  map[message_tcp.Topic][]*ClientBuf
	mq      *message_tcp.MessageQueue
	mu      *sync.RWMutex
	wg      *sync.WaitGroup
}

func (b *BrokerTCP) String() string {
	return fmt.Sprintf("Clients: %d Topics:%s", len(b.clients), b.topics)
}

func NewBrokerTCP() *BrokerTCP {
	return &BrokerTCP{
		mu:      &sync.RWMutex{},
		wg:      &sync.WaitGroup{},
		mq:      message_tcp.NewMessageQueue(),
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
				if err := b.ReadMessage(cb); err != nil {
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
	errChan := make(chan error)
	defer close(errChan)

	msg := b.mq.Dequeue()
	helpers.CPrintf(helpers.Magenta, "Got msg %s\n", msg)
	switch msg.Cmd {
	case message_tcp.Publish:
		return b.Publish(*msg, msg.Topic)
	case message_tcp.NewTopic:
		return b.AddTopic(msg.Topic)
	case message_tcp.DeleteTopic:
		return b.DeleteTopic(msg.Topic)
	case message_tcp.Subscribe:
		c, err := b.Client(msg.Sender)
		if err != nil {
			return err
		}

		return b.Subscribe(c, msg.Topic)
	case message_tcp.Unsubscribe:
		c, err := b.Client(msg.Sender)
		if err != nil {
			return err
		}

		return b.Unsubscribe(c, msg.Topic)
	case message_tcp.Topics:
		c, err := b.Client(msg.Sender)
		if err != nil {
			return err
		}
		return b.Topics(c)
	case "":
		return nil
	default:
		return fmt.Errorf("Unkown cmd\n")
	}
}

func (b *BrokerTCP) Client(n string) (*ClientBuf, error) {
	for _, c := range b.clients {
		if c.name == n {
			return c, nil
		}
	}

	return nil, fmt.Errorf("Client %s not found\n", n)
}

func (b *BrokerTCP) AddTopic(t message_tcp.Topic) error {
	if _, ok := b.topics[t]; ok {
		return fmt.Errorf("Topic \"%s\" already exists\n", t)
	}

	helpers.CPrintf(helpers.Magenta, "Added Topic")
	b.topics[t] = []*ClientBuf{}

	return nil
}

func (b *BrokerTCP) DeleteTopic(t message_tcp.Topic) error {
	if _, ok := b.topics[t]; ok {
		helpers.CPrintf(helpers.Magenta, "Deleting Topic")
		delete(b.topics, t)
		return nil
	}

	return fmt.Errorf("Topic \"%s\" doesn't exist\n", t)
}

func (b *BrokerTCP) Subscribe(cb *ClientBuf, t message_tcp.Topic) error {
	if bt, ok := b.topics[t]; ok {
		helpers.CPrintf(helpers.Magenta, "Subscribing")
		b.topics[t] = append(bt, cb)
		return nil
	}

	return fmt.Errorf("Topic \"%s\" doesn't exist\n", t)
}

func (b *BrokerTCP) Unsubscribe(cb *ClientBuf, t message_tcp.Topic) error {
	if clients, ok := b.topics[t]; ok {
		helpers.CPrintf(helpers.Magenta, "Unsubscribing")
		for i, c := range clients {
			if c == cb {
				b.topics[t] = append(clients[:i], clients[i+1:]...)
				return nil
			}
		}
		return fmt.Errorf("Client \"%s\" doesn't exist\n", cb.name)
	}

	return fmt.Errorf("Topic \"%s\" doesn't exist\n", t)
}

func (b *BrokerTCP) Publish(m message_tcp.Message, t message_tcp.Topic) error {
	b.wg.Add(1)
	defer b.wg.Done()

	if clients, ok := b.topics[t]; ok {
		helpers.CPrintf(helpers.Magenta, "Publishing")
		sender, err := b.Client(m.Sender)
		if err != nil {
			return err
		}
		for _, client := range clients {
			if sender != client {
				helpers.CPrintf(helpers.Magenta, "Done Publishing")
				b.SendMessage(client, m)
			}
		}
	}

	return fmt.Errorf("Topic \"%s\" doesn't exist\n", t)
}

func (b *BrokerTCP) Topics(s *ClientBuf) error {
	topics := ""
	for topic := range b.topics {
		topics += " " + string(topic)
	}
	m := message_tcp.Message{
		Sender:  "Broker",
		Topic:   message_tcp.AllTopics,
		Payload: topics,
	}

	helpers.CPrintf(helpers.Magenta, "Sending Topics")
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
		return fmt.Errorf("Broker failed to read: %v\n", err)
	}

	var msg message_tcp.Message

	err = json.Unmarshal([]byte(s), &msg)
	if err != nil {
		return fmt.Errorf("Failed to unmarshall on read: %v", err)
	}

	// queue in broker if the message is global queue in client if the message is to a topic

	fmt.Printf("Server got msg: %v\n", msg)

	b.mu.Lock()
	b.mq.Queue(&msg)
	b.mu.Unlock()

	return nil
}

func (b *BrokerTCP) Close() error {
	return fmt.Errorf("Error while closing the broker: %v\n", b.listen.Close())
}

func Run() {
	helpers.AwaitSIGINT()

	// port := "12345"
	fmt.Println("Started Broker")

	b := NewBrokerTCP()
	b.Init("")
	b.Serve()
}
