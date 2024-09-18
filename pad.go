package main

import (
	"fmt"
	"os"
	"pad/broker_tcp"
	"pad/client_tcp"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: go run main.go [appType] [whichApp]")
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
            fmt.Println("Invalid TCP application. Choose 'server' or 'client'.")
        }
    // case "grpc":
    //     switch whichApp {
    //     case "broker":
    //         .RunServer()
    //     case "client":
    //        grpc.RunClient()
    //     default:
    //         fmt.Println("Invalid gRPC application. Choose 'server' or 'client'.")
    //     }
    default:
        fmt.Println("Invalid application type. Choose 'tcp', 'grpc', 'app1', or 'app2'.")
    }
}
