package main

import (
	"bufio"
	"fmt"
	"messaging"
	"os"
	"strings"

	"strconv"

	"github.com/google/uuid"
)

func main() {
	// if console run follows argument "master"
	var isMaster = false
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "master" {
			isMaster = true
		}
	}

	printWelcome()

	var connParams = messaging.ConnParams{
		Ip:       "localhost",
		Port:     "3333",
		Protocol: "tcp",
	}
	var node = messaging.NewNode(connParams, isMaster)
	node.Run()

	printInfo(node)
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter topic#text:")
		text, _ := reader.ReadString('\n')
		msgArgs := strings.Split(text, "#")
		if len(msgArgs) != 2 {
			fmt.Println("Error: Invalid input. Hint: 'sport#We're watching a match.'")
		} else {
			var message = messaging.Message{Key: uuid.New(), Topic: msgArgs[0], Text: msgArgs[1]}
			node.SendMessage(message)
		}
	}
}

func printWelcome() {
	fmt.Println("*********************")
	fmt.Println("Welcome to TinyDFS")
	fmt.Println("*********************")

}

func printInfo(n messaging.Node) {
	fmt.Println(">>> ID: " + n.GetID().String())
	fmt.Println(">>> Election ID: " + strconv.Itoa(n.GetElectionID()))
	ip, err := n.GetIP()
	if err == nil {
		fmt.Println(">>> IP Address: " + ip)
	}
	fmt.Println("Node started!")
}
