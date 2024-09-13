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
	go func() {
		for {
			conn, err := b.listen.Accept()
			if err != nil {
				log.Fatalf("Listener failed to connect: %v", err)
			}

			buffer := make([]byte, 1024)

			_, err = conn.Read(buffer)
			fmt.Printf("\nServer got msg: %s\n\t", string(buffer))

			response := fmt.Sprintf("Server sent the message back: %s\n", string(buffer))
			conn.Write([]byte(response))

			conn.Close()
		}
	}()
}

func (b *BrokerTCP) Close() {
	b.listen.Close()
}

func main() {
	fmt.Println("started the server")
	b := BrokerTCP{}

	b.Open("")
	b.Serve()

	fmt.Println("ended server")
}
