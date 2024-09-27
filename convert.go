package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func convert(args []string) {
	if len(args) < 2 {
		fmt.Println(CONVERT_USAGE)
		return
	}

	sourceFile := args[0]
	targetFile := args[1]

	sourceFmt := strings.TrimPrefix(filepath.Ext(sourceFile), ".")
	targetFmt := strings.TrimPrefix(filepath.Ext(targetFile), ".")

	switch sourceFmt {
	case "lrc":
		convertLyrics(sourceFile, targetFile, targetFmt)
	case "srt":
		convertSRT(sourceFile, targetFile, targetFmt)
	default:
		fmt.Printf("unsupported source file format: %s\n", sourceFmt)
	}
}

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

func lrcToSrt(sourceFile, targetFile string) {
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

	for i, content := range sourceLyrics.Content {
		startTime := sourceLyrics.Timeline[i]
		var endTime Timestamp
		if i < len(sourceLyrics.Timeline)-1 {
			endTime = sourceLyrics.Timeline[i+1]
		} else {
			endTime = addSeconds(startTime, 3)
		}

		fmt.Fprintf(file, "%d\n", i+1)
		fmt.Fprintf(file, "%s --> %s\n", formatSRTTimestamp(startTime), formatSRTTimestamp(endTime))
		fmt.Fprintf(file, "%s\n\n", content)
	}
}

func srtToLrc(sourceFile, targetFile string) {
	srtEntries, err := parseSRT(sourceFile)
	if err != nil {
		fmt.Println("Error parsing source SRT file:", err)
		return
	}

	lyrics := Lyrics{
		Metadata: make(map[string]string),
		Timeline: make([]Timestamp, len(srtEntries)),
		Content:  make([]string, len(srtEntries)),
	}

	for i, entry := range srtEntries {
		lyrics.Timeline[i] = entry.StartTime
		lyrics.Content[i] = entry.Content
	}

	err = saveLyrics(targetFile, lyrics)
	if err != nil {
		fmt.Println("Error saving LRC file:", err)
		return
	}
}

func srtToTxt(sourceFile, targetFile string) {
	srtEntries, err := parseSRT(sourceFile)
	if err != nil {
		fmt.Println("Error parsing source SRT file:", err)
		return
	}

	file, err := os.Create(targetFile)
	if err != nil {
		fmt.Println("Error creating target file:", err)
		return
	}
	defer file.Close()

	for _, entry := range srtEntries {
		fmt.Fprintln(file, entry.Content)
	}
}
