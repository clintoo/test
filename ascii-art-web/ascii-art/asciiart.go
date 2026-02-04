package asciiart

import (
	"bufio"
	"os"
	"strings"
)

// Public entry point
func AsciiArt(text string, banner string) string {
	lines, err := readFontFile(banner)
	if err != nil {
		return ""
	}

	templates := parseTemplates(lines)
	return printAscii(text, templates)
}

// Reads the font file and returns its lines
func readFontFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Each character is 8 lines tall.
// The font file has:
// - 1 empty line at the top
// - then blocks of 8 lines per character
func parseTemplates(lines []string) [][]string {
	var templates [][]string

	// skip the first empty line
	lines = lines[1:]

	for i := 0; i+8 <= len(lines); i += 8 {
		char := make([]string, 8)
		copy(char, lines[i:i+8])
		templates = append(templates, char)
	}

	return templates
}

func printAscii(text string, templates [][]string) string {
	// Normalize Windows newlines
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")

	var b strings.Builder

	// Process each line separately
	for i, line := range lines {
		if line != "" {
			b.WriteString(printLine(line, templates))
		}
		// Print newline between lines (but not extra at the end)
		if i < len(lines)-1 {
			b.WriteByte('\n')
		}
	}

	return b.String()
}

// Print a single line of text in ASCII art
func printLine(s string, templates [][]string) string {
	var b strings.Builder
	indexes := returnIndex(s)

	// Each character is 8 lines tall
	for row := range 8 {
		for _, index := range indexes {
			if index >= 0 && index < len(templates) {
				b.WriteString(templates[index][row])
			}
		}
		b.WriteByte('\n')
	}

	return b.String()
}

// Returns the indexes of printable ASCII characters in the font templates
func returnIndex(s string) []int {
	indexes := make([]int, 0, len(s))

	for _, r := range s {
		// Only printable ASCII characters
		if r < 32 || r > 126 {
			continue
		}
		indexes = append(indexes, int(r-32))
	}

	return indexes
}
