package adapter

import (
	"context"
	"fmt"
	"goxa/src/adapter/entities"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

var busyPorts = []int{}

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

	// if strings.Contains(resp, "END_CONN") {
	// 	break
	// }

	// split string by : and get the second part
	// which is the port number
	portN, err := strconv.Atoi(strings.Split(resp, ":")[1])
	if err != nil {
		return nil, err
	}

	busyPorts = append(busyPorts, portN)

	return conn, nil
}

// Receiver is a function that receives a connection request
// put it in a channel and start a goroutine to listen for it
// and accept it. ch is the channel where the connection request
// is put.
func Receiver(ch chan<- string) {

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	portN := r.Intn(8000-5000) + 5000

	// If the port is already in use, generate a new one
	for {
		if !contains(portN, busyPorts) {
			break
		} else {
			portN = r.Intn(8000-5000) + 5000
			busyPorts = append(busyPorts, portN)
		}
	}

	port := fmt.Sprint(portN)

	fmt.Printf("\n< SERVER: Listening for connections on port %s\n> ", port)

	l, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Printf("\n< SERVER: Error listening: %s\n> ", err.Error())
	}

	for {

		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("\n< SERVER: Error accepting: %s\n> ", err.Error())
		}

		// Handle the connection in a new goroutine.
		go func() {
			fmt.Printf("\n< SERVER: Connection got from %s\n> ", conn.RemoteAddr().String())

			// Read the incoming connection into the buffer.
			buffer := make([]byte, 1024)
			_, err := conn.Read(buffer)
			if err != nil {
				fmt.Printf("\n< SERVER: Error reading: %s\n> ", err.Error())
			}

			// Send a response back to person contacting us.
			conn.Write([]byte("CONNECTED on port :" + port + ":"))
		}()
	}
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

	if _, err := conn.Write([]byte(a + "+" + b)); err != nil {
		return 0, err
	}

	buffer := make([]byte, 1024)
	if _, err := conn.Read(buffer); err != nil {
		return 0, err
	}

	resp := string(buffer)

	return strconv.Atoi(resp)
}
