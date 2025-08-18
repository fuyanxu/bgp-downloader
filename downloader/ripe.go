package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	baseURL = "https://data.ris.ripe.net"
)

// fileCache stores the file list for a specific monthURL to avoid duplicate requests
var fileCache = make(map[string][]string)

func downloadDailyData(collector, dataType string, date time.Time, outputDir string) error {
	// Format date components
	yyyyMM := date.Format("2006.01")
	_ = date.Format("20060102") // This is required by the specification but not used in this simplified version

	// Create the base URL for the day
	monthURL := fmt.Sprintf("%s/%s/%s", baseURL, collector, yyyyMM)

	// Get the list of files for the day
	files, err := GetMonthlyFileList(monthURL, date)
	if err != nil {
		return fmt.Errorf("failed to get file list for %s: %v", date.Format("2006-01-02"), err)
	}

	// Filter files based on data type
	var filteredFiles []string
	switch dataType {
	case "bview":
		for _, file := range files {
			if strings.Contains(file, "bview") {
				filteredFiles = append(filteredFiles, file)
			}
		}
	case "updates":
		for _, file := range files {
			if strings.Contains(file, "updates") {
				filteredFiles = append(filteredFiles, file)
			}
		}
	case "all":
		filteredFiles = files
	default:
		return fmt.Errorf("invalid data type: %s", dataType)
	}

	// Download each file
	for _, file := range filteredFiles {
		fileURL := fmt.Sprintf("%s/%s", monthURL, file)

		// Create subdirectory structure: ./collector/yyyy.mm/type
		var subDir string
		if strings.Contains(file, "bview") {
			subDir = filepath.Join(outputDir, "ripe", "bview", collector, yyyyMM)
		} else if strings.Contains(file, "updates") {
			subDir = filepath.Join(outputDir, "ripe", "updates", collector, yyyyMM)
		} else {
			subDir = filepath.Join(outputDir, "ripe", "unknown", collector, yyyyMM)
		}

		// Create the subdirectory if it doesn't exist
		if err := os.MkdirAll(subDir, 0755); err != nil {
			return fmt.Errorf("failed to create subdirectory: %v", err)
		}

		// Create the full output path
		outputPath := filepath.Join(subDir, file)

		if err := downloadFile(fileURL, outputPath); err != nil {
			return fmt.Errorf("failed to download %s: %v", file, err)
		}

		fmt.Printf("Downloaded: %s to %s\n", file, subDir)
	}

	return nil
}

func GetMonthlyFileList(monthURL string, date time.Time) ([]string, error) {
	// Extract monthURL from dayURL
	// urlParts := strings.Split(dayURL, "/")
	// if len(urlParts) < 5 {
	// 	return nil, fmt.Errorf("invalid URL format: %s", dayURL)
	// }
	// monthURL := strings.Join(urlParts[:len(urlParts)-1], "/")

	// Check if we have cached results for this monthURL
	if cachedFiles, exists := fileCache[monthURL]; exists {
		// Filter cached files based on the specified date
		dateStr := date.Format("20060102") // Format date as YYYYMMDD
		var filteredFiles []string
		for _, file := range cachedFiles {
			// Check if the file name contains the date string
			if strings.Contains(file, dateStr) {
				filteredFiles = append(filteredFiles, file)
			}
		}
		return filteredFiles, nil
	}

	// Make an HTTP GET request to dayURL
	resp, err := http.Get(monthURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL %s: %v", monthURL, err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status for %s: %s", monthURL, resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %v", monthURL, err)
	}

	// Parse the HTML response to extract .gz file links
	// Using a simple regex to find href attributes ending with .gz
	re := regexp.MustCompile(`href="([^"]+\.gz)`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	// Extract file names from the matches
	var files []string
	for _, match := range matches {
		if len(match) > 1 {
			files = append(files, match[1])
		}
	}

	// Cache the file list for this monthURL
	fileCache[monthURL] = files

	// Filter files based on the specified date
	dateStr := date.Format("20060102") // Format date as YYYYMMDD
	var filteredFiles []string
	for _, file := range files {
		// Check if the file name contains the date string
		if strings.Contains(file, dateStr) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles, nil
}

func downloadFile(url, outputPath string) error {
	// Check if file already exists
	if _, err := os.Stat(outputPath); err == nil {
		return nil // File already exists, skip download
	}

	// Create the file
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	//Retry parameters
	const maxRetries = 5
	retryDelay := 1 * time.Second

	// Retry loop
	for i:=0; i <maxRetries;i++ {
		resp, err := http.Get(url)
		if err != nil {
			if i == maxRetries {
				return err // Return the error if we've exhausted retries
			}
			fmt.Printf("Download attempt %d failed: %v. Retrying in %v...\n", i+1, err, retryDelay)
			time.Sleep(retryDelay)
			retryDelay *= 2 // Exponential backoff
			continue
		}
		defer resp.Body.Close()

		// Check server response
		if resp.StatusCode != http.StatusOK {
			if i == maxRetries {
				return fmt.Errorf("bad status: %s", resp.Status)
			}
			fmt.Printf("Download attempt %d failed with status %s. Retrying in %v...\n", i+1, resp.Status, retryDelay)
			time.Sleep(retryDelay)
			retryDelay *= 2 // Exponential backoff
			continue
		}

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			if i == maxRetries {
				return err
			}
			fmt.Printf("Failed to write to file on attempt %d: %v. Retrying in %v...\n", i+1, err, retryDelay)
			time.Sleep(retryDelay)
			retryDelay *= 2 // Exponential backoff
			continue
		}

		// Success
		fmt.Printf("Successfully downloaded %s\n", outputPath)
		return nil
	}

	return nil
}
