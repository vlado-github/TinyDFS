package main

import (
	"fmt"
	"strconv"

	"github.com/vlado-github/tinydfs/messaging"
)

func printWelcome() {
	fmt.Println("")
	fmt.Println("***********************************************");
	fmt.Println("***********************************************");
	fmt.Println(" _____  _               ______ ______  _____   ");
	fmt.Println("|_   _|(_)              |  _  \|  ___|/  ___|  ");
	fmt.Println("  | |   _  _ __   _   _ | | | || |_   \ `--.   ");
	fmt.Println("  | |  | || '_ \ | | | || | | ||  _|   `--. \  ");
	fmt.Println("  | |  | || | | || |_| || |/ / | |    /\__/ /  ");
	fmt.Println("  \_/  |_||_| |_| \__, ||___/  \_|    \____/   ");
	fmt.Println("                   __/ |                       ");
	fmt.Println("                  |___/                        ");
	fmt.Println("***********************************************");
	fmt.Println("***********************************************");
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
