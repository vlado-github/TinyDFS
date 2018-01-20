package main

import (
	"bufio"
	"fmt"
	"logging"
	"messaging"
	"os"
	"strings"

	"github.com/google/uuid"
)

func main() {
	defer close()

	// verbose output of logging to console is enabled
	logging.SetVerbose(true)

	// if console run command follows argument "master"
	var isLead = getArgs()

	// start a node and display info
	printWelcome()
	var host = startNode(isLead)
	printInfo(host)

	// actual implementation of node usage
	runApp(host)
}

func getArgs() bool {
	var isLead = false
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "leader" {
			isLead = true
		}
	}
	return isLead
}

func startNode(isLead bool) Host {
	var connParams = messaging.ConnParams{
		Ip:       "localhost",
		Port:     "3333",
		Protocol: "tcp",
	}
	var host = NewHost(connParams, isLead)
	onConnClosedCallback := func(nn messaging.Node) {
		var message = messaging.Message{Key: uuid.New(), Topic: "ConnClose for node: " + host.GetID().String(), Text: "Goodbye!"}
		nn.SendMessage(message)
	}
	host.RegisterNodeHandler(messaging.NODECONNCLOSED, onConnClosedCallback)
	host.Start()
	return host
}

func runApp(host Host) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter topic#text:")
		text, _ := reader.ReadString('\n')
		msgArgs := strings.Split(text, "#")
		if len(msgArgs) != 2 {
			fmt.Println("Error: Invalid input. Hint: 'sport#We're watching a match.'")
		} else {
			var message = messaging.Message{Key: uuid.New(), Topic: msgArgs[0], Text: msgArgs[1]}
			host.SendMessage(message)
		}
	}
}

func close() {
	logging.Close()
}
