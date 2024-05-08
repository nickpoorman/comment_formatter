package formatter

import (
	"fmt"
	"strings"

	"github.com/nickpoorman/comment_formatter/stringslice"
)

func Format(commentPrefix string, lines []string, lineNumber int, lineLength int) []string {
	if lines == nil || len(lines) == 0 {
		return lines
	}

	foundComment := IsCommentBlock(commentPrefix, lines, lineNumber)

	indentation := determineIndentation(lines[lineNumber])
	start, end := findCommentBlock(commentPrefix, lines, lineNumber)

	// Break the comment block into multiple comment blocks when there is an empty comment line separating them.
	blocks := splitCommentBlock(commentPrefix, lines[start:end+1])

	// Format each of the blocks that have been seperated by comments.
	formattedBlock := formatBlocks(commentPrefix, indentation, lineLength, blocks)

	// Edge case: If we found a comment, then we should at the very least have a single line comment
	if foundComment && len(formattedBlock) == 0 {
		formattedBlock = []string{indentLine(indentation, commentPrefix)}
	}

	// Update the lines
	return replaceLines(lines, start, end, formattedBlock)
}

// Format blocks that have been seperated by empty comment lines
func formatBlocks(commentPrefix string, indentation int, lineLength int, blocks [][]string) []string {

	var wrappedSubBlocks [][]string
	// Go through each block and wrap the text to the specified line length.
	for _, subBlock := range blocks {
		// Join the block of comment lines into a single string.
		wrappedSubBlock := formatBlock(commentPrefix, indentation, lineLength, subBlock)
		wrappedSubBlocks = append(wrappedSubBlocks, wrappedSubBlock)
	}

	// Join the blocks by inserting an empty comment line between them.
	joinedSubBlocks := joinSubBlocksWithEmptyCommentLines(commentPrefix, wrappedSubBlocks)

	// Flatten all the blocks
	return stringslice.Flatten(joinedSubBlocks)
}

// Format a single comment block
func formatBlock(commentPrefix string, indentation int, lineLength int, block []string) []string {
	blockSingleLine := joinBlock(commentPrefix, block)
	return wrapText(indentation, commentPrefix, blockSingleLine, lineLength)
}

func replaceLines(lines []string, start int, end int, replaceWith []string) []string {
	// Replace the lines from start to end with the wrapped block.
	// Start by removing the current lines from the slice.
	headLines := lines[:start]
	tailLines := lines[end+1:]

	var output []string
	output = append(output, headLines...)
	output = append(output, replaceWith...)
	output = append(output, tailLines...)
	return output
}

func isEmptyLine(commentPrefix string, line string) bool {
	return strings.TrimSpace(line) == commentPrefix
}

// Insert empty line between sub blocks.
func joinSubBlocksWithEmptyCommentLines(commentPrefix string, subBlocks [][]string) [][]string {
	var out [][]string

	for i, block := range subBlocks {
		// Append the block.
		out = append(out, block)

		// Append a empty comment line block only if we are not on the last block.
		if i != len(subBlocks)-1 {
			out = append(out, []string{commentPrefix})
		}
	}

	return out
}

// Break the comment block into multiple comment blocks when there is an empty comment line separating them.
func splitCommentBlock(commentPrefix string, lines []string) [][]string {
	return stringslice.Trim(stringslice.Split(lines, commentPrefix))
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

func IsCommentBlock(commentPrefix string, lines []string, lineNumber int) bool {
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
func joinBlock(commentPrefix string, lines []string) string {
	var block []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, commentPrefix)
		line = strings.TrimSpace(line)
		block = append(block, line)
	}
	return strings.Join(block, " ")
}

// Wraps text to the specified length.
func wrapText(indentation int, commentPrefix string, text string, lineLength int) []string {
	words := strings.Fields(text)
	var wrapped []string
	defaultLine := indentLine(indentation, commentPrefix)
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
