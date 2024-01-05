package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose output")
	help     = flag.Bool("help", false, "help")
	dataDir  = flag.String("dir", "", "directory to read files from")
	dataFile = flag.String("file", "", "file to read from")
	jsonFile = flag.String("json", "", "json file to read from")
)

func main() {
	flag.Parse()

	if *help {
		fmt.Println("Usage: go run main.go (-d <directory>|-f <file>) -j <json file>")
		return
	}

	jsonStr, err := loadJSON()
	if err != nil {
		fmt.Printf("Failed to load JSON: %s\n", err)
		return
	}

	err = json.Unmarshal([]byte(jsonStr), new((map[string]interface{})))
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON skeleton, please check your JSON is valid and try again: %s\n", err)
		return
	}

	files, err := loadData()
	if err != nil {
		fmt.Printf("Failed to load data: %s\n", err)
		return
	}

	startTime := time.Now()
	defer func() {
		if *verbose {
			fmt.Printf("Total time taken: %s\n", time.Since(startTime))
		}
	}()

	ch := make(chan string)
	for _, file := range files {
		if *verbose {
			fmt.Printf("Filling file %s\n", file)
		}

		go fill(file, jsonStr, ch)

		// Don't blast the API too hard
		time.Sleep(1 * time.Second)
	}

	for range files {
		if *verbose {
			fmt.Printf("Waiting for file %s to finish\n", <-ch)
		}
		fmt.Printf(<-ch)
	}
}

func fill(file File, jsonStr string, ch chan string) {
	startTime := time.Now()

	if *verbose {
		fmt.Println("Reading data from file")
	}

	if *verbose {
		fmt.Println("Requesting filled data from LM")
	}

	data, err := file.Load()
	if err != nil {
		ch <- fmt.Sprintf("\nFailed to load data from file %s: %s\n", file, err)
		return
	}

	result, err := requestFill(jsonStr, data)
	if err != nil {
		ch <- fmt.Sprintf("\nFailed to request fill for file %s: %s\n", file, err)
		return
	}

	ch <- result
	if *verbose {
		fmt.Printf("Time taken for file %s: %s\n", file, time.Since(startTime))
	}
}

func loadJSON() (string, error) {
	var jsonStr string

	// Get the json structure if the user provided a file
	if _, err := os.Stat(*jsonFile); err == nil {
		if *verbose {
			fmt.Printf("Reading JSON from file %s\n", *jsonFile)
		}
		jsonFile, err := os.Open(*jsonFile)
		if err != nil {
			return "", err
		}
		defer jsonFile.Close()

		bytes, err := io.ReadAll(jsonFile)

		if err != nil {
			return "", err
		}

		jsonStr = string(bytes)
		if *verbose {
			fmt.Printf("JSON: %s\n", jsonStr)
		}
	}

	if jsonStr == "" {
		return "", fmt.Errorf("please provide an empty JSON structure that you would like to be populated from the data you provide")
	}

	return jsonStr, nil
}

func loadData() ([]File, error) {
	var files []File

	// If the user provided a directory, read all files from that directory
	if *dataDir != "" {
		if *verbose {
			fmt.Printf("Reading data from directory %s\n", *dataDir)
		}
		dirFiles, err := os.ReadDir(*dataDir)
		if err != nil {
			return nil, err
		}

		for _, file := range dirFiles {
			newFile, err := NewFile(filepath.Join(*dataDir, file.Name()), filepath.Ext(file.Name()))
			if err != nil {
				return nil, err
			}

			files = append(files, newFile)
		}
	} else if *dataFile != "" {
		file, err := NewFile(*dataFile, filepath.Ext(*dataFile))
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	} else {
		return nil, fmt.Errorf("please provide a directory or file to read data from")
	}

	if *verbose {
		for _, file := range files {
			fmt.Printf("Data file found: %s\n", file.Path())
		}
	}

	return files, nil
}
