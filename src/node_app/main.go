package main

import (
	"fmt"
	"messaging"
	"bufio"
	"os"
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
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		node.SendMessage(text)
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