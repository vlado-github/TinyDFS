package messaging

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Pool struct {
	conns []net.Conn
}

var pool = Pool{}

func InitQueue(params connParams) {
	// Listen for incoming connections.
	l, err := net.Listen(params.protocol, params.ip+":"+params.port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + params.ip + ":" + params.port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		pool.conns = append(pool.conns, conn)
		fmt.Println("Client Connected...Total:",len(pool.conns))
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
			fmt.Print("\nMessage Received:", message+"\n")
			// sample process for string received
			newmessage := strings.ToUpper(message)
			// send new string back to client
			conn.Write([]byte(newmessage + "\n"))
		} else {
			conn.Close()
		}
	}
}

func sendMessages(){
	for {

	}
}

