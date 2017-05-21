package main

import (
	"fmt"
	"messaging"
	"bufio"
	"os"
	"github.com/google/uuid"
	"strings"
)

func main() {
	printWelcome()
	var connParams = messaging.ConnParams{
		Ip:"localhost",
		Port:"3333",
		Protocol:"tcp",
	}
	var node = messaging.NewNode(connParams)
	node.Run()
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter topic#text:")
		text, _ := reader.ReadString('\n')
		msgArgs := strings.Split(text,"#")
		if len(msgArgs) != 2 {
			fmt.Println("Error: Invalid input. Hint 'sport#We're watching a match.'")
		} else {
			var message = messaging.Message{Key:uuid.New(), Topic:msgArgs[0], Text:msgArgs[1]}
			node.SendMessage(message)
		}
	}
}

func printWelcome(){
	fmt.Println("*********************")
	fmt.Println("Welcome to TinyDFS")
	fmt.Println("*********************")
	fmt.Println("\n")
	fmt.Println("Node started!")
	fmt.Println("\n")
}