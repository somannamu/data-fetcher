package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

var MAX_ITERATIONS = 5

func TestApplication(t *testing.T) {

	// Set up test configuration
	config := Configuration{
		Frequency:  1,
		OutputMode: OverwriteOutput,
		APIURL:     "https://api.chucknorris.io/jokes/random",
	}

	// Create a wait group to synchronize goroutines
	var wg sync.WaitGroup

	// Increment the wait group count for each goroutine
	wg.Add(3)

	// Test case 1: Overwrite output file
	go func() {
		defer wg.Done()

		// Set up test directory
		testDir := createTempDir(t)
		defer cleanupTempDir(t, testDir)

		config.OutputMode = OverwriteOutput
		config.OutputFile = filepath.Join(testDir, "output.json")
		err := runApplication(config)
		if err != nil {
			t.Errorf("Application error (Overwrite): %s", err.Error())
		}

		// Check if the output file exists
		outputFile := filepath.Join(testDir, "output.json")
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Errorf("Output file does not exist (Overwrite): %s", outputFile)
		}

		// Read the content of the output file
		content, err := ioutil.ReadFile(outputFile)
		if err != nil {
			t.Errorf("Error reading output file (Overwrite): %s", err.Error())
		}

		// Ensure the content is not empty
		if len(content) == 0 {
			t.Error("Output file content is empty (Overwrite)")
		}
	}()

	// Test case 2: Create new output file with timestamp
	go func() {
		defer wg.Done()

		// Set up test directory
		testDir := createTempDir(t)
		defer cleanupTempDir(t, testDir)

		config.OutputMode = CreateNewOutput
		config.OutputFile = filepath.Join(testDir, "output_with_timestamp.json")
		err := runApplication(config)
		if err != nil {
			t.Errorf("Application error (CreateNew): %s", err.Error())
		}

		// Check if any file exists in the output directory
		files, err := ioutil.ReadDir(testDir)
		if err != nil {
			t.Errorf("Error reading output directory (CreateNew): %s", err.Error())
		}

		// Ensure at least one file exists in the output directory
		if len(files) == 0 {
			t.Error("No file exists in the output directory (CreateNew)")
		}
	}()

	// Test case 3: Append to existing output file
	go func() {
		defer wg.Done()

		// Set up test directory
		testDir := createTempDir(t)
		defer cleanupTempDir(t, testDir)

		config.OutputMode = AppendToOutput
		appendOutputFile := filepath.Join(testDir, "output_append.json")
		config.OutputFile = appendOutputFile
		err := runApplication(config)
		if err != nil {
			t.Errorf("Application error (Append): %s", err.Error())
		}

		// Check if the existing output file exists
		if _, err := os.Stat(appendOutputFile); os.IsNotExist(err) {
			t.Errorf("Existing output file does not exist: %s", appendOutputFile)
		}

		// Read the content of the existing output file
		appendContent, err := ioutil.ReadFile(appendOutputFile)
		if err != nil {
			t.Errorf("Error reading existing output file: %s", err.Error())
		}

		// Ensure the content is not empty
		if len(appendContent) == 0 {
			t.Error("Appended output file content is empty")
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()
}

func createTempDir(t *testing.T) string {
	tempDir, err := ioutil.TempDir("", "testdir")
	if err != nil {
		t.Fatal("Failed to create temporary directory")
	}
	return tempDir
}

func cleanupTempDir(t *testing.T, dirPath string) {
	err := os.RemoveAll(dirPath)
	if err != nil {
		t.Fatalf("Failed to clean up temporary directory: %s", err.Error())
	}
}

func runApplication(config Configuration) error {
	// Create the output directories if they don't exist
	err := createOutputDirectories(config.OutputFile)
	if err != nil {
		return fmt.Errorf("Error creating output directories: %s", err.Error())
	}

	// Perform periodic fetching of data
	for i := 0; i < MAX_ITERATIONS; i++ {

		// Fetch data from the API
		err := fetchDataFromAPI(config.APIURL, config.OutputFile, config.OutputMode)
		if err != nil {
			return fmt.Errorf("Error fetching data: %s", err.Error())
		}

		// Sleep for the specified frequency
		time.Sleep(time.Duration(config.Frequency) * time.Second)

	}

	return nil
}
