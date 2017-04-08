package messaging

import (
	"bufio"
	"fmt"
	"net"
	"log"
)

type Node interface {
	Run() error
	SendMessage(message string)
	CloseConn() error
}

type node struct {
	connParams ConnParams
	conn net.Conn
}

// Creates new instance of node
func NewNode(conn ConnParams) Node{
	return &node{
		connParams: conn,
	}
}

// Starts a node and connects to the queue
func (n *node) Run() error {
	// connects to queue
	var err error
	n.conn, err = net.Dial(n.connParams.Protocol, n.connParams.Ip + ":" + n.connParams.Port)

	if err != nil {
		fmt.Println("Error dialing:", err.Error())
		log.Fatal(err)
	} else {
		go n.receiveMessages()
	}
	return err
}

// Sends message to the queue
func (n *node) SendMessage(message string){
	fmt.Println("[Client] Sending: ", message)
	fmt.Fprintf(n.conn, message + "\n")
}

// Receives messages from the queue
func (n *node) receiveMessages(){
	for {
		message, _ := bufio.NewReader(n.conn).ReadString('\n')
		fmt.Println("[Client] Received: ", message)
		if message == "CONN_ACK\n" {
			fmt.Println("[Client] Connected")
		}
	}
}

// Close connection to the queue
func (n *node) CloseConn() error{
	err := n.conn.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("[Client] Connection closed.")
	return err
}


