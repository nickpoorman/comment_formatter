package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type SubBlock struct {
	lines []string
	sub   []string
	start int
	end   int
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <comment_prefix> <line_length> <line_number> <file_path>")
		os.Exit(1)
	}

	commentPrefix := os.Args[1]

	lineLength, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Invalid line length: %v\n", err)
		os.Exit(1)
	}

	lineNumber, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("Invalid line number: %v\n", err)
		os.Exit(1)
	}
	// Adjust the line number to be zero-based.
	lineNumber--

	filePath := os.Args[4]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	if lineNumber > len(lines) || lineNumber < 0 {
		fmt.Println("Line number out of bounds")
		os.Exit(1)
	}

	if !isCommentBlock(commentPrefix, lines, lineNumber) {
		fmt.Println("Line is not part of a comment block")
		os.Exit(1)
	}

	output := format(commentPrefix, lines, lineNumber, lineLength)

	writeFile(filePath, output)
}

func writeFile(filePath string, output []string) {
	// Overwrite the file with the output.
	err := os.WriteFile(filePath, []byte(strings.Join(output, "\n")), 0644)
	if err != nil {
		fmt.Printf("Failed to write file: %v\n", err)
		os.Exit(1)
	}
}
