package main

type Lyrics struct {
	Metadata map[string]string
	Timeline []string
	Content  []string
}

const (
	VERSION = "0.2.0"
	USAGE   = `Usage: lyc-cli [command] [options]
  Commands:
    sync	Synchronize timeline of two lyrics files
    convert	Convert lyrics file to another format
    help	Show help`

	SYNC_USAGE    = `Usage: lyc-cli sync <source> <target>`
	CONVERT_USAGE = `Usage: lyc-cli convert <source> <target>
  Note:
  Target format is determined by file extension. Supported formats:
    .txt	Plain text format(No meta/timeline tags)`
)
