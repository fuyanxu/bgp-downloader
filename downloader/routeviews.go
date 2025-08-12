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
	routeViewsBaseURL = "https://archive.routeviews.org/"
)

var routeviewsMap = map[string]string{
	"chicago":   "route-views.chicago/bgpdata",
	"isc":       "route-views.isc/bgpdata",
	"eqix":      "route-views.eqix/bgpdata",
	"rv":        "oix-route-views",
	"rv2":       "route-views2",
	"rv3":       "route-views3",
	"rv4":       "route-views4",
	"rv6":       "route-views6",
	"linx":      "route-views.linx/bgpdata",
	"napafrica": "route-views.napafrica/bgpdata",
	"sg":        "route-views.sg/bgpdata",
	"sydney":    "route-views.sydney/bgpdata",
	"saopaulo":  "route-views2.saopaulo/bgpdata",
	"ams":       "amsix.ams/bgpdata",
}

// fileCache stores the file list for a specific monthURL to avoid duplicate requests
var routeViewsFileCache = make(map[string][]string)

func downloadDailyRVData(collector, dataType string, date time.Time, outputDir string) error {
	// Format date components
	yyyyMM := date.Format("2006.01")
	_ = date.Format("20060102") // This is required by the specification but not used in this simplified version

	// Create the base URL for the day
	dayURL := fmt.Sprintf("%s/%s/%s", routeViewsBaseURL, routeviewsMap[collector], yyyyMM)
	rib_url := fmt.Sprintf("%s/%s", dayURL, "RIBS")
	updates_url := fmt.Sprintf("%s/%s", dayURL, "UPDATES")

	rib_files, err := GetRouteViewsDailyFileList(rib_url, date)
	if err != nil {
		return fmt.Errorf("failed to get file list for %s: %v", date.Format("2006-01-02"), err)
	}

	updates_files, err := GetRouteViewsDailyFileList(updates_url, date)
	if err != nil {
		return fmt.Errorf("failed to get file list for %s: %v", date.Format("2006-01-02"), err)
	}

	// Get the list of files for the day
	// files, err := GetRouteViewsDailyFileList(dayURL, date)
	// if err != nil {
	// 	return fmt.Errorf("failed to get file list for %s: %v", date.Format("2006-01-02"), err)
	// }

	// Filter files based on data type
	var filteredFiles []string
	switch dataType {
	case "rib":
		for _, file := range rib_files {
			filteredFiles = append(filteredFiles, file)
		}
	case "updates":
		for _, file := range updates_files {
			filteredFiles = append(filteredFiles, file)
		}
	case "all":
		filteredFiles = append(rib_files, updates_files...)
	default:
		return fmt.Errorf("invalid data type: %s", dataType)
	}

	dataType = strings.ToUpper(dataType)

	// Download each file
	for _, file := range filteredFiles {

		fileURL := fmt.Sprintf("%s/%s/%s", dayURL, dataType, file)

		// Create subdirectory structure: ./collector/yyyy.mm/type
		var subDir string
		if strings.Contains(file, "rib") {
			subDir = filepath.Join(outputDir, "routeviews", "ribs", collector, yyyyMM)
		} else if strings.Contains(file, "updates") {
			subDir = filepath.Join(outputDir, "routeviews", "updates", collector, yyyyMM)
		} else {
			subDir = filepath.Join(outputDir, "routeviews", "unknown", collector, yyyyMM)
		}

		// Create the subdirectory if it doesn't exist
		if err := os.MkdirAll(subDir, 0755); err != nil {
			return fmt.Errorf("failed to create subdirectory: %v", err)
		}

		// Create the full output path
		outputPath := filepath.Join(subDir, file)
		fmt.Println(fileURL)
		if err := downloadFile(fileURL, outputPath); err != nil {
			return fmt.Errorf("failed to download %s: %v", file, err)
		}

		fmt.Printf("Downloaded: %s to %s\n", file, subDir)
	}

	return nil
}

func GetRouteViewsDailyFileList(monthURL string, date time.Time) ([]string, error) {

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
	re := regexp.MustCompile(`href="([^"]+\.bz2)`)
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
