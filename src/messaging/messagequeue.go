package messaging

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"log"
	"github.com/google/uuid"
	"sync"
)

type MessageQueue interface{
	Run()
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

// Starts the queue and listens for incoming connections
func (queue *messagequeue) Run() {
	// Cretes instance of message buffer
	queue.messageBuffer = make(map[string]Message)

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
		queue.pool.conns = append(queue.pool.conns, conn)
		fmt.Println("Client Connected...Total:",len(queue.pool.conns))
		if err != nil {
			fmt.Println("[Queue] Error accepting: ", err.Error())
			log.Fatal(err)
			os.Exit(1)
		}
		fmt.Println("[Queue] Client Connected...")
		conn.Write([]byte("CONN_ACK\n"))
		go queue.receiveMessage(conn)
		go queue.sendingMessages()
	}
}

// Handles incoming messages
func (queue *messagequeue) receiveMessage(conn net.Conn) {
	for {
		scanner := bufio.NewScanner(conn)
		if ok := scanner.Scan(); ok {
			message := scanner.Text()
			key := uuid.New().String()

			mutex.Lock()
			queue.messageBuffer[key] = Message{message}
			fmt.Print("[Queue] Message Received:", queue.messageBuffer[key].text+"\n")
			mutex.Unlock()
		} else {
			fmt.Print("[Queue] Connection closed.")
			conn.Close()
			break
		}
	}
}

// Sends messages from buffer to all clientst
func (queue *messagequeue) sendingMessages(){
	for {
		mutex.Lock()
		for index,message := range queue.messageBuffer {
			for _,conn := range queue.pool.conns {
				// sends message to each client
				newmessage := strings.ToUpper(message.text)
				fmt.Print("[Queue] Sending: ", newmessage+"\n")
				_, err := conn.Write([]byte(newmessage + "\n"))
				if err != nil {
					fmt.Print(err.Error())
				}
			}
			delete(queue.messageBuffer, index)
		}
		mutex.Unlock()
	}
}

