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
	tinylogging.SetConfiguration(false, "../bin/log")

	// if console run command follows argument "queue"
	var isQueue = getArgs()

	// start a node and display info
	printWelcome()
	var host = startHost(isQueue)
	printInfo(host)

	// actual implementation of node usage
	runApp(host)
}

func getArgs() bool {
	var isQueue = false
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "queue" {
			isQueue = true
		}
	}
	return isQueue
}

func startHost(isQueue bool) tinydfs.Host {
	var connParams = messaging.ConnParams{
		Ip:       "localhost",
		Port:     "3333",
		Protocol: "tcp",
	}

	var host = tinydfs.NewHost(connParams, isQueue)
	onConnClosedCallback := func() {
		var message = messaging.Message{Key: uuid.New(), Topic: "ConnClose for node: " + host.GetID().String()}
		host.SendMessage(message)
	}
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
