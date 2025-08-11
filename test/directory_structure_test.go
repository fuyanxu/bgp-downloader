package main

import (
	"fmt"
	"bgp_downloader/downloader"
)

func main() {
	// Test downloading BGP data with new directory structure
	fmt.Println("Testing BGP data download with new directory structure...")
	
	err := downloader.DownloadBGPData("rrc00", "all", "2014-03-01", "2014-03-01", "./test-output")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("Download completed successfully with new directory structure!")
}