package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/unidoc/unipdf/v3/common/license"
)

var (
	verbose          = flag.Bool("verbose", false, "verbose output")
	help             = flag.Bool("help", false, "help")
	dataDir          = flag.String("dir", "", "directory to read files from")
	dataFile         = flag.String("file", "", "file to read from")
	jsonFile         = flag.String("json", "", "json file to read from")
	instructionsFile = flag.String("instructions", "", "instructions file to read from")
)

func init() {
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}

	flag.Parse()
}

func main() {
	if *help {
		fmt.Println("This program will fill a JSON structure with data from a file or directory of files.")
		fmt.Println("The JSON structure must be provided as a file, and the data must be provided as a file or directory of files.")
		fmt.Println("Extra instructions must be provided as a text file.")
		fmt.Println("Usage: go run main.go (-dir <directory>|-file <filename>) -json <json file> -instructions <instructions file>")
		fmt.Println("Example: go run main.go -dir ./data -json ./json.json -instructions ./instructions.txt")
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

	// Create a slice of [File] to hold all the files we want to read from.
	files, err := loadFiles()
	if err != nil {
		fmt.Printf("Failed to load data: %s\n", err)
		return
	}

	// Get any extra instructions from the user to be given to the LM.
	instructions, err := loadInstructions()
	if err != nil {
		fmt.Printf("Failed to load instructions: %s\n", err)
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

		go fill(file, jsonStr, instructions, ch)

		// Don't blast the API too hard
		time.Sleep(1 * time.Second)
	}

	results := []string{}
	for range files {
		if *verbose {
			fmt.Printf("Waiting for file %s to finish\n", <-ch)
		}

		results = append(results, <-ch)
	}

	// Print a json array of the items in results
	fmt.Printf("[%s", results[0])
	for _, result := range results[1:] {
		fmt.Printf(",%s", result)
	}
	fmt.Printf("]")
}

func fill(file File, jsonStr string, instructions string, ch chan string) {
	startTime := time.Now()

	if *verbose {
		fmt.Println("Reading data from file")
	}

	data, err := file.Load()
	if err != nil {
		ch <- fmt.Sprintf("\nFailed to load data from file %s: %s\n", file, err)
		return
	}

	if *verbose {
		fmt.Println("Requesting filled data from LM")
	}

	result, err := requestFill(jsonStr, data, instructions)
	if err != nil {
		ch <- fmt.Sprintf("\nFailed to request fill for file %s: %s\n", file, err)
		return
	}

	// Unmarshal the result and add a new filename key to it
	var resultJSON map[string]interface{}
	err = json.Unmarshal([]byte(result), &resultJSON)
	if err != nil {
		ch <- fmt.Sprintf("\nFailed to unmarshal result for file %s: %s\n", file, err)
		return

	}

	resultJSON["filename"] = file.Name()
	bytes, err := json.MarshalIndent(resultJSON, "", "\t")
	if err != nil {
		ch <- fmt.Sprintf("\nFailed to marshal result for file %s: %s\n", file, err)
		return
	}

	ch <- string(bytes)
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

func loadFiles() ([]File, error) {
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

func loadInstructions() (string, error) {
	instructionsFile, err := os.Open(*instructionsFile)
	if err != nil {
		return "", err
	}
	defer instructionsFile.Close()

	bytes, err := io.ReadAll(instructionsFile)
	if err != nil {
		return "", err
	}

	instructions := string(bytes)
	if *verbose {
		fmt.Printf("Instructions loaded: %s\n", instructions)
	}

	return instructions, nil
}
