package messaging

import (
	"fmt"
	"net"
	"os"
	"log"
	"github.com/google/uuid"
	"sync"
	"encoding/json"
)

type MessageQueue interface{
	Run()
	Status()
}

type messagequeue struct {
	connParams ConnParams
	pool Pool
	messageBuffer map[string]Message
}

var mutex = &sync.Mutex{}

// Creates new instance of the message queue
func NewQueue(conn ConnParams) MessageQueue{
	return &messagequeue{
		connParams: conn,
	}
}

func (queue *messagequeue) init() {
	// Cretes instance of message buffer and connection pool
	queue.messageBuffer = make(map[string]Message)
	queue.pool.conns = make(map[string]net.Conn)
}

// Starts the queue and listens for incoming connections
func (queue *messagequeue) Run() {
	queue.init();

	l, err := net.Listen(queue.connParams.Protocol, queue.connParams.Ip+":"+queue.connParams.Port)
	if err != nil {
		fmt.Println("[Queue] Error listening:", err.Error())
		log.Fatal(err)
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("[Queue] Listening on " + queue.connParams.Ip + ":" + queue.connParams.Port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		var poolKey = uuid.New().String()
		queue.pool.conns[poolKey] = conn
		queue.Status()
		if err != nil {
			fmt.Println("[Queue] Error accepting: ", err.Error())
			log.Fatal(err)
			os.Exit(1)
		}
		fmt.Println("[Queue] Client Connected...")

		var message = Message{Key:uuid.New(), Topic:"CONN_ACK", Text:"CONN_ACK"}
		encoder := json.NewEncoder(conn)
		encodeMessage(&message, encoder)

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
			fmt.Println("[Queue] Connection closed.")
			conn.Close()
			delete(queue.pool.conns, poolKey)
			queue.Status()
			break
		}

		mutex.Lock()
		var key = message.Key.String()
		queue.messageBuffer[key] = Message{Key: message.Key, Topic: message.Topic, Text: message.Text}
		fmt.Println("[Queue] Message Received:", queue.messageBuffer[key].Text)
		mutex.Unlock()
	}
}

// Sends messages from buffer to all clients
func (queue *messagequeue) sendingMessages(){
	for {
		mutex.Lock()
		for index,message := range queue.messageBuffer {
			for _,conn := range queue.pool.conns {
				if conn != nil {
					encoder := json.NewEncoder(conn)
					encodeMessage(&message, encoder)
					fmt.Print("[Queue] Sending: ", message.Text + "\n")
				}
			}
			delete(queue.messageBuffer, index)
		}
		mutex.Unlock()
	}
}

func (queue *messagequeue) Status(){
	fmt.Println("[Queue] Total connections:",len(queue.pool.conns))
}

