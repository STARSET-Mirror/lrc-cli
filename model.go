package main

type Lyrics struct {
	Metadata map[string]string
	Timeline []string
	Content  []string
}

const USAGE = `Usage: lyc-cli [command] [options]
Commands:
  sync	Synchronize timeline of two lyrics files
  help	Show help`

const SYNC_USAGE = `Usage: lyc-cli sync <source> <target>`
