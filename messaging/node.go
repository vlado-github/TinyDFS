package messaging

import (
	"encoding/json"
	"net"
	"persistance"
	"time"
	"tinylogging"

	"math/rand"

	"github.com/google/uuid"
)

// Node is a single unit for messaging,
// that communicates by sending messages over queue node (i.e. queue).
// Node also saves all received and sent messages.
type Node interface {
	Run() error
	SendMessage(message Message)
	ConnectToQueue(protocol string, address string) error
	IsConnectedToQueue() bool
	CloseConn() error

	GetID() uuid.UUID
	GetIP() (string, error)
	GetPort() string
	GetNumOfNodes() int

	RegisterHandler(HandlerType, NodeHandlerFunc)
	RegisterMessageHandler(HandlerType, MessageHandlerFunc)
}

type node struct {
	id                             uuid.UUID
	connParams                     ConnParams
	broadcastQueueConnParams       ConnParams
	conn                           net.Conn
	fileManager                    persistance.FileManager
	queue                          MessageQueue
	onConnectionClosedHandler      NodeHandlerFunc
	onQueueConnectionClosedHandler NodeHandlerFunc
	onMessageReceivedHandler       MessageHandlerFunc
	port                           string
}

const DefaultPort string = "3333"
const MaxNumberOfConnAttempts int = 10

// NewNode creates new instance of node
func NewNode(conn ConnParams, broadcastQueueConn ConnParams, port string) Node {
	rand.Seed(time.Now().Unix())
	uniqueID := uuid.New()
	fm := persistance.NewFileManager(getCurrentDirectory() + "//" + uniqueID.String())
	var msgQueue = NewQueue(conn)

	var p string
	if port == "" {
		p = DefaultPort
	}

	return &node{
		id:                       uniqueID,
		connParams:               conn,
		broadcastQueueConnParams: broadcastQueueConn,
		fileManager:              fm,
		queue:                    msgQueue,
		onQueueConnectionClosedHandler: NewNodeHandlerFunc(),
		onConnectionClosedHandler:      NewNodeHandlerFunc(),
		onMessageReceivedHandler:       NewMessageHandlerFunc(),
		port: p,
	}
}

// Returns the Node unique ID
func (n *node) GetID() uuid.UUID {
	return n.id
}

// Return current IP address of the device
func (n *node) GetIP() (string, error) {
	ipAddress := ""
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		tinylogging.AddError("Retrieving host's IP address failed.", err.Error())
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

func (n *node) GetPort() string {
	return n.port
}

// Return number of connected nodes to the queue
func (n *node) GetNumOfNodes() int {
	return n.queue.GetNumOfNodes()
}

// If node is queue than starts its own queue
// Runs node and connects to the broadcast queue
func (n *node) Run() error {
	// starts queue
	go n.queue.Run()
	// connects to broadcast queue
	return n.ConnectToQueue(n.broadcastQueueConnParams.Protocol, n.broadcastQueueConnParams.Ip+":"+n.broadcastQueueConnParams.Port)
}

// Connect to queue
func (n *node) ConnectToQueue(protocol string, address string) error {
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
			tinylogging.AddError("Error dialing: ", address, protocol, err.Error(), numOfAttempts, " attempts.")
			return err
		}
	}

	go n.receiveMessages()

	return err
}

func (n *node) IsConnectedToQueue() bool {
	return n.conn != nil
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
			tinylogging.AddError("Error: Decoding received message.", err.Error())
			break
		} else {
			tinylogging.AddInfo("[Client] Received: ", message.Topic, string(message.Payload))
			if message.Topic == CONN_ACK {
				tinylogging.AddInfo("[Client] Connected")
			} else {
				var guid = uuid.New()
				var cmd = persistance.Command{Key: guid, Text: string(message.Payload), Topic: message.Topic}
				n.fileManager.Write(cmd)
				n.onMessageReceivedHandler(message)
			}
		}
	}
	n.onQueueConnectionClosedHandler()
}

// Close connection to the queue
func (n *node) CloseConn() error {
	err := n.conn.Close()
	if err != nil {
		tinylogging.AddError("Close connection on node failed.", err.Error())
		return err
	}
	if n.queue != nil {
		err := n.queue.Close()
		if err != nil {
			tinylogging.AddError("Close message queue connection on node failed.", err.Error())
		}
		return err
	}
	n.onConnectionClosedHandler()
	return err
}

func (n *node) RegisterHandler(handlerType HandlerType, handlerFunc NodeHandlerFunc) {
	switch handlerType {
	case NODECONNCLOSED:
		{
			n.onConnectionClosedHandler = handlerFunc
			break
		}
	case QUEUECONNCLOSED:
		{
			n.onQueueConnectionClosedHandler = handlerFunc
			break
		}
	default:
		{
			break
		}
	}
}

func (n *node) RegisterMessageHandler(handlerType HandlerType, handlerFunc MessageHandlerFunc) {
	switch handlerType {
	case MESSAGERECEIVED:
		{
			n.onMessageReceivedHandler = handlerFunc
			break
		}
	default:
		{
			break
		}
	}
}
