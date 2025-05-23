package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	compiler "github.com/rabbytesoftware/quiver.compiler/compiler"
)

func main() {
	// Parse command line flags
	inputDir := flag.String("input", "", "Directory containing the Go project to compile")
	outputDir := flag.String("output", "", "Directory where the output .quiver file will be placed")
	fastMode := flag.Bool("fast", false, "Only compile for the current platform (faster)")
	
	flag.Parse()

	// Show build timestamp
	fmt.Printf("Quiver Compiler - %s - %s\n", *inputDir, time.Now().Format("2006-01-02 15:04:05"))
	
	// Validate flags
	if *inputDir == "" || *outputDir == "" {
		fmt.Println("Error: Both --input and --output flags are required")
		flag.PrintDefaults()
		os.Exit(1)
	}
	
	// Ensure input directory exists
	if _, err := os.Stat(*inputDir); os.IsNotExist(err) {
		fmt.Printf("Error: Input directory '%s' does not exist\n", *inputDir)
		os.Exit(1)
	}
	
	// Create output directory if it doesn't exist
	if _, err := os.Stat(*outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(*outputDir, 0755)
		if err != nil {
			fmt.Printf("Error creating output directory: %v\n", err)
			os.Exit(1)
		}
	}
	
	// Get the base folder name for the .watcher file
	folderName := filepath.Base(*inputDir)
	outputFile := filepath.Join(*outputDir, folderName+".quiver")
	
	// Create compiler and do the work
	compiler := compiler.NewCompiler(*inputDir, outputFile, *fastMode)
	err := compiler.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Successfully created %s\n", outputFile)
	if *fastMode {
		fmt.Println("Note: Built in fast mode - package only works on current platform")
	}
}