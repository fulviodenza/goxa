package main

import (
	"fmt"
	"goxa/src/adapter"
	"os"
)

func main() {

	r := make(chan string, 10)

	if len(os.Args) > 1 {
		command := os.Args[1]

		switch command {
		case "conn":
			adapter.NewConn(os.Args[2], os.Args[3])
		}
	}

	go adapter.Receiver(r)
	for {
		select {
		case msg := <-r:
			fmt.Println(msg)
		}
	}
}
