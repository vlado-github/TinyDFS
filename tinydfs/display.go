package main

import (
	"fmt"
	"strconv"
)

func printWelcome() {
	fmt.Println("")
	fmt.Println("	**************************	")
	fmt.Println("	*** Welcome to TinyDFS ***	")
	fmt.Println("	**************************	")
	fmt.Println("")
}

func printInfo(host Host) {
	fmt.Println(">>> ID: " + host.GetID().String())
	fmt.Println(">>> Election ID: " + strconv.Itoa(host.GetElectionID()))
	ip, err := host.GetIP()
	if err == nil {
		fmt.Println(">>> IP Address: " + ip)
	}
}
