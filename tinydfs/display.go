package main

import (
	"fmt"
	"messaging"
	"strconv"
)

func printWelcome() {
	fmt.Println("")
	fmt.Println("	**************************	")
	fmt.Println("	*** Welcome to TinyDFS ***	")
	fmt.Println("	**************************	")
	fmt.Println("")
}

func printInfo(n messaging.Node) {
	fmt.Println(">>> ID: " + n.GetID().String())
	fmt.Println(">>> Election ID: " + strconv.Itoa(n.GetElectionID()))
}

func printHelp() {
	fmt.Println("-listen or -l This arg is required, followed by port number for exchange queue")
	fmt.Println("-connect or -c This arg is required, followed by IP and port of broadcast queue")
}
