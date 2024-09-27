package main

import (
	"fmt"
	"os"
)

func lrcToTxt(sourceFile, targetFile string) {
	sourceLyrics, err := parseLyrics(sourceFile)
	if err != nil {
		fmt.Println("Error parsing source lyrics file:", err)
		return
	}

	file, err := os.Create(targetFile)
	if err != nil {
		fmt.Println("Error creating target file:", err)
		return
	}
	defer file.Close()

	for _, content := range sourceLyrics.Content {
		fmt.Fprintln(file, content)
	}
}
