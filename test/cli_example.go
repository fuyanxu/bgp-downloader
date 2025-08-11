package main

import (
	"fmt"
	"os"

	"bgp_downloader/cmd"
)

func main() {
	// Test the CLI interface
	fmt.Println("Testing CLI interface...")
	
	// Simulate command line arguments
	os.Args = []string{"bgp-downloader", "download", "-c", "rrc00", "-t", "bview", "-s", "2014-03-01", "-e", "2014-03-01", "-o", "./test-data"}
	
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	
	// Test with concurrency parameter
	fmt.Println("Testing CLI interface with concurrency parameter...")
	
	// Simulate command line arguments with concurrency
	os.Args = []string{"bgp-downloader", "download", "-c", "rrc00", "-t", "bview", "-s", "2014-03-01", "-e", "2014-03-01", "-o", "./test-data", "-n", "5"}
	
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("CLI test completed successfully!")
}