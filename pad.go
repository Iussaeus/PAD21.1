package main

import (
	"fmt"
	"os"
	"pad/broker_grpc"
	"pad/broker_tcp"
	"pad/client_grpc"
	"pad/client_tcp"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run pad [tcp|grpc] [broker|client] [Name]")
		return
	}

	appType := os.Args[1]
	whichApp := os.Args[2]

	switch appType {
	case "tcp":
		switch whichApp {
		case "broker":
			broker_tcp.Run()
		case "client":
			if len(os.Args) < 4 {
				fmt.Println("Usage: go run pad [tcp|grpc] [client] [Name]")
				return
			}

			name := os.Args[3]

			switch name {
			case "":
				fmt.Println("Enter a name for the client")
			case "Test1":
				client_tcp.Run(name, "")
			case "Test2":
				client_tcp.Run(name, "")
			default:
				fmt.Println("FOR TEST PURPOSES ONLY Test1 AND Test2 ARE ALLOWED")
			}

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
