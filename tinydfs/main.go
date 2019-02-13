package main

import (
	"bufio"
	"fmt"
	"logging"
	"messaging"
	"os"
	"strings"
	"utils"

	"github.com/google/uuid"
)

func main() {
	defer close()

	// verbose output of logging to console is enabled
	logging.SetConfiguration(logging.ALL, "../bin/log")

	// if console run command follows argument "master"
	var params = getParams()

	// start a node and display info
	printWelcome()
	var n = startNode(params[0], params[1], params[2])
	printInfo(n)

	// actual implementation of node usage
	runApp(n)
}

func getParams() []string {
	params := make([]string, 3)
	if len(os.Args) > 1 {
		arg0 := os.Args[1]
		arg1 := os.Args[2]
		arg2 := os.Args[3]
		arg3 := os.Args[4]
		arg4 := os.Args[5]
		if (arg0 == "-listen" || arg0 == "-l") && arg1 != "" {
			params[0] = arg1
		}
		if (arg2 == "-broadcast" || arg2 == "-b") && arg3 != "" && arg4 != "" {
			params[1] = arg3
			params[2] = arg4
		}
	}
	return params
}

func startNode(port string, broadcastQueueIP string, broadcastQueuePort string) messaging.Node {
	var deviceIP, err = utils.GetDeviceIpAddress()
	if err != nil {
		fmt.Println("Warning: Device IP not found.'")
		deviceIP = "localhost"
	}
	var connParams = messaging.ConnParams{
		Ip:       deviceIP,
		Port:     port,
		Protocol: "tcp",
	}
	logging.AddTrace(broadcastQueueIP, broadcastQueuePort)
	var broadcastConnParams = messaging.ConnParams{
		Ip:       broadcastQueueIP,
		Port:     broadcastQueuePort,
		Protocol: "tcp",
	}

	var n = messaging.NewNode(connParams, broadcastConnParams, true)
	onConnClosedCallback := func(nn messaging.Node) {
		var message = messaging.Message{Key: uuid.New(), Topic: "ConnClose for node: " + n.GetID().String(), Text: "Goodbye!"}
		nn.SendMessage(message)
	}
	n.RegisterHandler(messaging.NODECONNCLOSED, onConnClosedCallback)
	n.Run()
	return n
}

func runApp(n messaging.Node) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter topic#text:")
		text, _ := reader.ReadString('\n')
		msgArgs := strings.Split(text, "#")
		if len(msgArgs) != 2 {
			fmt.Println("Error: Invalid input. Hint: 'sport#We're watching a match.'")
		} else {
			var message = messaging.Message{Key: uuid.New(), Topic: msgArgs[0], Text: msgArgs[1]}
			n.SendMessage(message)
		}
	}
}

func close() {
	logging.Close()
}
