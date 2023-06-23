# Data Fetcher

Data Fetcher is a Go application that periodically fetches data from a publicly available API and saves the results as a JSON file on the hard disk.

## Features

- Periodically fetches data from a specified API endpoint
- Saves the fetched data as a JSON file on the hard disk
- Supports configuration of the fetch frequency and output file location
- Provides options to control how the output data is stored (overwrite, create new file with timestamp, or append to existing file)
- Uses only the Go standard library without any third-party dependencies

## Getting Started

### Prerequisites

- Go programming language
- Internet connection to access the API

### Installation

1. Clone the repository:

```bash
git clone https://github.com/somannamu/data-fetcher.git
```

2. Change into the project directory:

```bash
cd data-fetcher
```

### Usage

To run the application with default configuration:

```bash
go run main.go
```

To customize the configuration, you can provide command line arguments:

```bash
go run main.go -api-url https://api.chucknorris.io/jokes/random -frequency 30 -output path/to/output.json -output-mode append
```

#### Configuration Options

- `api-url` (optional): Specifies the API URL from where the data will be fetched.
- `frequency` (optional): Specifies the fetch frequency in seconds. Default is 60 seconds.
- `output` (optional): Specifies the output file path. Default is "output/output.json".
- `output-mode` (optional): Specifies the output mode. Options are "overwrite", "create", and "append". Default is "overwrite".

### Testing

To run the tests, execute the following command:

```bash
go test -v
```

## Contributing

Contributions are welcome! If you encounter any issues or have suggestions for improvement, please create a new issue or submit a pull request.
