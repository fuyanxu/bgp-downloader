package main

import (
	"fmt"
	"log"
	"time"

	downloader "bgp_downloader/downloader"
)

func main() {
	fmt.Println("Testing getDailyFileList function...")

	// Test with a known URL
	collector := "rrc00"
	date, err := time.Parse("2006-01-02", "2014-03-01")
	if err != nil {
		log.Fatalf("Failed to parse date: %v", err)
	}

	// Format date components
	yyyyMM := date.Format("2006.01")

	// Create the base URL for the day
	dayURL := fmt.Sprintf("https://data.ris.ripe.net/%s/%s", collector, yyyyMM)

	// Get the list of files for the day
	files, err := downloader.GetDailyFileList(dayURL, date)
	if err != nil {
		log.Fatalf("Failed to get file list: %v", err)
	}

	fmt.Printf("Found %d files:\n", len(files))
	for _, file := range files {
		fmt.Printf("- %s\n", file)
	}
}