package main

import (
	"fmt"
	"os"
	// "pad/broker_grpc"
	// "pad/broker_tcp"
	// "pad/client_grpc"
	// "pad/client_tcp"
	// "pad/message_tcp"

	dw "pad/data_warehouse"
	"pad/proxy"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run pad [proxy|client] [port]")
		return
	}

	appType := os.Args[1]
	whichApp := os.Args[2]

	switch appType {
	case "dw":
		argType := os.Args[3]
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run pad [proxy|client] [port] [type]")
			return
		}

		if argType == "1" {
			dw.Run(whichApp, "")
		}

		if argType == "2" {
			dw.Run(whichApp, "user=postgres password=yourpassword dbname=dw_db host=localhost sslmode=disable")
		}

	case "proxy":
		if len(os.Args) < 3 {
			proxy.Run(whichApp, nil)
		}
		ports := os.Args[3:]
		proxy.Run(whichApp, ports)

	// case "tcp":
	// 	switch whichApp {
	// 	case "broker":
	// 		broker_tcp.Run()
	// 	case "client":
	// 		if len(os.Args) == 3 {
	// 			client_tcp.Run("", "")
	// 			return
	// 		}
	//
	// 		command := os.Args[4]
	// 		args := os.Args[5:]
	//
	// 		name := os.Args[3]
	//
	// 		switch name {
	// 		case "":
	// 			fmt.Println("Enter a name for the client")
	// 		default:
	// 			client := client_tcp.NewClientTCP(name)
	// 			client.Connect("")
	//
	// 			// get the receiver name and payload
	//
	// 			var payload string
	// 			var topic string
	//
	// 			if len(args) > 1 {
	// 				topic = args[0]
	// 				payload = args[1]
	// 			}
	//
	// 			switch command {
	// 			case "publish":
	// 				client.Publish(payload, message_tcp.Topic(topic))
	// 			case "subscribe":
	// 				client.Subscribe(message_tcp.Topic(topic))
	// 			case "newtopic":
	// 				client.NewTopic(message_tcp.Topic(topic))
	// 			case "topics":
	// 				client.Topics()
	// 			case "unsubscribe":
	// 				client.Unsubscribe(message_tcp.Topic(topic))
	// 			case "delete":
	// 				client.DeleteTopic(message_tcp.Topic(topic))
	// 			case "receive":
	// 				for {
	// 					client.ReadMessage()
	// 				}
	// 			default:
	// 				fmt.Println("Unknown command")
	// 			}
	// 		}
	// 		// pad tcp [broker|client] [name] [publish|subscribe|topics|delete|unsubscribe|newtopic] [topic] [message]
	//
	// 	default:
	// 		fmt.Println("Invalid TCP application. Choose 'broker' or 'client'.")
	// 	}
	// case "grpc":
	// 	switch whichApp {
	// 	// case "broker":
	// 	// 	brokerGrpc.Run()
	// 	// case "client":
	// 	// 	clientGrpc.Run()
	// 	default:
	// 		fmt.Println("Invalid gRPC application. Choose 'broker' or 'client'.")
	// 	}
	default:
		fmt.Println("Usage: go run pad [tcp|grpc] [broker|client]")
	}
}
