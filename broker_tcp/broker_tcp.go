package broker_tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"pad/broker"
	"pad/client_tcp"
	"pad/helpers"
	"pad/message_tcp"
)

// TODO: make it deploy worthy

// TODO: persistent data storage
// TODO: remove clients from slice when they disconnect

var _ broker.Broker = (*BrokerTCP)(nil)

// mq is a client queue
type ClientBuf struct {
	Client *client_tcp.ClientTCP
	reader *bufio.Reader
	writer *bufio.Writer
	mq     *message_tcp.MessageQueue
	topics []message_tcp.Topic
	conn   net.Conn
}

func (cb *ClientBuf) String() string {
	return fmt.Sprintf("ClientBuf: conn=%v, r=%p, w=%p", cb.conn, cb.reader, cb.writer)
}

// mq is a global message queue
type BrokerTCP struct {
	listen  net.Listener
	clients []*ClientBuf
	mq      *message_tcp.MessageQueue
	mu      *sync.Mutex
	wg      *sync.WaitGroup
}

func (b *BrokerTCP) Open(port string) {
	if port == "" {
		port = "12345"
	}

	listen, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatalf("Broker failed to listen: %v", err)
	}

	b.listen = listen
	b.mu = &sync.Mutex{}
	b.wg = &sync.WaitGroup{}
	b.mq = message_tcp.NewMessageQueue()
}

// TODO: JSON serialization
// TODO:Subscriber/publsher stuff
// TODO:Depending on the type of the message it should call different funcs
// TODO: MessageQUEUE

// like if it is a subscribe message is it should subcribe the subscirebee to a topic
func (b *BrokerTCP) Accept() (*ClientBuf, error) {
	conn, err := b.listen.Accept()
	if err != nil {
		return nil, err
	}

	// err = conn.SetDeadline(<-time.NewTimer(time.Second * 10).C)
	// if err != nil {
	// 	log.Printf("failed to set deadline: %v", err)
	// }

	cb := &ClientBuf{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
		mq:     message_tcp.NewMessageQueue(),
	}

	if _, err = cb.writer.WriteString("Greetings\n"); err != nil {
		log.Printf("Failed to write to buffer ON CONNECTION: %v\n", err)
	}

	if err = cb.writer.Flush(); err != nil {
		log.Printf("Failed to flush buffer ON CONNECTION: %v\n", err)
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

			cb.Client = &client_tcp.ClientTCP{}
			json.Unmarshal([]byte(m), cb.Client)

			b.clients = append(b.clients, cb)

			b.printClientBuf()
			helpers.CPrintf(helpers.Blue, "Unmarshallized client data ON CONNECTION: %+v\n", *cb.Client)
		}
	}()
	b.wg.Wait()

	return cb, nil
}

func (b *BrokerTCP) printClientBuf() {
	helpers.CPrintf(helpers.Green, "ClientBufSlice:")
	for i, n := range b.clients {
		helpers.CPrintf(helpers.Green, "%d, %s\n", i, n)
	}
}

func (b *BrokerTCP) Serve() {
	for {
		cb, err := b.Accept()
		if err != nil {
			log.Printf("Broker Failed to accept conn: %v\n", err)
		}

		helpers.Wait(
			b.wg,
			func() {
				b.ReadMessage(cb.reader)
			},
			func() {
				b.SendMessage(cb.writer)
			},
		)

		// go func() {
		// 	for {
		// 		time.Sleep(time.Second * 1)
		// 		msg := fmt.Sprintln("Do you get the message?")
		//
		// 		_, err = writer.WriteString(msg)
		// 		if err != nil {
		// 			helpers.CPrintf(helpers.Blue, "Broker failed to write a periodic message: %v\n", err)
		// 			return
		// 		}
		//
		// 		err = writer.Flush()
		// 		if err != nil {
		// 			log.Printf("Writer failed to flush a periodic message: %v\n", err)
		// 			return
		// 		}
		// 	}
		// }()
	}
}

// TODO: func to send a message to a given client(by name or by reference)
// TODO: func to send a message to multiple clients
// TODO: Send message to the clients with respective topics
// or to a client topic which is kind of the same
// hat means that every client is subsribed to it's own topic but can write to it, only read from it
// also others can't read from another clients topic

func (b *BrokerTCP) SendMessage(w *bufio.Writer) {
	m := b.mq.Dequeue()

	response := fmt.Sprintf("Server sent the message back: %v\n", m)

	_, err := w.WriteString(response)
	if err != nil {
		helpers.CPrintf(helpers.Red, "Broker failed to write: %v\n", err)
		return
	}

	err = w.Flush()
	if err != nil {
		log.Printf("Writer failed to flush: %v\n", err)
		return
	}
}

func getClientsByTopic[T ~string](clients []*ClientBuf, t T) ([]*ClientBuf, bool) {
	c := []*ClientBuf{}
	for _, client := range clients {
		for _, topic := range client.topics {
			if string(topic) == string(t) {
				c = append(c, client)
			}
		}
	}
	if len(c) == 0 {
		return nil, false
	}
	return c, true
}

// TODO: check the type of the message and act accordingly
// TODO: actually implement the commands, if a command is passed it is executed

func (b *BrokerTCP) ReadMessage(r *bufio.Reader) {
	mStr, err := r.ReadString('\n')
	if err != nil {
		helpers.CPrintf(helpers.Red, "Broker failed to read: %v\n", err)
	}
	var msg message_tcp.Message
	err = json.Unmarshal([]byte(mStr), &msg)
	if err != nil {
		log.Printf("Failed to unmarshall on write: %v", err)
	}

	b.mq.Queue(msg)

	fmt.Printf("Server got msg: %v\n", msg)
	return
}

func (b *BrokerTCP) Close() {
	b.listen.Close()
}

func Run() {
	b := BrokerTCP{}
	b.Open("")
	b.Serve()
}
