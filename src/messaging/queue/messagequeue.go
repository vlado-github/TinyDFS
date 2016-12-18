package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// connection params
const (
	ConnHost = "localhost"
	ConnPort = "3333"
	ConnType = "tcp"
)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(ConnType, ConnHost+":"+ConnPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + ConnHost + ":" + ConnPort)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		fmt.Println("Client Connected...")
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go receiveMessage(conn)
	}
}

// Handles incoming messages.
func receiveMessage(conn net.Conn) {
	for {
		scanner := bufio.NewScanner(conn)
		if ok := scanner.Scan(); ok {
			message := scanner.Text()
			// output message received
			fmt.Print("Message Received:", message+"\n")
			// sample process for string received
			newmessage := strings.ToUpper(message)
			// send new string back to client
			conn.Write([]byte(newmessage + "\n"))
		} else {
			conn.Close()
		}
	}
}
