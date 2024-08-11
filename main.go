package main

import (
	"fmt"
	"os"
)

func main() {
	// parse args
	if len(os.Args) < 2 {
		fmt.Println(USAGE)
		return
	}
	switch os.Args[1] {
	case "sync":
		syncLyrics(os.Args[2:])
	case "help":
		fmt.Println(USAGE)
	default:
		fmt.Println("Unknown command")
		fmt.Println(USAGE)
	}
}
