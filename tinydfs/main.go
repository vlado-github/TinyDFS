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
	var isMaster = getArgs()

	// start a node and display info
	printWelcome()
	var n = startNode(isMaster)
	printInfo(n)

	// actual implementation of node usage
	runApp(n)
}

func getArgs() bool {
	var isMaster = false
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "master" {
			isMaster = true
		}
	}
	return isMaster
}

func startNode(isMaster bool) messaging.Node {
	var connParams = messaging.ConnParams{
		Ip:       "localhost",
		Port:     "3333",
		Protocol: "tcp",
	}
	var n = messaging.NewNode(connParams, isMaster)
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
