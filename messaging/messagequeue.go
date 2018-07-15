package messaging

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"tinylogging"

	"github.com/google/uuid"
)

// MessageQueue is a buffer that receives and broadcasts messages
// between distributed nodes.
type MessageQueue interface {
	Run()
	Status()
	GetNumOfNodes() int
	Close() error

	RegisterHandler(HandlerType, MsgQueueHandlerFunc)
}

type messagequeue struct {
	connParams                    ConnParams
	pool                          Pool
	messageBuffer                 map[string]Message
	onNodeConnectionOpenedHandler MsgQueueHandlerFunc
	onNodeConnectionClosedHandler MsgQueueHandlerFunc
}

var mutex = &sync.Mutex{}

// NewQueue creates new instance of the message queue
func NewQueue(conn ConnParams) MessageQueue {
	return &messagequeue{
		connParams:                    conn,
		onNodeConnectionOpenedHandler: NewMsgQueueHandlerFunc(),
		onNodeConnectionClosedHandler: NewMsgQueueHandlerFunc(),
	}
}

func (queue *messagequeue) init() {
	// Cretes instance of message buffer and connection pool
	queue.messageBuffer = make(map[string]Message)
	queue.pool.conns = make(map[string]net.Conn)
}

// Starts the queue and listens for incoming connections
func (queue *messagequeue) Run() {
	queue.init()

	l, err := net.Listen(queue.connParams.Protocol, queue.connParams.Ip+":"+queue.connParams.Port)
	if err != nil {
		tinylogging.AddError("[Queue] Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	tinylogging.AddInfo("[Queue] Listening on " + queue.connParams.Ip + ":" + queue.connParams.Port)
	fmt.Println("[Queue] Listening on " + queue.connParams.Ip + ":" + queue.connParams.Port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		var poolKey = uuid.New().String()
		queue.pool.conns[poolKey] = conn
		queue.Status()

		if err != nil {
			tinylogging.AddError("[Queue] Error accepting: ", err.Error())
			queue.onCloseConnection(err)
			os.Exit(1)
		}

		queue.onNewConnection()

		var message = Message{Key: uuid.New(), Topic: CONN_ACK}
		encoder := json.NewEncoder(conn)
		encodeMessage(&message, encoder)

		go queue.receiveMessage(conn, poolKey)
		go queue.sendingMessages()
	}
}

// Queue creates system messages to notify other nodes about changes
func (queue *messagequeue) addSystemMessageToQueue(topic string) {
	var key = uuid.New()
	var numOfNodes = queue.GetNumOfNodes()
	payload, err := NewBaseMessagePayload(numOfNodes, queue.getNodesIPAdresses()).ToByteArray()
	if err != nil {
		tinylogging.AddError("[Queue] addSystemMessageToQueue failed to create message payload for topic:", topic)
	} else {
		queue.messageBuffer[key.String()] = Message{Key: key, Topic: topic, Payload: payload}
	}
}

// Handles incoming messages
func (queue *messagequeue) receiveMessage(conn net.Conn, poolKey string) {
	decoder := json.NewDecoder(conn)
	for {
		var message = Message{}
		err := decodeMessage(&message, decoder)
		if err != nil {
			tinylogging.AddInfo("[Queue] Connection closed.")
			conn.Close()
			delete(queue.pool.conns, poolKey)
			queue.Status()
			queue.onCloseConnection(err)
			break
		}

		mutex.Lock()
		var key = uuid.New().String()
		queue.messageBuffer[key] = Message{Key: message.Key, Topic: message.Topic, Payload: message.Payload}
		tinylogging.AddInfo("[Queue] Message Received:", string(queue.messageBuffer[key].Payload))
		mutex.Unlock()
	}
}

// Sends messages from buffer to all clients
func (queue *messagequeue) sendingMessages() {
	for {
		mutex.Lock()
		for index, message := range queue.messageBuffer {
			for _, conn := range queue.pool.conns {
				if conn != nil {
					encoder := json.NewEncoder(conn)
					encodeMessage(&message, encoder)
					tinylogging.AddInfo("[Queue] Sending: ", string(message.Payload)+"\n")
				}
			}
			delete(queue.messageBuffer, index)
		}
		mutex.Unlock()
	}
}

func (queue *messagequeue) Status() {
	fmt.Println("[Queue] Total connections:", len(queue.pool.conns))
}

func (queue *messagequeue) GetNumOfNodes() int {
	return len(queue.pool.conns)
}

func (queue *messagequeue) getNodesIPAdresses() []string {
	ipAddresses := make([]string, len(queue.pool.conns))
	for _, value := range queue.pool.conns {
		ipAddresses = append(ipAddresses, value.RemoteAddr().String())
	}
	return ipAddresses
}

func (queue *messagequeue) Close() error {
	for _, conn := range queue.pool.conns {
		if conn != nil {
			err := conn.Close()
			return err
		}
	}
	return nil
}

func (queue *messagequeue) onNewConnection() {
	tinylogging.AddInfo("[Queue] Client Connected...")
	queue.onNodeConnectionOpenedHandler(queue)
	queue.addSystemMessageToQueue(CLIENT_CONN_OPENED)
}

func (queue *messagequeue) onCloseConnection(err error) {
	tinylogging.AddError("[Queue] Error trace: ", err.Error())
	queue.onNodeConnectionClosedHandler(queue)
	queue.addSystemMessageToQueue(CLIENT_CONN_CLOSED)
}

func (queue *messagequeue) RegisterHandler(handlerType HandlerType, handlerFunc MsgQueueHandlerFunc) {
	switch handlerType {
	case NODECONNOPENED:
		{
			queue.onNodeConnectionOpenedHandler = handlerFunc
			break
		}
	case NODECONNCLOSED:
		{
			queue.onNodeConnectionClosedHandler = handlerFunc
			break
		}
	default:
		{
			break
		}
	}
}
