package downloader

import (
	"fmt"
	"os"
	"time"
)

// DownloadBGPData is the main function to download BGP data
// It determines the source and calls the appropriate downloader
func DownloadBGPData(source, collector, dataType, startDate, endDate, outputDir string, maxConcurrency int) error {
	// For now, we only support RIPE
	// In the future, we can add support for RouteViews here
	if source == "ripe" {
		return downloadRipeData(collector, dataType, startDate, endDate, outputDir, maxConcurrency)
	}
	if source == "routeviews" {
		return downloadRouteViewsData(collector, dataType, startDate, endDate, outputDir, maxConcurrency)
	}
	return fmt.Errorf("invalid source: %s", source)
}

// downloadRipeData is the internal implementation for downloading RIPE data
func downloadRipeData(collector, dataType, startDate, endDate, outputDir string, maxConcurrency int) error {
	// Validate collector
	if !isValidCollector("ripe", collector) {
		return fmt.Errorf("invalid collector: %s", collector)
	}

	// Parse dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return fmt.Errorf("invalid start date: %v", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return fmt.Errorf("invalid end date: %v", err)
	}

	// Validate date range
	if start.After(end) {
		return fmt.Errorf("start date cannot be after end date")
	}

	// Create output directory if it doesn't exist
	if err := createOutputDir(outputDir); err != nil {
		return err
	}

	// Use default concurrency if not specified or invalid
	if maxConcurrency <= 0 {
		maxConcurrency = 10
	}

	// Create a channel to control concurrency
	semaphore := make(chan struct{}, maxConcurrency)

	// Create a channel to collect errors from goroutines
	errChan := make(chan error, 1)

	// Use a done channel to signal when all goroutines are finished
	done := make(chan struct{})

	// Counter for active goroutines
	var activeGoroutines int

	// Download data for each day in the range
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		activeGoroutines++
		go func(date time.Time) {
			// Acquire semaphore
			semaphore <- struct{}{}

			// Release semaphore when done
			defer func() { <-semaphore }()

			// Perform the download
			if err := downloadDailyRipeData(collector, dataType, date, outputDir); err != nil {
				// Send error to channel, but only if no error has been sent yet
				select {
				case errChan <- err:
				default:
				}
			}

			// Decrement active goroutines counter
			activeGoroutines--

			// If no more active goroutines, signal done
			if activeGoroutines == 0 {
				done <- struct{}{}
			}
		}(d)
	}

	// Wait for either all downloads to complete or an error to occur
	select {
	case <-done:
		return nil
	case err := <-errChan:
		return err
	}
}

func downloadRouteViewsData(collector, dataType, startDate, endDate, outputDir string, maxConcurrency int) error {
	// Validate collector
	if !isValidCollector("routeviews", collector) {
		return fmt.Errorf("invalid collector: %s", collector)
	}

	// Parse dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return fmt.Errorf("invalid start date: %v", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return fmt.Errorf("invalid end date: %v", err)
	}

	// Validate date range
	if start.After(end) {
		return fmt.Errorf("start date cannot be after end date")
	}

	// Create output directory if it doesn't exist
	if err := createOutputDir(outputDir); err != nil {
		return err
	}

	// Use default concurrency if not specified or invalid
	if maxConcurrency <= 0 {
		maxConcurrency = 10
	}

	// Create a channel to control concurrency
	semaphore := make(chan struct{}, maxConcurrency)

	// Create a channel to collect errors from goroutines
	errChan := make(chan error, 1)

	// Use a done channel to signal when all goroutines are finished
	done := make(chan struct{})

	// Counter for active goroutines
	var activeGoroutines int

	// Download data for each day in the range
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		activeGoroutines++
		go func(date time.Time) {
			// Acquire semaphore
			semaphore <- struct{}{}

			// Release semaphore when done
			defer func() { <-semaphore }()

			// Perform the download
			if err := downloadDailyRouteViewsData(collector, dataType, date, outputDir); err != nil {
				// Send error to channel, but only if no error has been sent yet
				select {
				case errChan <- err:
				default:
				}
			}

			// Decrement active goroutines counter
			activeGoroutines--

			// If no more active goroutines, signal done
			if activeGoroutines == 0 {
				done <- struct{}{}
			}
		}(d)
	}

	// Wait for either all downloads to complete or an error to occur
	select {
	case <-done:
		return nil
	case err := <-errChan:
		return err
	}
}

// isValidCollector checks if the collector is valid for RIPE
func isValidCollector(source string, collector string) bool {
	switch source {
	case "ripe":
		validCollectors := map[string]bool{
			"rrc00": true, "rrc01": true, "rrc02": true, "rrc03": true, "rrc04": true,
			"rrc05": true, "rrc06": true, "rrc07": true, "rrc08": true, "rrc09": true,
			"rrc10": true, "rrc11": true, "rrc12": true, "rrc13": true, "rrc14": true,
			"rrc15": true, "rrc16": true, "rrc17": true, "rrc18": true, "rrc19": true,
			"rrc20": true, "rrc21": true, "rrc22": true, "rrc23": true, "rrc24": true,
			"rrc25": true, "rrc26": true,
		}
		return validCollectors[collector]
	case "routeviews":
		validCollectors := map[string]bool{
			"chicago":   true,
			"isc":       true,
			"eqix":      true,
			"rv":        true,
			"rv2":       true,
			"rv3":       true,
			"rv4":       true,
			"rv6":       true,
			"linx":      true,
			"napafrica": true,
			"sg":        true,
			"sydney":    true,
			"saopaulo":  true,
			"ams":       true,
		}
		return validCollectors[collector]
	}
	return false
}

// createOutputDir creates the output directory if it doesn't exist
func createOutputDir(outputDir string) error {
	return os.MkdirAll(outputDir, 0755)
}

// downloadDailyRipeData downloads data for a specific day from RIPE
func downloadDailyRipeData(collector, dataType string, date time.Time, outputDir string) error {
	return downloadDailyData(collector, dataType, date, outputDir)
}

func downloadDailyRouteViewsData(collector, dataType string, date time.Time, outputDir string) error {
	return downloadDailyRVData(collector, dataType, date, outputDir)
}
