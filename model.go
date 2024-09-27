package main

type Timestamp struct {
	Hours        int
	Minutes      int
	Seconds      int
	Milliseconds int
}

type Lyrics struct {
	Metadata map[string]string
	Timeline []Timestamp
	Content  []string
}

type SRTEntry struct {
	Number    int
	StartTime Timestamp
	EndTime   Timestamp
	Content   string
}

const (
	VERSION = "0.3.0"
	USAGE   = `Usage: lyc-cli [command] [options]
  Commands:
    sync	Synchronize timeline of two lyrics files
    convert	Convert lyrics file to another format
    fmt		Format lyrics file
    help	Show help`

	SYNC_USAGE    = `Usage: lyc-cli sync <source> <target>`
	CONVERT_USAGE = `Usage: lyc-cli convert <source> <target>
  Note:
  Target format is determined by file extension. Supported formats:
    .txt	Plain text format(No meta/timeline tags, only support as target format)
    .srt	SubRip Subtitle format
    .lrc	LRC format`
)
