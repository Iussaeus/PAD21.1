package main

import (
	"fmt"
	"os"
	"pad/broker_tcp"
	"pad/client_tcp"
	"pad/broker_grpc"
	"pad/client_grpc"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run pad [tcp|grpc] [broker|client]")
		return
	}

	appType := os.Args[1]
	whichApp := os.Args[2]

	switch appType {
	case "tcp":
		switch whichApp {
		case "broker":
			brokerTcp.Run()
		case "client":
			clientTcp.Run()
		default:
			fmt.Println("Invalid TCP application. Choose 'broker' or 'client'.")
		}
	case "grpc":
		switch whichApp {
		case "broker":
			brokerGrpc.Run()
		case "client":
			clientGrpc.Run()
		default:
			fmt.Println("Invalid gRPC application. Choose 'broker' or 'client'.")
		}
	default:
		fmt.Println("Usage: go run pad [tcp|grpc] [broker|client]")
	}
}
