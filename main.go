package main

import (
	"time"
	"warn/pkg/stats"
)

func main() {
	for {
		stats.CheckSystemUsage()
		time.Sleep(5 * time.Second) // Check every 10 seconds
	}
}
