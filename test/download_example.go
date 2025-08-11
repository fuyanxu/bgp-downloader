package main

import (
	"bgp_downloader/downloader"
	"fmt"
)

func main() {
	// Test downloading BGP data
	fmt.Println("Testing BGP data download...")

	err := downloader.DownloadBGPData("rrc00", "bview", "2014-03-01", "2014-03-01", "./test-data")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Download completed successfully!")
}