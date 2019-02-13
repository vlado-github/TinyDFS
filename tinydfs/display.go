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
