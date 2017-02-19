package messaging

import (
	"bufio"
	"fmt"
	"net"
)

var conn net.Conn;

func InitNode(params connParams) error {
	// connect to queue
	var err error
	conn, err = net.Dial(params.protocol, params.ip + ":" + params.port)

	if err != nil {
		fmt.Println("Error dialing:", err.Error())
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

func CloseConn(){
	conn.Close();
}


