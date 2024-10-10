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
	var isContent bool
	var contentBuffer strings.Builder

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			if currentEntry.Number != 0 {
				currentEntry.Content = contentBuffer.String()
				entries = append(entries, currentEntry)
				currentEntry = SRTEntry{}
				isContent = false
				contentBuffer.Reset()
			}
			continue
		}

		if currentEntry.Number == 0 {
			currentEntry.Number, _ = strconv.Atoi(line)
		} else if isEntryTimeStampUnset(currentEntry) {
			times := strings.Split(line, " --> ")
			if len(times) == 2 {
				currentEntry.StartTime = parseSRTTimestamp(times[0])
				currentEntry.EndTime = parseSRTTimestamp(times[1])
				isContent = true
			}
		} else if isContent {
			if contentBuffer.Len() > 0 {
				contentBuffer.WriteString("\n")
			}
			contentBuffer.WriteString(line)
		}
	}

	if currentEntry.Number != 0 {
		currentEntry.Content = contentBuffer.String()
		entries = append(entries, currentEntry)
	}

	return entries, scanner.Err()
}

func isEntryTimeStampUnset(currentEntry SRTEntry) bool {
	return currentEntry.StartTime.Hours == 0 && currentEntry.StartTime.Minutes == 0 &&
		currentEntry.StartTime.Seconds == 0 && currentEntry.StartTime.Milliseconds == 0 &&
		currentEntry.EndTime.Hours == 0 && currentEntry.EndTime.Minutes == 0 &&
		currentEntry.EndTime.Seconds == 0 && currentEntry.EndTime.Milliseconds == 0
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
	parts := strings.Split(timeStr, ",")
	if len(parts) != 2 {
		return Timestamp{}
	}

	timeParts := strings.Split(parts[0], ":")
	if len(timeParts) != 3 {
		return Timestamp{}
	}

	hours, _ := strconv.Atoi(timeParts[0])
	minutes, _ := strconv.Atoi(timeParts[1])
	seconds, _ := strconv.Atoi(timeParts[2])
	milliseconds, _ := strconv.Atoi(parts[1])

	return Timestamp{
		Hours:        hours,
		Minutes:      minutes,
		Seconds:      seconds,
		Milliseconds: milliseconds,
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
