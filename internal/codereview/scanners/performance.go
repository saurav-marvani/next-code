package scanners

import (
	"context"
	"regexp"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// PerformanceScanner performs performance-focused code analysis
type PerformanceScanner struct {
	enabled bool
	rules   []*performanceRule
}

type performanceRule struct {
	name     string
	pattern  *regexp.Regexp
	severity codereview.Severity
	message  string
}

// NewPerformanceScanner creates a new performance scanner
func NewPerformanceScanner() *PerformanceScanner {
	scanner := &PerformanceScanner{
		enabled: true,
		rules:   make([]*performanceRule, 0),
	}
	scanner.initializeRules()
	return scanner
}

func (s *PerformanceScanner) initializeRules() {
	rules := []struct {
		name     string
		pattern  string
		severity codereview.Severity
		message  string
	}{
		{
			name:     "N+1 Query Problem",
			pattern:  `(?i)for\s+.*\{[^}]*(?:query|find|select)`,
			severity: codereview.SeverityHigh,
			message:  "Potential N+1 query problem: consider using batch queries or eager loading",
		},
		{
			name:     "Unnecessary Loop",
			pattern:  `(?i)for\s+.*\{\s*$`,
			severity: codereview.SeverityMedium,
			message:  "Verify this loop is necessary; consider using map/filter/reduce",
		},
		{
			name:     "String Concatenation in Loop",
			pattern:  `(?i)for\s+.*\{[^}]*\s*\+\s*=\s*["\']`,
			severity: codereview.SeverityMedium,
			message:  "String concatenation in loop is inefficient: use StringBuilder or buffer",
		},
		{
			name:     "Synchronous I/O",
			pattern:  `(?i)(readFile|readLine|fetch)\s*\(.*\)(?!\s*await|\s*then|\s*\))`,
			severity: codereview.SeverityMedium,
			message:  "Consider using async I/O to avoid blocking operations",
		},
		{
			name:     "Memory Leak Potential",
			pattern:  `(?i)(?:setInterval|setTimeout|addEventListener).*function`,
			severity: codereview.SeverityMedium,
			message:  "Ensure timers and event listeners are properly cleaned up",
		},
		{
			name:     "Large Object Creation",
			pattern:  `(?i)new\s+(?:Array|Object)\s*\(\s*\d{4,}`,
			severity: codereview.SeverityLow,
			message:  "Creating large objects in memory: verify this is necessary",
		},
	}

	for _, r := range rules {
		if pattern, err := regexp.Compile(r.pattern); err == nil {
			s.rules = append(s.rules, &performanceRule{
				name:     r.name,
				pattern:  pattern,
				severity: r.severity,
				message:  r.message,
			})
		}
	}
}

// GetType returns the scanner type
func (s *PerformanceScanner) GetType() codereview.ScannerType {
	return codereview.ScannerPerformance
}

// Scan performs performance analysis on the code
func (s *PerformanceScanner) Scan(ctx context.Context, req *codereview.ReviewRequest) ([]codereview.Finding, error) {
	findings := make([]codereview.Finding, 0)

	for _, file := range req.Files {
		if !isPerformanceRelevantFile(file.Path) {
			continue
		}

		fileFindings := s.scanContent(file.NewBlob, file.Path)
		findings = append(findings, fileFindings...)
	}

	return findings, nil
}

func (s *PerformanceScanner) scanContent(content string, path string) []codereview.Finding {
	findings := make([]codereview.Finding, 0)

	for lineNum, line := range scanLines(content) {
		for _, rule := range s.rules {
			if rule.pattern.MatchString(line) {
				finding := codereview.Finding{
					ID:          generateID(),
					File:        path,
					Line:        lineNum,
					Type:        "performance",
					Rule:        rule.name,
					Severity:    rule.severity,
					Title:       rule.name,
					Message:     rule.message,
					Description: rule.message,
					Code:        line,
					Suggestion:  generatePerformanceSuggestion(rule.name),
					Tags:        []string{"performance", "optimization"},
				}
				findings = append(findings, finding)
			}
		}
	}

	return findings
}

// IsEnabled returns whether the scanner is enabled
func (s *PerformanceScanner) IsEnabled() bool {
	return s.enabled
}

// SetEnabled sets whether the scanner is enabled
func (s *PerformanceScanner) SetEnabled(enabled bool) {
	s.enabled = enabled
}

func isPerformanceRelevantFile(path string) bool {
	ext := getFileExtension(path)
	relevantExts := map[string]bool{
		"go": true, "js": true, "ts": true, "py": true, "java": true,
		"php": true, "rb": true, "cs": true, "cpp": true,
	}
	return relevantExts[ext]
}

func generatePerformanceSuggestion(ruleName string) string {
	suggestions := map[string]string{
		"N+1 Query Problem":        "Use batch queries or eager loading to fetch related data in a single query",
		"Unnecessary Loop":         "Consider using higher-order functions like map, filter, or reduce",
		"String Concatenation in Loop": "Use StringBuilder (Java), string.Builder (Go), or array join instead of string +=",
		"Synchronous I/O":          "Use async/await or callbacks to handle I/O without blocking the thread",
		"Memory Leak Potential":    "Always clean up timers with clearInterval/clearTimeout and remove event listeners",
		"Large Object Creation":    "Pre-allocate or use lazy initialization instead of creating large objects upfront",
	}
	if suggestion, ok := suggestions[ruleName]; ok {
		return suggestion
	}
	return "Review this code for performance opportunities"
}
