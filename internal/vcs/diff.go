package vcs

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// DiffParser parses unified diff format
type DiffParser struct {
	reader *bufio.Reader
}

// NewDiffParser creates a new diff parser
func NewDiffParser(r io.Reader) *DiffParser {
	return &DiffParser{
		reader: bufio.NewReader(r),
	}
}

// Parse parses a unified diff and returns file diffs
func (p *DiffParser) Parse() ([]FileDiff, error) {
	var files []FileDiff
	var currentFile *FileDiff

	for {
		line, err := p.reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}

		line = strings.TrimSuffix(line, "\n")

		// File header lines
		if strings.HasPrefix(line, "diff --git") {
			if currentFile != nil {
				files = append(files, *currentFile)
			}
			currentFile = &FileDiff{Status: FileModified}
		}

		// Extract file paths
		if strings.HasPrefix(line, "--- a/") {
			path := strings.TrimPrefix(line, "--- a/")
			if currentFile != nil {
				currentFile.OldPath = path
			}
		}
		if strings.HasPrefix(line, "+++ b/") {
			path := strings.TrimPrefix(line, "+++ b/")
			if currentFile != nil {
				currentFile.Path = path
			}
		}

		// File status detection
		if strings.HasPrefix(line, "new file mode") {
			if currentFile != nil {
				currentFile.Status = FileAdded
			}
		}
		if strings.HasPrefix(line, "deleted file mode") {
			if currentFile != nil {
				currentFile.Status = FileDeleted
			}
		}
		if strings.HasPrefix(line, "rename from") || strings.HasPrefix(line, "copy from") {
			if currentFile != nil {
				currentFile.Status = FileRenamed
			}
		}

		// Count additions and deletions
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			if currentFile != nil {
				currentFile.Additions++
				currentFile.Changes++
			}
		}
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			if currentFile != nil {
				currentFile.Deletions++
				currentFile.Changes++
			}
		}

		if err == io.EOF {
			if currentFile != nil {
				files = append(files, *currentFile)
			}
			break
		}
	}

	return files, nil
}

// ParseHunk parses a hunk header and returns the parsed hunk info
func ParseHunk(header string) (*DiffHunk, error) {
	// Format: @@ -oldStart,oldLines +newStart,newLines @@
	pattern := regexp.MustCompile(`@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@`)
	matches := pattern.FindStringSubmatch(header)

	if len(matches) < 5 {
		return nil, fmt.Errorf("invalid hunk header: %s", header)
	}

	oldStart, _ := strconv.Atoi(matches[1])
	oldLines := 1
	if matches[2] != "" {
		oldLines, _ = strconv.Atoi(matches[2])
	}

	newStart, _ := strconv.Atoi(matches[3])
	newLines := 1
	if matches[4] != "" {
		newLines, _ = strconv.Atoi(matches[4])
	}

	return &DiffHunk{
		OldStart: oldStart,
		OldLines: oldLines,
		NewStart: newStart,
		NewLines: newLines,
		Header:   header,
		Lines:    []DiffLine{},
	}, nil
}

// GetLineNumber returns the appropriate line number based on diff type
func (h *DiffHunk) GetLineNumber(diffType LineType) int {
	switch diffType {
	case LineAdded:
		return h.NewStart
	case LineRemoved:
		return h.OldStart
	case LineContext:
		return h.NewStart
	default:
		return h.NewStart
	}
}

// GetFileExtension returns the file extension from a file path
func GetFileExtension(path string) string {
	parts := strings.Split(path, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}

// IsTextFile determines if a file is likely a text file
func IsTextFile(path string) bool {
	ext := GetFileExtension(path)

	// Binary file extensions to exclude
	binaryExts := map[string]bool{
		"png": true, "jpg": true, "jpeg": true, "gif": true, "pdf": true,
		"zip": true, "tar": true, "gz": true, "bin": true, "exe": true,
		"o":   true, "a": true, "so": true, "dylib": true,
	}

	if binaryExts[ext] {
		return false
	}

	return true
}

// NormalizePath normalizes a file path for comparison
func NormalizePath(path string) string {
	return strings.TrimSpace(path)
}
