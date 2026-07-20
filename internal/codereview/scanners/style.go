package scanners

import (
	"context"
	"regexp"
	"strings"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// StyleScanner performs code style and formatting analysis
type StyleScanner struct {
	enabled bool
	rules   []*styleRule
}

type styleRule struct {
	name     string
	pattern  *regexp.Regexp
	severity codereview.Severity
	message  string
}

// NewStyleScanner creates a new style scanner
func NewStyleScanner() *StyleScanner {
	scanner := &StyleScanner{
		enabled: true,
		rules:   make([]*styleRule, 0),
	}
	scanner.initializeRules()
	return scanner
}

func (s *StyleScanner) initializeRules() {
	rules := []struct {
		name     string
		pattern  string
		severity codereview.Severity
		message  string
	}{
		{
			name:     "Missing Function Documentation",
			pattern:  `(?i)^func\s+(?:[A-Z]\w*)\s*\(`,
			severity: codereview.SeverityLow,
			message:  "Public function missing documentation comment",
		},
		{
			name:     "Inconsistent Naming",
			pattern:  `(?i)(?:var|const)\s+(?:[a-z]{1,2}|CamelCase_with_underscores)`,
			severity: codereview.SeverityLow,
			message:  "Variable name doesn't follow convention: use camelCase or snake_case",
		},
		{
			name:     "Excessive Nesting",
			pattern:  `\s{16,}(?:if|for|while|function)`,
			severity: codereview.SeverityMedium,
			message:  "Excessive nesting (4+ levels): consider extracting methods or using guards",
		},
		{
			name:     "Magic Numbers",
			pattern:  `(?i)=[^"']*\b(?:0[1-9]|[1-9]\d+)\b(?![0-9])`,
			severity: codereview.SeverityLow,
			message:  "Magic number detected: define as a named constant",
		},
		{
			name:     "Trailing Whitespace",
			pattern:  `\s+$`,
			severity: codereview.SeverityInfo,
			message:  "Line contains trailing whitespace",
		},
		{
			name:     "Missing Blank Line",
			pattern:  `(?i)(?:^import|^const|^var|^func).*\n(?![\n]|$)`,
			severity: codereview.SeverityInfo,
			message:  "Missing blank line between logical sections",
		},
	}

	for _, r := range rules {
		if pattern, err := regexp.Compile(r.pattern); err == nil {
			s.rules = append(s.rules, &styleRule{
				name:     r.name,
				pattern:  pattern,
				severity: r.severity,
				message:  r.message,
			})
		}
	}
}

// GetType returns the scanner type
func (s *StyleScanner) GetType() codereview.ScannerType {
	return codereview.ScannerStyle
}

// Scan performs style analysis on the code
func (s *StyleScanner) Scan(ctx context.Context, req *codereview.ReviewRequest) ([]codereview.Finding, error) {
	findings := make([]codereview.Finding, 0)

	for _, file := range req.Files {
		if !isCodeFile(file.Path) {
			continue
		}

		fileFindings := s.scanContent(file.NewBlob, file.Path)
		findings = append(findings, fileFindings...)
	}

	return findings, nil
}

func (s *StyleScanner) scanContent(content string, path string) []codereview.Finding {
	findings := make([]codereview.Finding, 0)

	for lineNum, line := range scanLines(content) {
		// Check trailing whitespace
		if strings.HasSuffix(line, " ") || strings.HasSuffix(line, "\t") {
			findings = append(findings, codereview.Finding{
				ID:          generateID(),
				File:        path,
				Line:        lineNum,
				Type:        "style",
				Rule:        "Trailing Whitespace",
				Severity:    codereview.SeverityInfo,
				Title:       "Trailing Whitespace",
				Message:     "Line contains trailing whitespace",
				Code:        line,
				Suggestion:  "Remove trailing whitespace from end of line",
				Tags:        []string{"style", "formatting"},
			})
		}

		// Check for magic numbers
		if !strings.Contains(line, "//") && !strings.Contains(line, "const") {
			for _, rule := range s.rules {
				if rule.pattern.MatchString(line) {
					finding := codereview.Finding{
						ID:          generateID(),
						File:        path,
						Line:        lineNum,
						Type:        "style",
						Rule:        rule.name,
						Severity:    rule.severity,
						Title:       rule.name,
						Message:     rule.message,
						Code:        line,
						Suggestion:  generateStyleSuggestion(rule.name),
						Tags:        []string{"style"},
					}
					findings = append(findings, finding)
					break
				}
			}
		}
	}

	return findings
}

// IsEnabled returns whether the scanner is enabled
func (s *StyleScanner) IsEnabled() bool {
	return s.enabled
}

// SetEnabled sets whether the scanner is enabled
func (s *StyleScanner) SetEnabled(enabled bool) {
	s.enabled = enabled
}

func generateStyleSuggestion(ruleName string) string {
	suggestions := map[string]string{
		"Missing Function Documentation": "Add a comment above the function describing what it does",
		"Inconsistent Naming":            "Use consistent naming conventions: camelCase for variables, PascalCase for types",
		"Excessive Nesting":              "Extract nested blocks into separate functions to reduce nesting depth",
		"Magic Numbers":                  "Define named constants instead of using magic numbers",
		"Trailing Whitespace":            "Remove trailing whitespace from lines",
		"Missing Blank Line":             "Add blank lines between logical sections (imports, constants, functions)",
	}
	if suggestion, ok := suggestions[ruleName]; ok {
		return suggestion
	}
	return "Review this code for style improvements"
}

// Metrics for style analysis
type StyleMetrics struct {
	AverageLineLength      int
	AverageNestingDepth    float64
	CyclomaticComplexity   int
	DocumentationCoverage  float64
	StyleViolationCount    int
}
