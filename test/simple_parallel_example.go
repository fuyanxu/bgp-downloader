package main

import (
	"fmt"
	"time"
	"bgp_downloader/downloader"
)

func main() {
	// Test the parallel download functionality with multiple days
	fmt.Println("Testing parallel download functionality with multiple days...")
	
	// Record start time
	start := time.Now()
	
	// Download data for multiple days with concurrency limit
	err := downloader.DownloadBGPData("rrc00", "bview", "2014-03-01", "2014-03-05", "./test-data-multi", 5)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Record end time
	elapsed := time.Since(start)
	fmt.Printf("Parallel download of 5 days completed in %s\n", elapsed)
	
	// Test with higher concurrency
	fmt.Println("\nTesting with higher concurrency...")
	start = time.Now()
	err = downloader.DownloadBGPData("rrc00", "bview", "2014-03-01", "2014-03-05", "./test-data-multi-high", 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Record end time
	elapsed = time.Since(start)
	fmt.Printf("High concurrency download of 5 days completed in %s\n", elapsed)
	
	fmt.Println("Multi-day parallel download test completed successfully!")
}