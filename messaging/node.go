package messaging

import (
	"encoding/json"
	"logging"
	"net"
	"persistance"
	"time"

	"math/rand"

	"github.com/google/uuid"
)

// Node is a single unit of distributed storage,
// that communicates via the master node (i.e. queue).
type Node interface {
	Run() error
	SendMessage(message Message)
	ConnectToQueue(protocol string, address string) error
	CloseConn() error

	GetID() uuid.UUID
	GetElectionID() int

	RegisterHandler(HandlerType, NodeHandlerFunc)
}

type node struct {
	id                        uuid.UUID
	electionID                int
	broadcastQueueConnParams  ConnParams
	exchangeQueueConnParams   ConnParams
	conn                      net.Conn
	fileManager               persistance.FileManager
	queue                     MessageQueue
	onConnectionClosedHandler NodeHandlerFunc
	persistanceEnabled        bool
}

const MaxNumberOfConnAttempts int = 10

// NewNode creates new instance of node
func NewNode(exchangeQueueConn ConnParams, broadcastQueueConn ConnParams, persistanceEnabled bool) Node {
	rand.Seed(time.Now().Unix())
	uniqueID := uuid.New()
	randomID := rand.Int()
	fm := persistance.NewFileManager(getCurrentDirectory() + "//" + uniqueID.String())
	msgQueue := NewQueue(exchangeQueueConn)

	return &node{
		id:                       uniqueID,
		electionID:               randomID,
		exchangeQueueConnParams:  exchangeQueueConn,
		fileManager:              fm,
		broadcastQueueConnParams: broadcastQueueConn,
		persistanceEnabled:       persistanceEnabled,
		queue:                    msgQueue,
		onConnectionClosedHandler: NewHandlerFunc(),
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

// If node is master than starts a queue
// Runs node and connects to the queue
func (n *node) Run() error {
	// run exchange queue
	go n.queue.Run()

	// connects to broadcast queue
	return n.ConnectToQueue(n.broadcastQueueConnParams.Protocol, n.broadcastQueueConnParams.Ip+":"+n.broadcastQueueConnParams.Port)
}

// Connects to queue
func (n *node) ConnectToQueue(protocol string, address string) error {
	if n.conn != nil {
		n.conn.Close()
	}
	numOfAttempts := 0
	var err error
	n.conn, err = net.Dial(protocol, address)
	numOfAttempts++

	if err != nil {
		var isConnected = false
		for numOfAttempts <= MaxNumberOfConnAttempts {
			n.conn, err = net.Dial(protocol, address)
			if err == nil {
				isConnected = true
				break
			}
			numOfAttempts++
		}
		if !isConnected {
			logging.AddError("[Node] Error dialing: ", address, protocol, err.Error(), numOfAttempts, " attempts.")
			return err
		}
	}

	go n.receiveMessages()

	return err
}

// Sends message to the queue
func (n *node) SendMessage(message Message) {
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
			logging.AddError("Error: Queue connection is closed.", err.Error())
			break
		} else {
			logging.AddInfo("[Client] Received: ", message.Topic, message.Text)
			if message.Text == "CONN_ACK" {
				logging.AddInfo("[Client] Connected")
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
	n.onConnectionClosedHandler(n)
	err := n.conn.Close()
	if err != nil {
		logging.AddError("Close connection on node failed.", err.Error())
		return err
	}
	if n.queue != nil {
		err := n.queue.Close()
		logging.AddError("Close message queue connection on node failed.", err.Error())
		return err
	}
	return err
}

func (n *node) RegisterHandler(handlerType HandlerType, handlerFunc NodeHandlerFunc) {
	switch handlerType {
	case NODECONNCLOSED:
		{
			n.onConnectionClosedHandler = handlerFunc
			break
		}
	default:
		{
			break
		}
	}
}
