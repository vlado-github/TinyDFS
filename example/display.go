package main

import (
	"fmt"
	"strconv"
	"tinydfs"
)

func printWelcome() {
	fmt.Println("")
	fmt.Println("	*******************************	")
	fmt.Println("	*** Welcome to TinyDFS_Chat ***	")
	fmt.Println("	*******************************	")
	fmt.Println("")
}

func printInfo(host tinydfs.Host) {
	fmt.Println(">>> ID: " + host.GetID().String())
	fmt.Println(">>> Election ID: " + strconv.Itoa(host.GetElectionID()))
	ip, err := host.GetIP()
	if err == nil {
		fmt.Println(">>> IP Address: " + ip)
	}
}
