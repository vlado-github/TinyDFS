package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vlado-github/tinydfs/logging"
	"github.com/vlado-github/tinydfs/messaging"
	"github.com/vlado-github/tinydfs/utils"

	"github.com/google/uuid"
)

func main() {
	defer close()

	// verbose output of logging to console is enabled
	logging.SetConfiguration(logging.ALL, "../bin/log")

	var params = getParams()
	if isValid(params) {
		// start a node
		printWelcome()
		var n = startNode(params[0], params[1], params[2])
		printInfo(n)
		// run application
		runApp(n)
	}
}

func getParams() []string {
	params := make([]string, 3)
	if len(os.Args) > 1 {
		arg0 := os.Args[1]
		if arg0 == "-help" || arg0 == "-h" {
			printHelp()
		} else {
			arg1 := os.Args[2]
			arg2 := os.Args[3]
			arg3 := os.Args[4]
			arg4 := os.Args[5]
			if (arg0 == "-listen" || arg0 == "-l") && arg1 != "" {
				params[0] = arg1
			}
			if (arg2 == "-connect" || arg2 == "-c") && arg3 != "" && arg4 != "" {
				params[1] = arg3
				params[2] = arg4
			}
		}
	}
	return params
}

func isValid(params []string) bool {
	if len(params) > 0 && params[0] != "" && params[1] != "" && params[2] != "" {
		return true
	}
	return false
}

func startNode(port string, broadcastQueueIP string, broadcastQueuePort string) messaging.Node {
	var deviceIP, err = utils.GetDeviceIpAddress()
	if err != nil {
		logging.AddWarning("Warning: Device IP not found.'")
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
			logging.AddError("Error: Invalid input. Hint: 'sport#We're watching a match.'")
		} else {
			var message = messaging.Message{Key: uuid.New(), Topic: msgArgs[0], Payload: []byte(msgArgs[1])}
			n.SendMessage(message)
		}
	}
}

func close() {
	logging.Close()
}
