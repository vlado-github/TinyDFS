package messaging

import (
	"encoding/json"
	"logging"

	"net"
	"os"
	"sync"

	"github.com/google/uuid"
)

// MessageQueue is a buffer that receives and broadcasts messages
// between distributed nodes.
type MessageQueue interface {
	Run()
	Status()
	Close() error

	RegisterHandler(HandlerType, MsgQueueHandlerFunc)
}

type messagequeue struct {
	connParams        ConnParams
	pool              Pool
	messageBuffer     map[string]Message
	onMessageReceived MsgQueueHandlerFunc
	networkRegistry   NetworkRegistry
}

var mutex = &sync.Mutex{}

// NewQueue creates new instance of the message queue
func NewQueue(conn ConnParams) MessageQueue {
	return &messagequeue{
		connParams:        conn,
		onMessageReceived: NewMsgQueueHandlerFunc(),
	}
}

// Cretes instance of message buffer, connection pool and network registry
func (queue *messagequeue) init() {
	queue.messageBuffer = make(map[string]Message)
	queue.pool.conns = make(map[string]net.Conn)
	queue.networkRegistry = NewNetworkRegistry()
}

// Starts the queue and listens for incoming connections
func (queue *messagequeue) Run() {
	queue.init()

	l, err := net.Listen(queue.connParams.Protocol, queue.connParams.Ip+":"+queue.connParams.Port)
	if err != nil {
		logging.AddError("[Queue] Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	logging.AddInfo("[Queue] Listening on " + queue.connParams.Ip + ":" + queue.connParams.Port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		var poolKey = uuid.New().String()
		queue.pool.conns[poolKey] = conn
		if err != nil {
			logging.AddError("[Queue] Error accepting: ", err.Error())
			os.Exit(1)
		}
		queue.onNewConnection(conn)

		go queue.receiveMessage(conn, poolKey)
		go queue.sendingMessages()
	}
}

// Handles incoming messages
func (queue *messagequeue) receiveMessage(conn net.Conn, poolKey string) {
	decoder := json.NewDecoder(conn)
	for {
		var message = Message{}
		err := decodeMessage(&message, decoder)
		if err != nil {
			logging.AddInfo("[Queue] Connection closed.")
			conn.Close()
			delete(queue.pool.conns, poolKey)
			queue.removeFromNetworkRegistry(conn)
			break
		}

		if message.Topic == CONN_ACK_REPLY {
			logging.AddInfo("[Queue] Message Received:", message.Topic, string(message.Payload))
			queue.onNewNetworkNode(message)
		} else {
			key := queue.addMessage(message)
			logging.AddInfo("[Queue] Message Received:", string(queue.messageBuffer[key].Payload))
		}
	}
}

// Adds received message to buffer
func (queue *messagequeue) addMessage(message Message) string {
	mutex.Lock()
	var key = uuid.New().String()
	queue.messageBuffer[key] = message
	mutex.Unlock()
	return key
}

// Sends messages from buffer to all nodes
func (queue *messagequeue) sendingMessages() {
	for {
		mutex.Lock()
		for index, message := range queue.messageBuffer {
			for _, conn := range queue.pool.conns {
				if conn != nil {
					encoder := json.NewEncoder(conn)
					encodeMessage(&message, encoder)
					logging.AddInfo("[Queue] Sending: ", string(message.Payload)+"\n")
				}
			}
			delete(queue.messageBuffer, index)
		}
		mutex.Unlock()
	}
}

// Prints current network status
func (queue *messagequeue) Status() {
	logging.AddInfo("[Queue] Total connections:", len(queue.pool.conns))
	networkList, _ := queue.networkRegistry.ToString()
	logging.AddInfo("[Queue] NetworkRegistry: ", networkList)
}

// Closes all connections to nodes
func (queue *messagequeue) Close() error {
	for _, conn := range queue.pool.conns {
		if conn != nil {
			err := conn.Close()
			return err
		}
	}
	return nil
}

// Sends connection ack message to node
func (queue *messagequeue) onNewConnection(conn net.Conn) {
	logging.AddInfo("[Queue] Client Connected...")
	var message = Message{Key: uuid.New(), Topic: CONN_ACK, Payload: []byte(conn.RemoteAddr().String())}
	encoder := json.NewEncoder(conn)
	encodeMessage(&message, encoder)
}

// Adds new node info to network registry
func (queue *messagequeue) onNewNetworkNode(message Message) {
	var networkTuple *networktuple
	err := json.Unmarshal(message.Payload, &networkTuple)
	if err != nil {
		logging.AddError("Message has invalid format.", err.Error())
	}
	queue.networkRegistry.AddItem(networkTuple)
	queue.onNetworkChanged()
}

// Notfies all nodes in network about network change
func (queue *messagequeue) onNetworkChanged() {
	queue.Status()
	payload, err := queue.networkRegistry.ToByteArray()
	if err != nil {
		logging.AddError("Json serialization failed.", err.Error())
	}
	logging.AddInfo(string(payload))
	var message = Message{Key: uuid.New(), Topic: NETWORK_CHANGED, Payload: payload}
	queue.addMessage(message)
}

// Remove closed node from network registry
func (queue *messagequeue) removeFromNetworkRegistry(conn net.Conn) {
	_, port, _ := net.SplitHostPort(conn.RemoteAddr().String())
	networkItem, index := queue.networkRegistry.GetItemByRemoteAddPort(port)
	if networkItem != nil {
		queue.networkRegistry.RemoveItem(index)
	}
	queue.onNetworkChanged()
}

func (queue *messagequeue) RegisterHandler(handlerType HandlerType, handlerFunc MsgQueueHandlerFunc) {
	switch handlerType {
	default:
		{
			break
		}
	}
}
