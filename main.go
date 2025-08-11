package main

import (
	"bgp_downloader/cmd"
	"fmt"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("BGP Downloader finished successfully.")
}
