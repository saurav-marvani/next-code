package scanners

import (
	"context"
	"regexp"
	"strings"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// AdvancedScanner detects complex bugs that simple patterns can't catch
type AdvancedScanner struct {
	enabled bool
}

// NewAdvancedScanner creates a new advanced scanner
func NewAdvancedScanner() *AdvancedScanner {
	return &AdvancedScanner{enabled: true}
}

func (a *AdvancedScanner) Scan(ctx context.Context, req *codereview.ReviewRequest) ([]codereview.Finding, error) {
	var findings []codereview.Finding

	// Scan for race conditions
	findings = append(findings, a.detectRaceConditions(req.Files)...)

	// Scan for resource leaks
	findings = append(findings, a.detectResourceLeaks(req.Files)...)

	// Scan for logic errors
	findings = append(findings, a.detectLogicErrors(req.Files)...)

	// Scan for concurrency issues
	findings = append(findings, a.detectConcurrencyIssues(req.Files)...)

	return findings, nil
}

func (a *AdvancedScanner) detectRaceConditions(files []codereview.FileDiff) []codereview.Finding {
	var findings []codereview.Finding

	// Pattern: Unsynchronized access to shared variables
	for _, file := range files {
		lines := strings.Split(file.Patch, "\n")
		for i, line := range lines {
			// Look for go routines accessing shared state without locks
			if strings.Contains(line, "go ") && i > 0 {
				// Check if mutex lock is present nearby
				context := strings.Join(lines[max(0, i-5):min(len(lines), i+5)], "\n")
				if !strings.Contains(context, "Lock()") && !strings.Contains(context, "sync") {
					findings = append(findings, codereview.Finding{
						File:        file.Path,
						Line:        i,
						Type:        "correctness",
						Rule:        "potential-race-condition",
						Severity:    codereview.SeverityHigh,
						Title:       "Potential Race Condition",
						Description: "Goroutine accesses shared state without synchronization",
						Message:     "Unsynchronized goroutine access detected",
						Tags:        []string{"concurrency", "race-condition"},
					})
				}
			}
		}
	}

	return findings
}

func (a *AdvancedScanner) detectResourceLeaks(files []codereview.FileDiff) []codereview.Finding {
	var findings []codereview.Finding

	for _, file := range files {
		// Pattern: File/Database/Network opens without closes
		patterns := []*regexp.Regexp{
			regexp.MustCompile(`(?i)open.*\(.*\)\s*(?:,|:=)`),
			regexp.MustCompile(`(?i)db\.query\(`),
			regexp.MustCompile(`(?i)net\.dial\(`),
		}

		lines := strings.Split(file.Patch, "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, "+") { // Only new lines
				for _, pattern := range patterns {
					if pattern.MatchString(line) {
						// Check if defer close/cleanup exists
						if !strings.Contains(strings.Join(lines[i:min(len(lines), i+10)], "\n"), "defer") &&
							!strings.Contains(strings.Join(lines[i:min(len(lines), i+10)], "\n"), "Close()") {
							findings = append(findings, codereview.Finding{
								File:        file.Path,
								Line:        i,
								Type:        "correctness",
								Rule:        "resource-leak",
								Severity:    codereview.SeverityMedium,
								Title:       "Potential Resource Leak",
								Description: "Resource acquired without guaranteed cleanup",
								Message:     "Missing defer or cleanup for resource",
								Tags:        []string{"resource-leak", "memory"},
							})
						}
					}
				}
			}
		}
	}

	return findings
}

func (a *AdvancedScanner) detectLogicErrors(files []codereview.FileDiff) []codereview.Finding {
	var findings []codereview.Finding

	for _, file := range files {
		lines := strings.Split(file.Patch, "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, "+") {
				// Off-by-one errors in loops
				if strings.Contains(line, "< len") || strings.Contains(line, "> len") {
					if strings.Contains(line, "for i := 0") {
						findings = append(findings, codereview.Finding{
							File:        file.Path,
							Line:        i,
							Type:        "correctness",
							Rule:        "loop-bound-error",
							Severity:    codereview.SeverityMedium,
							Title:       "Possible Loop Boundary Error",
							Description: "Loop boundary might cause off-by-one error",
							Message:     "Review loop boundaries carefully",
							Tags:        []string{"loop", "off-by-one"},
						})
					}
				}

				// Unreachable code after return
				if strings.Contains(line, "return") {
					if i+1 < len(lines) && !strings.HasPrefix(lines[i+1], "}") {
						nextLine := strings.TrimSpace(lines[i+1])
						if nextLine != "" && !strings.HasPrefix(nextLine, "//") {
							findings = append(findings, codereview.Finding{
								File:        file.Path,
								Line:        i + 1,
								Type:        "correctness",
								Rule:        "unreachable-code",
								Severity:    codereview.SeverityLow,
								Title:       "Unreachable Code",
								Description: "Code after return statement",
								Message:     "This code is unreachable",
								Tags:        []string{"dead-code"},
							})
						}
					}
				}
			}
		}
	}

	return findings
}

func (a *AdvancedScanner) detectConcurrencyIssues(files []codereview.FileDiff) []codereview.Finding {
	var findings []codereview.Finding

	for _, file := range files {
		lines := strings.Split(file.Patch, "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, "+") {
				// Context timeouts not checked
				if strings.Contains(line, "ctx") && strings.Contains(line, "Done()") {
					if !strings.Contains(strings.Join(lines[max(0, i-3):min(len(lines), i+3)], "\n"), "select") {
						findings = append(findings, codereview.Finding{
							File:        file.Path,
							Line:        i,
							Type:        "correctness",
							Rule:        "context-timeout-not-checked",
							Severity:    codereview.SeverityMedium,
							Title:       "Context Timeout Not Handled",
							Description: "Context timeout might not be properly handled",
							Message:     "Ensure context.Done() is checked in select statement",
							Tags:        []string{"context", "concurrency"},
						})
					}
				}
			}
		}
	}

	return findings
}

func (a *AdvancedScanner) SetEnabled(enabled bool) {
	a.enabled = enabled
}

func (a *AdvancedScanner) Enabled() bool {
	return a.enabled
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
