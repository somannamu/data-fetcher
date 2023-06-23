package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Configuration represents the application configuration
type Configuration struct {
	Frequency  int // Frequency must be specified in seconds
	OutputFile string
	OutputMode string
	APIURL     string
}

// OutputMode enumeration
const (
	OverwriteOutput = "overwrite"
	CreateNewOutput = "create" // Based on Timestamp
	AppendToOutput  = "append"
)

func parseArguments() Configuration {
	// Define default values for configuration
	defaultFrequency := 60 // 60 seconds
	defaultOutputFile := "output/output.json"
	defaultOutputMode := OverwriteOutput
	defaultAPIURL := "https://api.chucknorris.io/jokes/random"

	// Parse command-line arguments
	frequency := flag.Int("frequency", defaultFrequency, "Frequency of data fetching in seconds")
	outputFile := flag.String("output", defaultOutputFile, "Output file path")
	outputMode := flag.String("output-mode", defaultOutputMode, "Output file mode (overwrite = Overwrite, create = Create new file with timestamp, append = Append to existing file)")
	apiURL := flag.String("api-url", defaultAPIURL, "API URL")
	flag.Parse()

	// Return the parsed configuration
	return Configuration{
		Frequency:  *frequency,
		OutputFile: *outputFile,
		OutputMode: *outputMode,
		APIURL:     *apiURL,
	}
}

func createOutputDirectories(outputFile string) error {
	dir := filepath.Dir(outputFile)

	// Check if the output path includes directories
	if dir != "." && dir != string(filepath.Separator) {
		// Create the output directories if they don't exist
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func fetchDataFromAPI(apiURL string, outputFile string, outputMode string) error {
	// Make HTTP request to the API
	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Discard the response body
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Save a placeholder message to the output file
	err = saveDataToFile([]byte("API response ignored"), outputFile, outputMode)
	if err != nil {
		return err
	}

	fmt.Println("Data saved to file:", outputFile)

	return nil
}

func saveDataToFile(data []byte, outputFile string, outputMode string) error {
	file, err := openOutputFile(outputFile, outputMode)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write data to file
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func openOutputFile(outputFile string, outputMode string) (*os.File, error) {
	if outputMode == CreateNewOutput {
		// Append timestamp to the output file
		dir := filepath.Dir(outputFile)
		baseName := filepath.Base(outputFile)
		ext := filepath.Ext(outputFile)
		fileName := baseName[:len(baseName)-len(ext)]
		timeStamp := time.Now().Format("20060102150405") // YYYYMMDDHHMMSS
		outputFile = filepath.Join(dir, fmt.Sprintf("%s_%s%s", fileName, timeStamp, ext))
	}

	if outputMode == AppendToOutput {
		// Open file in append mode
		return os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	}

	// Open file in write mode, creating a new file if it doesn't exist
	return os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
}

func main() {
	// Parse command-line arguments
	config := parseArguments()

	// Create the output directories if they don't exist
	err := createOutputDirectories(config.OutputFile)
	if err != nil {
		fmt.Printf("Error creating output directories: %s\n", err.Error())
		return
	}

	// Perform periodic fetching of data
	for {
		// Fetch data from the API
		err := fetchDataFromAPI(config.APIURL, config.OutputFile, config.OutputMode)
		if err != nil {
			fmt.Printf("Error fetching data: %s\n", err.Error())
		}

		// Sleep for the specified frequency
		time.Sleep(time.Duration(config.Frequency) * time.Second)
	}
}
