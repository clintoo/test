package asciiart

import (
	"os"
	"strings"
)

const (
	asciiStart      = 32  // Space character (first printable ASCII)
	asciiEnd        = 126 // Tilde character (last printable ASCII)
	characterHeight = 8   // Each ASCII art character is 8 lines tall
	linesPerChar    = 9   // 8 art lines + 1 separator line per character
)

// AsciiArt converts text to ASCII art using the specified banner file
func AsciiArt(text string, bannerPath string) string {
	// Load the banner font
	fontMap, err := loadBanner(bannerPath)
	if err != nil {
		return ""
	}

	// Render the text
	return render(text, fontMap)
}

// loadBanner reads a banner file and returns a map of characters to their ASCII art
func loadBanner(path string) (map[rune][]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Normalize line endings (handle both \r\n and \n)
	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	lines := strings.Split(content, "\n")
	
	// Remove the first empty line if present
	if len(lines) > 0 && lines[0] == "" {
		lines = lines[1:]
	}

	fontMap := make(map[rune][]string)

	// Each character has 8 lines of art followed by 1 separator line
	for ascii := asciiStart; ascii <= asciiEnd; ascii++ {
		index := (ascii - asciiStart) * linesPerChar
		
		// Ensure we have enough lines for this character
		if index+characterHeight > len(lines) {
			break
		}

		// Extract the 8 lines for this character
		charLines := make([]string, characterHeight)
		for i := 0; i < characterHeight; i++ {
			charLines[i] = lines[index+i]
		}
		
		fontMap[rune(ascii)] = charLines
	}

	return fontMap, nil
}

// render converts the input text to ASCII art
func render(text string, fontMap map[rune][]string) string {
	if text == "" {
		return ""
	}

	// Normalize line endings (handle both \r\n and \n)
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for _, line := range lines {
		if line == "" {
			// Empty line - output 8 empty rows to match character height
			for row := 0; row < characterHeight; row++ {
				result.WriteString("\n")
			}
		} else {
			// Render each of the 8 rows for this line of text
			for row := 0; row < characterHeight; row++ {
				for _, char := range line {
					if artLines, exists := fontMap[char]; exists {
						result.WriteString(artLines[row])
					}
				}
				result.WriteString("\n")
			}
		}
	}

	// Remove trailing newline
	output := result.String()
	return strings.TrimSuffix(output, "\n")
}
