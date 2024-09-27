package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func parseTimestamp(timeStr string) (Timestamp, error) {
	// 移除方括号
	timeStr = strings.Trim(timeStr, "[]")

	parts := strings.Split(timeStr, ":")
	var hours, minutes, seconds, milliseconds int
	var err error

	switch len(parts) {
	case 2: // 分钟:秒.毫秒
		minutes, err = strconv.Atoi(parts[0])
		if err != nil {
			return Timestamp{}, err
		}
		secParts := strings.Split(parts[1], ".")
		seconds, err = strconv.Atoi(secParts[0])
		if err != nil {
			return Timestamp{}, err
		}
		if len(secParts) > 1 {
			milliseconds, err = strconv.Atoi(secParts[1])
			if err != nil {
				return Timestamp{}, err
			}
			// 根据毫秒的位数调整
			switch len(secParts[1]) {
			case 1:
				milliseconds *= 100
			case 2:
				milliseconds *= 10
			}
		}
	case 3: // 小时:分钟:秒.毫秒
		hours, err = strconv.Atoi(parts[0])
		if err != nil {
			return Timestamp{}, err
		}
		minutes, err = strconv.Atoi(parts[1])
		if err != nil {
			return Timestamp{}, err
		}
		secParts := strings.Split(parts[2], ".")
		seconds, err = strconv.Atoi(secParts[0])
		if err != nil {
			return Timestamp{}, err
		}
		if len(secParts) > 1 {
			milliseconds, err = strconv.Atoi(secParts[1])
			if err != nil {
				return Timestamp{}, err
			}
			// 根据毫秒的位数调整
			switch len(secParts[1]) {
			case 1:
				milliseconds *= 100
			case 2:
				milliseconds *= 10
			}
		}
	default:
		return Timestamp{}, fmt.Errorf("无效的时间戳格式")
	}

	return Timestamp{Hours: hours, Minutes: minutes, Seconds: seconds, Milliseconds: milliseconds}, nil
}

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
				// Timeline
				timeStr := timeLineRegex.FindString(line)
				timestamp, err := parseTimestamp(timeStr)
				if err != nil {
					return Lyrics{}, err
				}
				lyrics.Timeline = append(lyrics.Timeline, timestamp)
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
		fmt.Fprintf(file, "%s %s\n", timestampToString(lyrics.Timeline[i]), lyrics.Content[i])
	}

	return nil
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

func convertLyrics(args []string) {
	if len(args) < 2 {
		fmt.Println(CONVERT_USAGE)
		return
	}

	sourceFile := args[0]
	targetFile := args[1]

	targetFmt := strings.TrimPrefix(filepath.Ext(targetFile), ".")

	switch targetFmt {
	case "txt":
		lrcToTxt(sourceFile, targetFile)
	default:
		fmt.Println("Unsupported target format:", targetFmt)
	}
}

func fmtLyrics(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: lyc-cli fmt <source>")
		return
	}

	sourceFile := args[0]

	sourceLyrics, err := parseLyrics(sourceFile)
	if err != nil {
		fmt.Println("Error parsing source lyrics file:", err)
		return
	}

	// save to target (source_name_fmt.lrc)
	targetFile := strings.TrimSuffix(sourceFile, ".lrc") + "_fmt.lrc"
	err = saveLyrics(targetFile, sourceLyrics)
	if err != nil {
		fmt.Println("Error saving formatted lyrics file:", err)
		return
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

func timestampToString(ts Timestamp) string {
	if ts.Hours > 0 {
		return fmt.Sprintf("[%02d:%02d:%02d.%03d]", ts.Hours, ts.Minutes, ts.Seconds, ts.Milliseconds)
	}
	return fmt.Sprintf("[%02d:%02d.%03d]", ts.Minutes, ts.Seconds, ts.Milliseconds)
}
