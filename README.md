# BGP Downloader

A tool to download BGP data from RIPE and RouteViews repositories.

## Features

- Download BGP data from RIPE repository and RouteViews (planned)
- Support for RouteViews (planned)
- Customizable collector, data type, and date range
- Custom output directory
- Organizes downloaded files in subdirectories by collector, date, and type (./collector/yyyy.mm/type)## Features

- Parses HTML pages to extract file links for download
- Filters files by specific date rather than returning all files for a month
- Implements caching mechanism to avoid redundant HTTP requests for the same month
- Supports parallel downloads with a default concurrency limit of 10 goroutines for improved performance Installation

On Windows:
```bash
go build -o bgp-downloader.exe
```

On Unix-like systems (Linux, macOS):
```bash
go build -o bgp-downloader
```

To compile for Linux from any platform:
```bash
# On Unix-like systems:
GOOS=linux GOARCH=amd64 go build -o bgp-downloader-linux

# On Windows (PowerShell):
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o bgp-downloader-linux
```

## Usage

```bash
./bgp-downloader download [flags]
```

### Flags

- `-S, --source string` - Download source (ripe, routeviews) (default "ripe")
- `-c, --collector string` - Collector name (rrc00-rrc26) (default "rrc00")
- `-t, --type string` - Data type (bview/ribs, updates, all) (default "bview")
- `-s, --start-date string` - Start date (YYYY-MM-DD) (required)
- `-e, --end-date string` - End date (YYYY-MM-DD) (required)
- `-o, --output string` - Output directory (default ".")
- `-n, --concurrency int` - Maximum number of concurrent downloads (default 10)

### Collector Optional values

- ripe: rrc00-rrc26
- routeviews: "chicago", "isc", "eqix", "rv", "rv2", "rv3", "rv4", "rv6", "linx", "napafrica", "sg", "sydney", "saopaulo", "ams"

### Examples

Download bview data for a specific day:

```bash
bgp-downloader download -S ripe -c rrc00 -t bview -s 2014-03-01 -e 2014-03-01 -o ./data
```

To specify the maximum number of concurrent downloads, use the `-n` or `--concurrency` flag:

```bash
bgp-downloader download -c rrc00 -t bview -s 2014-03-01 -e 2014-03-01 -o ./data -n 20
```

Download both bview and updates data for a date range:

```bash
./bgp-downloader download -c rrc00 -t all -s 2014-03-01 -e 2014-03-03 -o ./data
```

## Testing

### Example Usage

```bash
go run test/cli_example.go
```

This will run the CLI example which demonstrates how to use the downloader with command line arguments.

To test the new directory structure feature:

```bash
go run test/directory_structure_test.go
```

This will download files and organize them in the new subdirectory structure.

## License

MIT

## Author

*Yanxu Fu* <fuyanxu@bupt.edu.com>
Created on 2025-08-02
Finish the code with help from Qwen3 Coder.