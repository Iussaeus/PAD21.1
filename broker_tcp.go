package broker

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// TODO: persistent data storage
// TODO: change the way the clients are stored
// actually store the clients

type BrokerTCP struct {
	listen  net.Listener
	clients map[string]net.Conn
	reader  bufio.Reader
	writer  bufio.Writer
}

// TODO: make it deploy worthy

func (b *BrokerTCP) Open(port string) {
	if port == "" {
		port = "8080"
	}

	listen, err := net.Listen("tcp", "0.0.0.0:"+port)
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
	fmt.Println("Serving")
	for {
		conn, err := b.listen.Accept()
		defer conn.Close()

		if err != nil {
			log.Fatalf("Listener failed to connect: %v", err)
		}

		b.reader = *bufio.NewReader(conn)
		b.writer = *bufio.NewWriter(conn)

		scanner := bufio.NewScanner(conn)
		if scanner.Scan() {
			response := scanner.Text()
			fmt.Printf("Received response: %s\n", response)
		} else if err := scanner.Err(); err != nil {
			log.Fatalf("Client Failed to read %v", err)
		}
		// go func() {
		// 	for {
		// 		b.writer.WriteString("Got conneciton bruh")
		//
		// 		b.writer.Flush()
		// 		var msg string
		// 		msg, err = b.reader.ReadString('\n')
		// 		if err != nil {
		// 			log.Fatalf("Broker failed to read: %v", err)
		// 		}
		// 		fmt.Println(b.reader.Buffered())
		//
		// 		fmt.Printf("\nServer got msg: %s\n\t", msg)
		// 		response := fmt.Sprintf("Server sent the message back: %s\n", msg)
		// 		_, err = b.writer.WriteString(response)
		// 		if err != nil {
		// 			log.Fatalf("Broker failed to write: %v", err)
		// 		}
		// 		b.writer.Flush()
		// 		b.listen.Close()
		// 	}
		// }()
		fmt.Println("End of Serving")
	}
}

func (b *BrokerTCP) Close() {
	b.listen.Close()
}

func main() {
	b := BrokerTCP{}

	b.Open("")
	b.Serve()
}
