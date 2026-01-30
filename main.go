package main

import (
	"fmt"
	"os"
)

func main() {
	// Simple argument check
	if len(os.Args) == 1 {
		fmt.Println("8-bit Retro Style Image Converter")
		fmt.Println("Hint: Run with -h or --help for instructions.")
		return
	}

	// Help flags
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		printUsage()
		return
	}

	// Validate argument count
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: Missing arguments.")
		printUsage()
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	// 1. Open input
	fmt.Printf("Opening image: %s\n", inputPath)
	f, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// 2. Decode
	src, err := decodeImage(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode image: %v\n", err)
		os.Exit(1)
	}

	// 3. Process
	fmt.Println("Applying 8-bit retro effect (Standard Palette)...")
	processed := ProcessRetro8Bit(src)

	// 4. Save
	fmt.Printf("Saving result to: %s\n", outputPath)
	err = saveImage(outputPath, processed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save image: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Success! Conversion complete.")
}

func printUsage() {
	fmt.Printf("Usage: %s <input_image> <output_image>\n", os.Args[0])
	fmt.Println("\nArguments:")
	fmt.Println("  <input_image>   Path to the PNG or JPEG file you want to convert.")
	fmt.Println("  <output_image>  Path where the 8-bit style image will be saved.")
	fmt.Println("\nFlags:")
	fmt.Println("  -h, --help      Show this help information.")
}
