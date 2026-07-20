package scanners

import (
	"context"
	"regexp"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// SecurityScanner performs security-focused code analysis
type SecurityScanner struct {
	enabled bool
	rules   []*securityRule
}

type securityRule struct {
	name     string
	pattern  *regexp.Regexp
	severity codereview.Severity
	message  string
}

// NewSecurityScanner creates a new security scanner
func NewSecurityScanner() *SecurityScanner {
	scanner := &SecurityScanner{
		enabled: true,
		rules:   make([]*securityRule, 0),
	}
	scanner.initializeRules()
	return scanner
}

func (s *SecurityScanner) initializeRules() {
	rules := []struct {
		name     string
		pattern  string
		severity codereview.Severity
		message  string
	}{
		{
			name:     "SQL Injection",
			pattern:  `(?i)(\+|concat|format)\s*.*(?:query|sql|execute)`,
			severity: codereview.SeverityCritical,
			message:  "Potential SQL injection vulnerability: use parameterized queries",
		},
		{
			name:     "Hardcoded Credentials",
			pattern:  `(?i)(password|api_?key|secret|token)\s*=\s*["\'][\w\-\.]+["\']`,
			severity: codereview.SeverityCritical,
			message:  "Hardcoded credentials detected: use environment variables or secrets manager",
		},
		{
			name:     "XSS Vulnerability",
			pattern:  `(?i)innerHTML\s*=|dangerouslySetInnerHTML`,
			severity: codereview.SeverityHigh,
			message:  "Potential XSS vulnerability: sanitize user input before rendering",
		},
		{
			name:     "Command Injection",
			pattern:  `(?i)(exec|system|sh\.Cmd)\s*\(\s*.*(?:user|input|request)`,
			severity: codereview.SeverityHigh,
			message:  "Potential command injection: use safe APIs, avoid shell execution with user input",
		},
		{
			name:     "Insecure Random",
			pattern:  `(?i)math\.random\(\)|rand\.Intn`,
			severity: codereview.SeverityHigh,
			message:  "Insecure random number generation: use crypto/rand for security-sensitive operations",
		},
	}

	for _, r := range rules {
		if pattern, err := regexp.Compile(r.pattern); err == nil {
			s.rules = append(s.rules, &securityRule{
				name:     r.name,
				pattern:  pattern,
				severity: r.severity,
				message:  r.message,
			})
		}
	}
}

// GetType returns the scanner type
func (s *SecurityScanner) GetType() codereview.ScannerType {
	return codereview.ScannerSecurity
}

// Scan performs security analysis on the code
func (s *SecurityScanner) Scan(ctx context.Context, req *codereview.ReviewRequest) ([]codereview.Finding, error) {
	findings := make([]codereview.Finding, 0)

	for _, file := range req.Files {
		if !isSecurityRelevantFile(file.Path) {
			continue
		}

		// Scan the new content
		fileFindings := s.scanContent(file.NewBlob, file.Path)
		findings = append(findings, fileFindings...)
	}

	return findings, nil
}

func (s *SecurityScanner) scanContent(content string, path string) []codereview.Finding {
	findings := make([]codereview.Finding, 0)

	for lineNum, line := range scanLines(content) {
		for _, rule := range s.rules {
			if rule.pattern.MatchString(line) {
				finding := codereview.Finding{
					ID:          generateID(),
					File:        path,
					Line:        lineNum,
					Type:        "security",
					Rule:        rule.name,
					Severity:    rule.severity,
					Title:       rule.name,
					Message:     rule.message,
					Description: rule.message,
					Code:        line,
					Suggestion:  generateSecuritySuggestion(rule.name),
					Tags:        []string{"security"},
				}
				findings = append(findings, finding)
			}
		}
	}

	return findings
}

// IsEnabled returns whether the scanner is enabled
func (s *SecurityScanner) IsEnabled() bool {
	return s.enabled
}

// SetEnabled sets whether the scanner is enabled
func (s *SecurityScanner) SetEnabled(enabled bool) {
	s.enabled = enabled
}

func isSecurityRelevantFile(path string) bool {
	// Check if file is code (not binary, docs, etc)
	ext := getFileExt(path)
	securityExts := map[string]bool{
		"go": true, "js": true, "ts": true, "py": true, "java": true,
		"php": true, "rb": true, "cs": true, "cpp": true, "c": true,
	}
	return securityExts[ext]
}

func getFileExt(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
	}
	return ""
}

func generateSecuritySuggestion(ruleName string) string {
	suggestions := map[string]string{
		"SQL Injection":        "Use parameterized queries or prepared statements to prevent SQL injection",
		"Hardcoded Credentials": "Move credentials to environment variables, configuration files, or a secrets manager",
		"XSS Vulnerability":    "Sanitize and escape user input before rendering, or use safe DOM APIs",
		"Command Injection":    "Avoid using shell execution with user input; use safe APIs or input validation",
		"Insecure Random":      "Use crypto/rand for cryptographic operations instead of math/rand",
	}
	if suggestion, ok := suggestions[ruleName]; ok {
		return suggestion
	}
	return "Review this code for potential security issues"
}
