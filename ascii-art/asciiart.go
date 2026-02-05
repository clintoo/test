package asciiart

import (
	"os"
	"strings"
)

const (
	asciiStart      = 32  // Space character
	asciiEnd        = 126 // Tilde character
	characterHeight = 8
	linesPerChar    = 9 // 8 art lines + 1 separator
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

	lines := strings.Split(string(data), "\n")
	
	// Remove the first empty line
	if len(lines) > 0 && lines[0] == "" {
		lines = lines[1:]
	}

	fontMap := make(map[rune][]string)

	// Each character has 8 lines of art followed by 1 separator line
	for ascii := asciiStart; ascii <= asciiEnd; ascii++ {
		index := (ascii - asciiStart) * linesPerChar
		
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
	// Normalize line endings
	text = strings.ReplaceAll(text, "\r\n", "\n")
	
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for _, line := range lines {
		if line == "" {
			// Empty line - just add a newline
			result.WriteString("\n")
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
