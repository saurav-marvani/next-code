package scanners

import (
	"context"
	"regexp"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// CorrectnessScanner performs logic and correctness analysis
type CorrectnessScanner struct {
	enabled bool
	rules   []*correctnessRule
}

type correctnessRule struct {
	name     string
	pattern  *regexp.Regexp
	severity codereview.Severity
	message  string
}

// NewCorrectnessScanner creates a new correctness scanner
func NewCorrectnessScanner() *CorrectnessScanner {
	scanner := &CorrectnessScanner{
		enabled: true,
		rules:   make([]*correctnessRule, 0),
	}
	scanner.initializeRules()
	return scanner
}

func (s *CorrectnessScanner) initializeRules() {
	rules := []struct {
		name     string
		pattern  string
		severity codereview.Severity
		message  string
	}{
		{
			name:     "Null/Nil Dereference",
			pattern:  `(?i)(?:obj|result|data|response)\s*[\[\.]`,
			severity: codereview.SeverityHigh,
			message:  "Potential null/nil dereference: add null check before accessing",
		},
		{
			name:     "Missing Error Check",
			pattern:  `(?i)(?:err|error)\s*:=.*(?:json\.Unmarshal|ioutil\.ReadAll|os\.Open)`,
			severity: codereview.SeverityHigh,
			message:  "Error returned but not checked: add error handling",
		},
		{
			name:     "Missing Break in Switch",
			pattern:  `(?i)case\s+.*:(?:\n\s*[^}])*(?:\n\s*case|\n\s*default)`,
			severity: codereview.SeverityMedium,
			message:  "Missing break statement: add break or use fallthrough intentionally",
		},
		{
			name:     "Unreachable Code",
			pattern:  `(?i)(?:return|throw)\s*;[^}]*[a-zA-Z]`,
			severity: codereview.SeverityMedium,
			message:  "Code after return/throw statement is unreachable",
		},
		{
			name:     "Logic Error - Assignment vs Comparison",
			pattern:  `(?i)if\s*\([^)]*\s=\s[^=]`,
			severity: codereview.SeverityHigh,
			message:  "Assignment in condition (= instead of ==): verify this is intentional",
		},
		{
			name:     "Off-by-One Error",
			pattern:  `(?i)(?:i\s*<|i\s*<=)\s*(?:length|size|count)`,
			severity: codereview.SeverityMedium,
			message:  "Verify loop bounds to prevent off-by-one errors",
		},
	}

	for _, r := range rules {
		if pattern, err := regexp.Compile(r.pattern); err == nil {
			s.rules = append(s.rules, &correctnessRule{
				name:     r.name,
				pattern:  pattern,
				severity: r.severity,
				message:  r.message,
			})
		}
	}
}

// GetType returns the scanner type
func (s *CorrectnessScanner) GetType() codereview.ScannerType {
	return codereview.ScannerCorrectness
}

// Scan performs correctness analysis on the code
func (s *CorrectnessScanner) Scan(ctx context.Context, req *codereview.ReviewRequest) ([]codereview.Finding, error) {
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

func (s *CorrectnessScanner) scanContent(content string, path string) []codereview.Finding {
	findings := make([]codereview.Finding, 0)

	for lineNum, line := range scanLines(content) {
		for _, rule := range s.rules {
			if rule.pattern.MatchString(line) {
				finding := codereview.Finding{
					ID:          generateID(),
					File:        path,
					Line:        lineNum,
					Type:        "correctness",
					Rule:        rule.name,
					Severity:    rule.severity,
					Title:       rule.name,
					Message:     rule.message,
					Description: rule.message,
					Code:        line,
					Suggestion:  generateCorrectnessSuggestion(rule.name),
					Tags:        []string{"correctness", "logic"},
				}
				findings = append(findings, finding)
			}
		}
	}

	return findings
}

// IsEnabled returns whether the scanner is enabled
func (s *CorrectnessScanner) IsEnabled() bool {
	return s.enabled
}

// SetEnabled sets whether the scanner is enabled
func (s *CorrectnessScanner) SetEnabled(enabled bool) {
	s.enabled = enabled
}

func isCodeFile(path string) bool {
	ext := getFileExtension(path)
	codeExts := map[string]bool{
		"go": true, "js": true, "ts": true, "py": true, "java": true,
		"php": true, "rb": true, "cs": true, "cpp": true, "c": true,
	}
	return codeExts[ext]
}

func generateCorrectnessSuggestion(ruleName string) string {
	suggestions := map[string]string{
		"Null/Nil Dereference":           "Add null/nil check before accessing object properties or using values",
		"Missing Error Check":            "Check if err != nil after function calls that return errors",
		"Missing Break in Switch":        "Add break statement or use fallthrough comment if intentional",
		"Unreachable Code":               "Remove unreachable code after return/throw statements",
		"Logic Error - Assignment vs Comparison": "Use == for comparison, not = for assignment in conditions",
		"Off-by-One Error":               "Verify loop boundaries: < vs <= for array length to avoid off-by-one errors",
	}
	if suggestion, ok := suggestions[ruleName]; ok {
		return suggestion
	}
	return "Review this code for correctness issues"
}
