package messaging

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func InitNode(params connParams) {
	// connect to this socket
	conn, err := net.Dial(params.protocol, params.ip+":"+params.port)
	if err != nil {
		fmt.Println("Error dialing:", err.Error())
		os.Exit(1)
	}
	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nText to send: ")
		text, _ := reader.ReadString('\n')
		// send to socket
		fmt.Fprintf(conn, text+"\n")
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}
