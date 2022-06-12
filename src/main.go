package main

import (
	"bufio"
	"fmt"
	"goxa/src/adapter"
	"net"
	"os"
	"strings"
)

func main() {

	var conn []net.Conn
	r := make(chan string, 10)
	go adapter.Receiver(r)

	for {
		fmt.Print("> ")
		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadSlice('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		// trim the newline
		line = line[:len(line)-1]

		// split the line by spaces
		// and get the first part
		// which is the command
		command := strings.Split(string(line), " ")[0]

		// if user enters "conn"
		// if user enters "exit"
		if command == "exit" {
			conn[0].Close()
		} else if command == "conn" {
			// get the second part
			// which is the ip address
			ip := strings.Split(string(line), " ")[1]
			// get the third part
			// which is the port number
			port := strings.Split(string(line), " ")[2]

			// create a new connection
			c, err := adapter.NewConn(ip, port)
			conn = append(conn, c)

			if err != nil {
				fmt.Println(err)
				continue
			}

			// if the connection is successful
			// then print the connection details
			// else print the error
		} else if command == "add" {
			firstNumber := strings.Split(string(line), " ")[1]
			secondNumber := strings.Split(string(line), " ")[2]
			sum, err := adapter.Add(firstNumber, secondNumber, conn[0])
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Printf("\n< %s + %s = %d", firstNumber, secondNumber, sum)
		} else {
			fmt.Println("Invalid command")
		}

	}
}
