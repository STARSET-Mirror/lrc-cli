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
	case "convert":
		convertLyrics(os.Args[2:])
	case "fmt":
		fmtLyrics(os.Args[2:])
	case "version":
		fmt.Printf("lyc-cli version %s\n", VERSION)
	case "help":
		fmt.Println(USAGE)
	default:
		fmt.Println("Unknown command")
		fmt.Println(USAGE)
	}
}
