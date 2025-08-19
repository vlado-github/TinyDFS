package main

import (
	"fmt"
	"strconv"

	"github.com/vlado-github/tinydfs/messaging"
)

func printWelcome() {
	fmt.Println("")
	fmt.Println(`
	*******************************************
	*******************************************
	
	  _____ _               ____  _____ ____  
	 |_   _(_)_ __  _   _  |  _ \|  ___/ ___| 
	   | | | | '_ \| | | | | | | | |_  \___ \  
	   | | | | | | | |_| | | |_| |  _|  ___) | 
	   |_| |_|_| |_|\__, | |____/|_|   |____/  
	                |___/                    

    
	*******************************************
	*******************************************`);
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
