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
type Pool struct {
	conns []net.Conn
}

var pool = Pool{}

var messageBuffer map[string]Message
var mutex = &sync.Mutex{}

func InitQueue(params ConnParams) {
	// Crete message buffer
	messageBuffer = make(map[string]Message)

	// Listen for incoming connections.
	l, err := net.Listen(params.protocol, params.ip+":"+params.port)
	if err != nil {
		fmt.Println("[Queue] Error listening:", err.Error())
		log.Fatal(err)
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("[Queue] Listening on " + params.ip + ":" + params.port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		pool.conns = append(pool.conns, conn)
		fmt.Println("Client Connected...Total:",len(pool.conns))
		if err != nil {
			fmt.Println("[Queue] Error accepting: ", err.Error())
			log.Fatal(err)
			os.Exit(1)
		}
		fmt.Println("[Queue] Client Connected...")
		conn.Write([]byte("CONN_ACK\n"))
		go receiveMessage(conn)
		go sendingMessages()
	}
}

// Handles incoming messages.
func receiveMessage(conn net.Conn) {
	for {
		scanner := bufio.NewScanner(conn)
		if ok := scanner.Scan(); ok {
			message := scanner.Text()
			key := uuid.New().String()

			mutex.Lock()
			messageBuffer[key] = Message{message}
			fmt.Print("[Queue] Message Received:", messageBuffer[key].text+"\n")
			mutex.Unlock()
		} else {
			fmt.Print("[Queue] Connection closed.")
			conn.Close()
			break
		}
	}
}

// Sending messages from buffer as broadcast
func sendingMessages(){
	for {
		mutex.Lock()
		for index,message := range messageBuffer {
			for _,conn := range pool.conns {
				// sends message to each client
				newmessage := strings.ToUpper(message.text)
				fmt.Print("[Queue] Sending: ", newmessage+"\n")
				_, err := conn.Write([]byte(newmessage + "\n"))
				if err != nil {
					fmt.Print(err.Error())
				}
			}
			delete(messageBuffer, index)
		}
		mutex.Unlock()
	}
}

