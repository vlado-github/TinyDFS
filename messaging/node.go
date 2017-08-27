package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"persistance"

	"github.com/google/uuid"
)

type Node interface {
	Run() error
	SendMessage(message Message)
	CloseConn() error
}

type node struct {
	connParams  ConnParams
	conn        net.Conn
	fileManager persistance.FileManager
}

// Creates new instance of node
func NewNode(conn ConnParams, fm persistance.FileManager) Node {
	return &node{
		connParams:  conn,
		fileManager: fm,
	}
}

// Starts a node and connects to the queue
func (n *node) Run() error {
	// connects to queue
	var err error
	n.conn, err = net.Dial(n.connParams.Protocol, n.connParams.Ip+":"+n.connParams.Port)

	if err != nil {
		fmt.Println("Error dialing:", err.Error())
		log.Fatal(err)
	} else {
		go n.receiveMessages()
	}
	return err
}

// Sends message to the queue
func (n *node) SendMessage(message Message) {
	//fmt.Print("[Client] Sending: ", message.Text+"\n")
	encoder := json.NewEncoder(n.conn)
	encodeMessage(&message, encoder)
}

// Receives messages from the queue
func (n *node) receiveMessages() {
	decoder := json.NewDecoder(n.conn)
	for {
		var message Message
		err := decodeMessage(&message, decoder)
		if err != nil {
			fmt.Println("Error: Queue connection is closed.", err.Error())
			break
		} else {
			fmt.Println("[Client] Received: ", message.Topic, message.Text)
			if message.Text == "CONN_ACK" {
				fmt.Println("[Client] Connected")
			} else {
				var guid = uuid.New()
				var cmd = persistance.Command{Key: guid, Text: message.Text, Topic: message.Topic}
				n.fileManager.Write(cmd)
			}
		}
	}
}

// Close connection to the queue
func (n *node) CloseConn() error {
	err := n.conn.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("[Client] Connection closed.")
	return err
}
