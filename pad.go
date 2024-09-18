package main

import (
	"pad/broker"
)

func main() {
	b := broker.BrokerTCP{}
	b.Open("")
	b.Serve()
}
