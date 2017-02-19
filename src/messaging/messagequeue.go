package messaging

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func InitQueue(params connParams) {
	// Listen for incoming connections.
	l, err := net.Listen(params.protocol, params.ip+":"+params.port)
	if err != nil {
		fmt.Println("[Queue] Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("[Queue] Listening on " + params.ip + ":" + params.port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("[Queue] Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("[Queue] Client Connected...")
		conn.Write([]byte("CONN_ACK\n"))
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
			fmt.Print("[Queue] Message Received:", message+"\n")
			// sample process for string received
			newmessage := strings.ToUpper(message)
			// send new string back to client
			conn.Write([]byte(newmessage + "\n"))
		} else {
			fmt.Print("[Queue] Connection closed.")
			conn.Close()
			break
		}
	}
}
