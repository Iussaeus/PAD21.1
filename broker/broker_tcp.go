package broker

import (
	"fmt"
	"log"
	"net"
	client "pad/client"
)

type BrokerTCP struct {
	ip      net.IP
	port    string
	clients map[string]client.ClientTCP
}

func (b BrokerTCP) Open() {
	fmt.Println("it's a struct")
}

func (b BrokerTCP) Close() {
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
