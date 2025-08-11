package main

import (
	"fmt"
	"time"

	"bgp_downloader/downloader"
)

func main() {
	// Test with different concurrency levels
	testConcurrency(5)
	testConcurrency(10)
}

func testConcurrency(concurrency int) {
	fmt.Printf("\n=== Testing with concurrency level: %d ===\n", concurrency)
	
	start := time.Now()
	err := downloader.DownloadBGPData("rrc06", "bview", "2014-03-01", "2014-03-03", "test-data-parallel", concurrency)
	if err != nil {
		fmt.Printf("Download completed with errors: %v\n", err)
	}
	
	elapsed := time.Since(start)
	fmt.Printf("Download completed with concurrency %d in %s\n", concurrency, elapsed)
}