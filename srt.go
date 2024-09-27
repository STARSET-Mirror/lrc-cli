package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func parseSRT(filePath string) ([]SRTEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var entries []SRTEntry
	var currentEntry SRTEntry

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			if currentEntry.Number != 0 {
				entries = append(entries, currentEntry)
				currentEntry = SRTEntry{}
			}
			continue
		}

		if currentEntry.Number == 0 {
			currentEntry.Number, _ = strconv.Atoi(line)
		} else if currentEntry.StartTime.Hours == 0 {
			times := strings.Split(line, " --> ")
			currentEntry.StartTime = parseSRTTimestamp(times[0])
			currentEntry.EndTime = parseSRTTimestamp(times[1])
		} else {
			currentEntry.Content += line + "\n"
		}
	}

	if currentEntry.Number != 0 {
		entries = append(entries, currentEntry)
	}

	return entries, scanner.Err()
}

func convertSRT(sourceFile, targetFile, targetFmt string) {
	switch targetFmt {
	case "txt":
		srtToTxt(sourceFile, targetFile)
	case "lrc":
		srtToLrc(sourceFile, targetFile)
	default:
		fmt.Printf("unsupported target format: %s\n", targetFmt)
	}
}

func formatSRTTimestamp(ts Timestamp) string {
	return fmt.Sprintf("%02d:%02d:%02d,%03d", ts.Hours, ts.Minutes, ts.Seconds, ts.Milliseconds)
}

func parseSRTTimestamp(timeStr string) Timestamp {
	t, _ := time.Parse("15:04:05,000", timeStr)
	return Timestamp{
		Hours:        t.Hour(),
		Minutes:      t.Minute(),
		Seconds:      t.Second(),
		Milliseconds: t.Nanosecond() / 1e6,
	}
}

// basically for the last line of lrc
func addSeconds(ts Timestamp, seconds int) Timestamp {
	t := time.Date(0, 1, 1, ts.Hours, ts.Minutes, ts.Seconds, ts.Milliseconds*1e6, time.UTC)
	t = t.Add(time.Duration(seconds) * time.Second)
	return Timestamp{
		Hours:        t.Hour(),
		Minutes:      t.Minute(),
		Seconds:      t.Second(),
		Milliseconds: t.Nanosecond() / 1e6,
	}
}
