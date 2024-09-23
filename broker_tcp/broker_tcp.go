package broker_tcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"pad/client_tcp"
	"pad/message_tcp"
	"pad/helpers"
	"sync"
	// "time"
)

// TODO: make it deploy worthy

// TODO: persistent data storage
// TODO: change the way the clients are stored
// actually store the clients

type BrokerTCP struct {
	listen       net.Listener
	clients      map[net.Conn]clientTcp.ClientTCP
	messageQueue *message_tcp.MessageQueue
	reader       *bufio.Reader
	writer       *bufio.Writer
	mutex        *sync.Mutex
	wg           *sync.WaitGroup
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
	b.mutex = &sync.Mutex{}
	b.wg = &sync.WaitGroup{}
}

// TODO: JSON serialization
// TODO:Subscriber/publsher stuff
// TODO:Depending on the type of the message it should call different funcs
// TODO: MessageQUEUE
// like if it is a subscribe message is it should subcribe the subscirebee to a topic

func (b *BrokerTCP) Accept() (*bufio.Reader, *bufio.Writer, error) {
	conn, err := b.listen.Accept()
	if err != nil {
		return nil, nil, err
	}

	// err = conn.SetDeadline(<-time.NewTimer(time.Second * 10).C)
	// if err != nil {
	// 	log.Printf("failed to set deadline: %v", err)
	// }

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	msg := fmt.Sprintf("Greetings\n")
	writer.WriteString(msg)
	if err != nil {
		log.Printf("Failed to write to buffer ON CONNECTION: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Printf("Failed to flush buffer ON CONNECTION: %v", err)
	}

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		msg, err = reader.ReadString('\n')
		if err != nil {
			log.Printf("Failed to read buffer ON CONNECTION: %v", err)
		}
		msg = msg[:len(msg)-1]
		helpers.CPrintf(helpers.Blue, "Got client data: %v\n", msg)
	}()
	b.wg.Wait()

	return reader, writer, nil
}
func (b *BrokerTCP) Serve() {
	for {
		reader, writer, err := b.Accept()

		msgchan := make(chan string)

		go func() {
			msg, err := reader.ReadString('\n')
			if err != nil {
				helpers.CPrintf(helpers.Red, "Broker failed to read: %v\n", err)
			}

			if msg != "" {
				msgchan <- msg
				fmt.Printf("\nServer got msg: %s\n", msg)
			}
		}()

		go func() {
			response := fmt.Sprintf("Server sent the message back: %s\n", <-msgchan)
			_, err = writer.WriteString(response)
			if err != nil {
				helpers.CPrintf(helpers.Blue, "Broker failed to write: %v\n", err)
				return
			}

			err = writer.Flush()
			if err != nil {
				log.Printf("Writer failed to flush: %v\n", err)
				return
			}
		}()

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

func (b *BrokerTCP) Close() {
	b.listen.Close()
}

func Run() {
	b := BrokerTCP{}
	b.Open("")
	b.Serve()
}
