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

func NewConn(ip string, port string) (*entities.Conn, error) {

	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	e := &entities.Conn{
		IP:   ip,
		Port: port,
	}

	fmt.Println(ip, port)

	conn, err := d.DialContext(ctx, "tcp", ip+":"+port)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	if _, err := conn.Write([]byte("CONNECT " + e.IP + ":" + fmt.Sprint(e.Port))); err != nil {
		return nil, err
	}

	buffer := make([]byte, 1024)
	if _, err := conn.Read(buffer); err != nil {
		return nil, err
	}

	resp := string(buffer)
	fmt.Println(resp)

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
	fmt.Println(busyPorts)

	return e, nil
}

// Receiver is a function that receives a connection request
// put it in a channel and start a goroutine to listen for it
// and accept it. ch is the channel where the connection request
// is put.
func Receiver(ch chan<- string) {

	fmt.Println("SONO QUI")

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	// Generate random number between 5000 and 8000
	// to use as a port
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

	ch <- "Listening for connections on port " + port

	// Create a new listener
	l, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		ch <- "Error listening: " + err.Error()
	}

	for {

		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			ch <- "Error"
		}

		// Handle the connection in a new goroutine.
		go func() {
			ch <- "Connection got from" + conn.RemoteAddr().String()

			// Read the incoming connection into the buffer.
			buffer := make([]byte, 1024)
			_, err := conn.Read(buffer)
			if err != nil {
				ch <- "Error"
			}

			fmt.Println(string(buffer))

			// Send a response back to person contacting us.
			conn.Write([]byte("CONNECTED on port :" + port + ":"))
			ch <- string(buffer)
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
