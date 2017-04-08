package main

import (
	"messaging"
	"fmt"
)

func main() {
	printWelcome()
	var connParams = messaging.ConnParams{
		Ip:"localhost",
		Port:"3333",
		Protocol:"tcp",
	}
	var queue = messaging.NewQueue(connParams)
	queue.Run()
}

func printWelcome(){
	fmt.Println("*********************")
	fmt.Println("Welcome to TinyDFS")
	fmt.Println("*********************")
	fmt.Println("\n")
	fmt.Println("Message queue started!")
	fmt.Println("\n")
}
