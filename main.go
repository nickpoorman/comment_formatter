package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

	start, end := findCommentBlock(commentPrefix, lines, lineNumber)

	// Join the block of comment lines into a single string.
	block := joinBlock(commentPrefix, lines, start, end)

	intentation := determineIndentation(lines[lineNumber])
	wrappedBlock := wrapText(intentation, commentPrefix, block, lineLength)

	// Replace the lines from start to end with the wrapped block.
	// Start by removing the current lines from the slice.
	headLines := lines[:start]
	tailLines := lines[end+1:]

	var output []string
	output = append(output, headLines...)
	output = append(output, wrappedBlock...)
	output = append(output, tailLines...)

	// Overwrite the file with the output.
	err = os.WriteFile(filePath, []byte(strings.Join(output, "\n")), 0644)
	if err != nil {
		fmt.Printf("Failed to write file: %v\n", err)
		os.Exit(1)
	}
}

func indentLine(indentation int, line string) string {
	return fmt.Sprintf("%s%s", strings.Repeat(" ", indentation), line)
}

func indentLines(indentation int, lines []string) []string {
	var indented []string
	for _, line := range lines {
		indented = append(indented, indentLine(indentation, line))
	}
	return indented
}

func determineIndentation(line string) int {
	// Determine the amount of leading whitespace on the line.
	indentation := 0
	for _, r := range line {
		if r == ' ' || r == '\t' {
			indentation++
		} else {
			break
		}
	}
	return indentation
}

func isCommentBlock(commentPrefix string, lines []string, lineNumber int) bool {
	return strings.HasPrefix(strings.TrimSpace(lines[lineNumber]), commentPrefix)
}

// Finds the block of comment lines around the specified line number.
func findCommentBlock(commentPrefix string, lines []string, lineNumber int) (int, int) {
	start, end := lineNumber, lineNumber
	// Go backwards
	for start > 0 && strings.HasPrefix(strings.TrimSpace(lines[start-1]), commentPrefix) {
		start--
	}
	for end < len(lines)-1 && strings.HasPrefix(strings.TrimSpace(lines[end+1]), commentPrefix) {
		end++
	}
	return start, end
}

// Join all the lines in the block into a single string by removing the comment prefix.
func joinBlock(commentPrefix string, lines []string, start, end int) string {
	var block []string
	for i := start; i <= end; i++ {
		line := strings.TrimSpace(lines[i])
		line = strings.TrimPrefix(line, commentPrefix)
		line = strings.TrimSpace(line)
		block = append(block, line)
	}
	return strings.Join(block, " ")
}

// Wraps text to the specified length.
func wrapText(intentation int, commentPrefix string, text string, lineLength int) []string {
	words := strings.Fields(text)
	var wrapped []string
	defaultLine := indentLine(intentation, commentPrefix)
	line := defaultLine

	for _, word := range words {
		tmpLine := fmt.Sprintf("%s %s", line, word)
		if len(tmpLine) > lineLength {
			// Add the current line to the wrapped text.
			// wrapped = fmt.Sprintf("%s\n%s", wrapped, line)
			wrapped = append(wrapped, line)
			// Start a new line with the comment prefix and current word.
			line = fmt.Sprintf("%s %s", defaultLine, word)
		} else {
			// It's not too long, so add the word to the current line.
			line = tmpLine
		}
	}
	// Add the last line to the wrapped text if it's not empty.
	if len(line) > len(defaultLine) {
		// wrapped = fmt.Sprintf("%s\n%s", wrapped, line)
		wrapped = append(wrapped, line)
	}
	return wrapped
}
