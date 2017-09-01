package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"persistance"
	"time"

	"math/rand"

	"github.com/google/uuid"
)

type Node interface {
	Run() error
	SendMessage(message Message)
	CloseConn() error

	GetID() uuid.UUID
	GetElectionID() int
	GetIP() (string, error)
}

type node struct {
	id          uuid.UUID
	electionID  int
	connParams  ConnParams
	conn        net.Conn
	fileManager persistance.FileManager
	isMaster    bool
	queue       MessageQueue
}

// Creates new instance of node
func NewNode(conn ConnParams, master bool) Node {
	rand.Seed(time.Now().Unix())
	uniqueID := uuid.New()
	randomID := rand.Int()
	fm := persistance.NewFileManager(getCurrentDirectory() + "//" + uniqueID.String())
	var msgQueue MessageQueue
	if master {
		msgQueue = NewQueue(conn)
	}

	return &node{
		id:          uniqueID,
		electionID:  randomID,
		connParams:  conn,
		fileManager: fm,
		isMaster:    master,
		queue:       msgQueue,
	}
}

// Returns the Node unique ID
func (n *node) GetID() uuid.UUID {
	return n.id
}

// Returns the Node master-election ID
func (n *node) GetElectionID() int {
	return n.electionID
}

// Return current IP address of the device
func (n *node) GetIP() (string, error) {
	ipAddress := ""
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error: Get current IP address." + err.Error())
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipAddress = ipnet.IP.String()
			}
		}
	}

	return ipAddress, err
}

// If node is master than starts a queue
// Runs node and connects to the queue
func (n *node) Run() error {
	if n.isMaster {
		go n.queue.Run()
	}

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

// Close connection to the master
func (n *node) CloseConn() error {
	err := n.conn.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("[Client] Connection closed.")
	return err
}
