package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func parseLyrics(filePath string) (Lyrics, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Lyrics{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lyrics := Lyrics{
		Metadata: make(map[string]string),
	}
	timeLineRegex := regexp.MustCompile(`\[((\d+:)?\d+:\d+(\.\d+)?)\]`)
	tagRegex := regexp.MustCompile(`\[(\w+):(.+)\]`)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[") && strings.Contains(line, "]") {
			if timeLineRegex.MatchString(line) {
				// Timeline Tag
				time := timeLineRegex.FindString(line)
				lyrics.Timeline = append(lyrics.Timeline, time)
				// Content
				content := timeLineRegex.ReplaceAllString(line, "")
				lyrics.Content = append(lyrics.Content, strings.TrimSpace(content))
			} else {
				// Metadata
				matches := tagRegex.FindStringSubmatch(line)
				if len(matches) == 3 {
					lyrics.Metadata[matches[1]] = strings.TrimSpace(matches[2])
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return Lyrics{}, err
	}

	return lyrics, nil
}

func syncLyrics(args []string) {
	if len(args) < 2 {
		fmt.Println(SYNC_USAGE)
		return
	}

	sourceFile := args[0]
	targetFile := args[1]

	sourceLyrics, err := parseLyrics(sourceFile)
	if err != nil {
		fmt.Println("Error parsing source lyrics file:", err)
		return
	}

	targetLyrics, err := parseLyrics(targetFile)
	if err != nil {
		fmt.Println("Error parsing target lyrics file:", err)
		return
	}

	// Sync timeline
	if len(sourceLyrics.Timeline) != len(targetLyrics.Timeline) {
		fmt.Println("Warning: Timeline length mismatch")
		return
	}
	minLength := len(sourceLyrics.Timeline)
	if len(targetLyrics.Timeline) < minLength {
		minLength = len(targetLyrics.Timeline)
	}

	for i := 0; i < minLength; i++ {
		targetLyrics.Timeline[i] = sourceLyrics.Timeline[i]
	}

	// save to target, name it as "<filename>_synced.lrc"
	targetFileName := strings.TrimSuffix(targetFile, ".lrc") + "_synced.lrc"
	err = saveLyrics(targetFileName, targetLyrics)
	if err != nil {
		fmt.Println("Error saving synced lyrics file:", err)
		return
	}
}

// func printLyricsInfo(lyrics Lyrics) {
// 	fmt.Println("Metadata:")
// 	for key, value := range lyrics.Metadata {
// 		fmt.Printf("%s: %s\n", key, value)
// 	}

// 	fmt.Println("\nTimeline:")
// 	for _, time := range lyrics.Timeline {
// 		fmt.Println(time)
// 	}

// 	fmt.Println("\nLyrics Content:")
// 	for _, content := range lyrics.Content {
// 		fmt.Println(content)
// 	}
// }

func saveLyrics(filePath string, lyrics Lyrics) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write metadata
	for key, value := range lyrics.Metadata {
		fmt.Fprintf(file, "[%s: %s]\n", key, value)
	}

	// Write timeline and content
	for i := 0; i < len(lyrics.Timeline); i++ {
		fmt.Fprintf(file, "%s %s\n", lyrics.Timeline[i], lyrics.Content[i])
	}

	return nil
}
