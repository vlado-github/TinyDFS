package messaging

import (
	"bufio"
	"fmt"
	"net"
	"log"
)

var conn net.Conn;

func InitNode(params ConnParams) error {
	// connect to queue
	var err error
	conn, err = net.Dial(params.protocol, params.ip + ":" + params.port)

	if err != nil {
		fmt.Println("Error dialing:", err.Error())
		log.Fatal(err)
	} else {
		for {
			ack_msg := ReceiveEvent()
			if ack_msg == "CONN_ACK\n" {
				fmt.Println("[Client] Connected")
				break;
			}
		}
	}
	return err
}

// sends message to queue
func SendEvent(message string){
	fmt.Println("[Client] Sending: ", message)
	fmt.Fprintf(conn, message + "\n")
}

// receives message from queue
func ReceiveEvent() string{
	event, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("[Client] Received: ", event)
	return  event;
}

func CloseConn() error{
	err := conn.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("[Client] Connection closed.")
	return err
}

