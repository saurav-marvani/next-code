package scanners

import (
	"fmt"
	"strings"
	"time"
)

// scanLines splits content into lines for analysis
func scanLines(content string) map[int]string {
	lines := make(map[int]string)
	for i, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(line) != "" {
			lines[i+1] = line
		}
	}
	return lines
}

// generateID generates a unique ID for findings
func generateID() string {
	return fmt.Sprintf("finding_%d", time.Now().UnixNano())
}

// countLines returns the number of lines in content
func countLines(content string) int {
	return len(strings.Split(content, "\n"))
}

// getLineContent returns a specific line from content
func getLineContent(content string, lineNum int) string {
	lines := strings.Split(content, "\n")
	if lineNum > 0 && lineNum <= len(lines) {
		return lines[lineNum-1]
	}
	return ""
}

// getFileExt extracts file extension
func getFileExtension(path string) string {
	parts := strings.Split(path, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}
