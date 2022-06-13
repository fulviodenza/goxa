package adapter

import (
	"context"
	"fmt"
	"goxa/src/adapter/entities"
	"net"
	"strconv"
	"strings"
	"time"
)

func NewConn(ip string, port string) (net.Conn, error) {

	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	e := &entities.Conn{
		IP:   ip,
		Port: port,
	}

	conn, err := d.DialContext(ctx, "tcp", ip+":"+port)

	if err != nil {
		return nil, err
	}

	if _, err := conn.Write([]byte("CONNECT " + e.IP + ":" + fmt.Sprint(e.Port))); err != nil {
		return nil, err
	}

	buffer := make([]byte, 1024)
	if _, err := conn.Read(buffer); err != nil {
		return nil, err
	}

	resp := string(buffer)
	fmt.Printf("< %s\n> ", resp)

	return conn, nil
}

// Receiver is a function that receives a connection request
// put it in a channel and start a goroutine to listen for it
// and accept it. ch is the channel where the connection request
// is put.
func Receiver(ch chan<- string) {

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		fmt.Printf("\n< SERVER: Error listening: %s\n> ", err.Error())
	}

	_, port, err := net.SplitHostPort(l.Addr().String())
	fmt.Printf("\n< SERVER: Listening for connections on port %s\n> ", port)

	for {

		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("\n< SERVER: Error accepting: %s\n> ", err.Error())
		}

		// Handle the connection in a new goroutine.
		go func() {
			fmt.Printf("\n< SERVER: Connection got from %s\n> ", conn.RemoteAddr().String())

			for {
				// Read the incoming connection into the buffer.
				buffer := make([]byte, 1024)
				_, err := conn.Read(buffer)
				if err != nil {
					fmt.Printf("\n< SERVER: Error reading: %s\n> ", err.Error())
					break
				}

				// Check if first bytes are CONNECT
				if string(buffer[0:7]) == "CONNECT" {
					// Send a response back to person contacting us.
					conn.Write([]byte("CONNECTED on port: " + port))
				} else {

					// Handle the buffer
					b, _ := handleBuffer(buffer)
					// Send a response back to person contacting us.
					conn.Write(b)
				}
			}
		}()
	}
}

func handleBuffer(buffer []byte) ([]byte, error) {

	// Split the buffer by spaces
	split := strings.Split(string(buffer), " ")
	switch split[0] {
	case "ADD":
		v1, err := strconv.Atoi(split[1])
		if err != nil {
			return nil, err
		}
		v2, err := strconv.Atoi(split[2])
		if err != nil {
			return nil, err
		}
		return []byte(fmt.Sprint(v1 + v2)), nil
	}
	return []byte("EMPTY_RESPONSE"), nil
}

func contains(n int, a []int) bool {
	for a := range a {
		if a == n {
			return true
		}
	}

	return false
}

func Add(a, b string, conn net.Conn) (int, error) {

	if _, err := conn.Write([]byte("ADD " + a + " " + b + " ")); err != nil {
		return 0, err
	}

	buffer := make([]byte, 1024)
	if _, err := conn.Read(buffer); err != nil {
		return 0, err
	}

	resp := string(buffer[0])

	return strconv.Atoi(resp)
}

// func Collatz(n1, n2 int, conn net.Conn) (bool, error) {

// }
