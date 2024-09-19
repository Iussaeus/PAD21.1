package brokerTcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	// "time"
)

// TODO: persistent data storage
// TODO: change the way the clients are stored
// actually store the clients

type BrokerTCP struct {
	listen  net.Listener
	clients map[string]net.Conn
	reader  *bufio.Reader
	writer  *bufio.Writer
}

// TODO: make it deploy worthy

func (b *BrokerTCP) Open(port string) {
	if port == "" {
		port = "12345"
	}

	listen, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatalf("Broker failed to listen: %v", err)
	}

	b.listen = listen
}

// TODO: JSON serialization
// TODO:Subscriber/publsher stuff
// TODO:Depending on the type of the message it should call different funcs
// like if it is a subscribe message is it should subcribe the subscirebee to a topic

func (b *BrokerTCP) Serve() {
	for {
		conn, err := b.listen.Accept()
		if err != nil {
			log.Printf("Listener failed to connect: %v", err)
		}

		// err = conn.SetDeadline(<-time.NewTimer(time.Second * 10).C)
		// if err != nil {
		// 	log.Printf("failed to set deadline: %v", err)
		// }
		//
		reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(conn)
		go func(conn net.Conn) {
			msg, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Broker failed to read: %v", err)
			}
			fmt.Printf("\nServer got msg: %s\n\t", msg)
		}(conn)

		go func(conn net.Conn) {
			response := fmt.Sprintf("Server sent the message back: %s\n", "runn")
			_, err = writer.WriteString(response)
			if err != nil {
				log.Printf("Broker failed to write: %v\n", err)
				return
			}

			err = writer.Flush()
			if err != nil {
				log.Printf("Writer failed to flush: %v\n", err)
				return
			}

			// fmt.Println("serving conn: ", conn)
			// i, err := b.writer.WriteString("Got conneciton bruh\n")
			// if err != nil {
			// 	log.Printf("writer failed: %v, writen %v bytes\n", err, i)
			// }
			//
			// // err = b.writer.Flush()
			// if err != nil {
			// 	log.Printf("Writer failed to flush: %v\n", err)
			// }
		}(conn)
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
