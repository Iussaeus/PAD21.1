package main

import (
	"fmt"
	"os"

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

	default:
		fmt.Println("Usage: go run pad [proxy|client] [port]")
	}
}
