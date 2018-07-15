package main

import (
	"bufio"
	"fmt"
	"messaging"
	"os"
	"strings"
	"tinydfs"
	"tinylogging"

	"github.com/google/uuid"
)

func main() {
	defer close()

	// verbose output of logging to console is enabled
	// log directory specified
	tinylogging.SetConfiguration(tinylogging.TRACE, "../bin/log")

	// if console run command follows argument "-p 3334 -queue 192.168.1.55 3333"
	var params = getParams()

	// start a node and display info
	printWelcome()
	var host = startHost(params[0], params[1], params[2])
	printInfo(host)

	// actual implementation of node usage
	runApp(host)
}

func getParams() []string {
	params := make([]string, 3)
	fmt.Println(len(os.Args))
	if len(os.Args) > 1 {
		arg0 := os.Args[1]
		arg1 := os.Args[2]
		arg2 := os.Args[3]
		arg3 := os.Args[4]
		arg4 := os.Args[5]
		if arg0 == "-p" && arg1 != "" {
			params[0] = arg1
		}
		if arg2 == "-queue" && arg3 != "" && arg4 != "" {
			params[1] = arg3
			params[2] = arg4
		}
	}
	return params
}

func startHost(port string, broadcastQueueIP string, broadcastQueuePort string) tinydfs.Host {
	var connParams = messaging.ConnParams{
		Ip:       "localhost",
		Port:     port,
		Protocol: "tcp",
	}
	tinylogging.AddTrace(broadcastQueueIP, broadcastQueuePort)
	var broadcastConnParams = messaging.ConnParams{
		Ip:       broadcastQueueIP,
		Port:     broadcastQueuePort,
		Protocol: "tcp",
	}

	var host = tinydfs.NewHost(connParams, broadcastConnParams, true, port)
	onConnClosedCallback := func() {
		var message = messaging.Message{Key: uuid.New(), Topic: "ConnClose for node: " + host.GetID().String()}
		host.SendMessage(message)
	}
	onQueueConnClosedCallback := func() {
		fmt.Println("TODO: implement on queue connection closed handler...")
	}
	host.RegisterNodeHandler(messaging.QUEUECONNCLOSED, onQueueConnClosedCallback)
	host.RegisterNodeHandler(messaging.NODECONNCLOSED, onConnClosedCallback)
	host.Start()
	return host
}

func runApp(host tinydfs.Host) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter topic#text:")
		text, _ := reader.ReadString('\n')
		msgArgs := strings.Split(text, "#")
		if len(msgArgs) != 2 {
			fmt.Println("Error: Invalid input. Hint: 'sport#We're watching a match.'")
		} else {
			var message = messaging.Message{Key: uuid.New(), Topic: msgArgs[0], Payload: []byte(msgArgs[1])}
			host.SendMessage(message)
		}
	}
}

func close() {
	tinylogging.Close()
}
