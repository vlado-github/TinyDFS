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
	GetExchangeNetworkTuple() NetworkTuple

	RegisterNodeHandler(HandlerType, NodeHandlerFunc)
	RegisterQueueHandler(HandlerType, MsgQueueHandlerFunc)
}

type node struct {
	id                        uuid.UUID
	electionID                int
	remoteAddressPort         string
	broadcastQueueConnParams  ConnParams
	exchangeQueueConnParams   ConnParams
	conn                      net.Conn
	fileManager               persistance.FileManager
	queue                     MessageQueue
	onConnectionClosedHandler NodeHandlerFunc
	onConnectionOpenedHandler NodeHandlerFunc
	persistanceEnabled        bool
	networkRegistry           NetworkRegistry
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
		onConnectionOpenedHandler: NewHandlerFunc(),
		networkRegistry:           NewNetworkRegistry(),
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

// Returns the Node exchange ip and port
func (n *node) GetExchangeNetworkTuple() NetworkTuple {
	return NewNetworkTuple(n.GetID().String(), n.exchangeQueueConnParams.Ip, n.exchangeQueueConnParams.Port)
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
				n.onConnectionOpenedHandler(n)
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
			logging.AddInfo("[Client] Received: ", message.Topic, string(message.Payload))
			if message.Topic == CONN_ACK {
				n.onConnectionAcknowledged(message)
			} else if message.Topic == NETWORK_CHANGED {

			} else {
				var guid = uuid.New()
				var cmd = persistance.Command{Key: guid, Text: string(message.Payload), Topic: message.Topic}
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

// After connection is ack from queue side, node sends Id details to queue
// to update network registry.
func (n *node) onConnectionAcknowledged(message Message) {
	if message.Topic != CONN_ACK {
		return
	}
	ip, port, _ := net.SplitHostPort(string(message.Payload))
	logging.AddInfo("[Client] Connected.", ip, port)
	n.remoteAddressPort = port
	networkTuple := NewNetworkTuple(n.GetID().String(), ip, port)
	payload, err := json.Marshal(networkTuple)
	if err != nil {
		logging.AddError("Json serialization failed.", err)
	}
	ackReply := Message{Key: uuid.New(), Topic: CONN_ACK_REPLY, Payload: []byte(payload)}
	n.SendMessage(ackReply)
}

// Queue notifies nodes about network updates
func (n *node) onNetworkChanged(message Message) {
	if message.Topic != NETWORK_CHANGED {
		return
	}
	err := n.networkRegistry.FromByteArray(message.Payload)
	if err != nil {
		logging.AddError("OnNetworkChanged invalid message format.", err.Error())
	}
}

func (n *node) RegisterNodeHandler(handlerType HandlerType, handlerFunc NodeHandlerFunc) {
	switch handlerType {
	case NODECONNCLOSED:
		{
			n.onConnectionClosedHandler = handlerFunc
			break
		}
	case NODECONNOPENED:
		{
			n.onConnectionOpenedHandler = handlerFunc
		}
	default:
		{
			break
		}
	}
}

func (n *node) RegisterQueueHandler(handlerType HandlerType, handlerFunc MsgQueueHandlerFunc) {
	switch handlerType {
	case NETWORKCHANGED:
		{
			n.queue.RegisterHandler(handlerType, handlerFunc)
			break
		}
	default:
		{
			break
		}
	}
}
