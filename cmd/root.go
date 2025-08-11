package cmd

import (
	"fmt"
	"os"

	"bgp_downloader/downloader"

	"github.com/spf13/cobra"
)

var (
	collector   string
	dataType    string
	startDate   string
	endDate     string
	outputDir   string
	concurrency int
	source      string
)

var rootCmd = &cobra.Command{
	Use:   "bgp-downloader",
	Short: "A tool to download BGP data",
	Long:  `A tool to download BGP data from RIPE and RouteViews repositories.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify a subcommand. Use --help for more information.")
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download BGP data",
	Run: func(cmd *cobra.Command, args []string) {
		err := downloader.DownloadBGPData(source, collector, dataType, startDate, endDate, outputDir, concurrency)
		if err != nil {
			fmt.Printf("Error downloading BGP data: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Download command flags
	downloadCmd.Flags().StringVarP(&source, "source", "S", "ripe", "Source (ripe, routeviews)")
	downloadCmd.Flags().StringVarP(&collector, "collector", "c", "rrc00", "Collector name (rrc00-rrc26)")
	downloadCmd.Flags().StringVarP(&dataType, "type", "t", "bview", "Data type (bview/rib, updates, all)")
	downloadCmd.Flags().StringVarP(&startDate, "start-date", "s", "", "Start date (YYYY-MM-DD) (required)")
	downloadCmd.Flags().StringVarP(&endDate, "end-date", "e", "", "End date (YYYY-MM-DD) (required)")
	downloadCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	downloadCmd.Flags().IntVarP(&concurrency, "concurrency", "n", 10, "Maximum number of concurrent downloads")

	downloadCmd.MarkFlagRequired("start-date")
	downloadCmd.MarkFlagRequired("end-date")
}
