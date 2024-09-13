package broker_tcp

import (
	"fmt"
	"log"
	"net"
	client "pad/client"
)

type Server interface {
	Open() error
	Close() error
}

type Broker struct {
	ip      net.IP
	port    string
	clients []client.Client
}

func (b Broker) Open() {
	fmt.Println("it's a struct")
}

func (b Broker) Close() {
	fmt.Println("it's a struct")
}

func openConn() {
	port := ":8080"
	l, err := net.Listen("tcp", "0.0.0.0"+port)

	if err != nil {
		log.Fatal(err)
	}
	l.Close()
}
