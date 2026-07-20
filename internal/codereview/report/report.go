package report

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// Reporter generates code review reports
type Reporter struct {
	result *codereview.ReviewResult
}

// NewReporter creates a new reporter
func NewReporter(result *codereview.ReviewResult) *Reporter {
	return &Reporter{
		result: result,
	}
}

// GenerateMarkdown generates a markdown report
func (r *Reporter) GenerateMarkdown() string {
	report := ""

	// Title
	report += fmt.Sprintf("# Code Review Report\n")
	report += fmt.Sprintf("**PR:** #%d - %s\n\n", r.result.PRNumber, r.result.PRTitle)

	// Summary
	report += r.generateSummary()

	// Statistics
	report += r.generateStatistics()

	// Findings by severity
	report += r.generateFindingsByType()

	// Suggestions
	report += r.generateSuggestions()

	// Metrics
	report += r.generateMetrics()

	return report
}

// generateSummary generates the summary section
func (r *Reporter) generateSummary() string {
	summary := "## Summary\n\n"

	if r.result.Summary != "" {
		summary += r.result.Summary + "\n\n"
	}

	summary += "| Metric | Value |\n"
	summary += "|--------|-------|\n"
	summary += fmt.Sprintf("| Files Changed | %d |\n", r.result.Statistics.TotalFiles)
	summary += fmt.Sprintf("| Total Changes | +%d -%d |\n", r.result.Statistics.AddedLines, r.result.Statistics.RemovedLines)
	summary += fmt.Sprintf("| Issues Found | %d |\n", r.result.Statistics.FindingCount)
	summary += fmt.Sprintf("| Critical | %d |\n", r.result.Statistics.CriticalCount)
	summary += fmt.Sprintf("| High | %d |\n", r.result.Statistics.HighCount)
	summary += fmt.Sprintf("| Medium | %d |\n", r.result.Statistics.MediumCount)

	summary += "\n"
	return summary
}

// generateStatistics generates the statistics section
func (r *Reporter) generateStatistics() string {
	stats := "## Statistics\n\n"

	stats += "### Finding Distribution\n\n"
	stats += fmt.Sprintf("- **Critical:** %d\n", r.result.Statistics.CriticalCount)
	stats += fmt.Sprintf("- **High:** %d\n", r.result.Statistics.HighCount)
	stats += fmt.Sprintf("- **Medium:** %d\n", r.result.Statistics.MediumCount)
	stats += fmt.Sprintf("- **Low:** %d\n", r.result.Statistics.LowCount)
	stats += fmt.Sprintf("- **Info:** %d\n\n", r.result.Statistics.InfoCount)

	stats += "### Code Changes\n\n"
	stats += fmt.Sprintf("- **Files:** %d\n", r.result.Statistics.TotalFiles)
	stats += fmt.Sprintf("- **Added:** %d lines\n", r.result.Statistics.AddedLines)
	stats += fmt.Sprintf("- **Removed:** %d lines\n", r.result.Statistics.RemovedLines)
	stats += fmt.Sprintf("- **Total Changes:** %d\n\n", r.result.Statistics.TotalChanges)

	return stats
}

// generateFindingsByType generates findings organized by type
func (r *Reporter) generateFindingsByType() string {
	report := "## Findings\n\n"

	// Group findings by type
	findingsByType := make(map[string][]codereview.Finding)
	for _, finding := range r.result.Findings {
		findingsByType[finding.Type] = append(findingsByType[finding.Type], finding)
	}

	// Report findings by type
	for _, ftype := range []string{"security", "performance", "correctness", "style"} {
		findings := findingsByType[ftype]
		if len(findings) == 0 {
			continue
		}

		report += fmt.Sprintf("### %s (%d)\n\n", strings.Title(ftype), len(findings))

		// Sort by severity
		sort.Slice(findings, func(i, j int) bool {
			return severityRank(findings[i].Severity) > severityRank(findings[j].Severity)
		})

		for _, finding := range findings {
			report += fmt.Sprintf("#### %s [%s]\n", finding.Title, finding.Severity)
			report += fmt.Sprintf("**File:** %s:%d\n", finding.File, finding.Line)
			report += fmt.Sprintf("**Rule:** %s\n", finding.Rule)
			report += fmt.Sprintf("**Message:** %s\n\n", finding.Message)

			if finding.Code != "" {
				report += fmt.Sprintf("```\n%s\n```\n\n", finding.Code)
			}

			if finding.Suggestion != "" {
				report += fmt.Sprintf("**Suggestion:** %s\n\n", finding.Suggestion)
			}

			report += "\n"
		}
	}

	return report
}

// generateSuggestions generates the suggestions section
func (r *Reporter) generateSuggestions() string {
	if len(r.result.Suggestions) == 0 {
		return ""
	}

	report := "## Suggestions\n\n"

	// Sort suggestions by priority
	sort.Slice(r.result.Suggestions, func(i, j int) bool {
		return r.result.Suggestions[i].Priority > r.result.Suggestions[j].Priority
	})

	for i, suggestion := range r.result.Suggestions {
		if i >= 10 {
			report += fmt.Sprintf("\n... and %d more suggestions\n", len(r.result.Suggestions)-10)
			break
		}

		report += fmt.Sprintf("### %s\n", suggestion.Title)
		report += fmt.Sprintf("**Type:** %s\n", suggestion.Type)
		report += fmt.Sprintf("**Priority:** %d/10\n\n", suggestion.Priority)
		report += fmt.Sprintf("%s\n\n", suggestion.Body)
	}

	return report
}

// generateMetrics generates the metrics section
func (r *Reporter) generateMetrics() string {
	metrics := "## Metrics\n\n"

	m := r.result.Metrics

	metrics += fmt.Sprintf("- **Review Duration:** %d seconds\n", m.ReviewDuration)
	metrics += fmt.Sprintf("- **Security Risk:** %.1f%%\n", m.SecurityRisk)
	metrics += fmt.Sprintf("- **Performance Risk:** %.1f%%\n", m.PerformanceRisk)

	if m.TechnicalDebt > 0 {
		metrics += fmt.Sprintf("- **Technical Debt:** %.1f%%\n", m.TechnicalDebt)
	}

	if m.TestCoverage > 0 {
		metrics += fmt.Sprintf("- **Test Coverage:** %.1f%%\n", m.TestCoverage)
	}

	metrics += "\n"
	return metrics
}

// severityRank returns a numeric rank for sorting
func severityRank(severity codereview.Severity) int {
	switch severity {
	case codereview.SeverityCritical:
		return 5
	case codereview.SeverityHigh:
		return 4
	case codereview.SeverityMedium:
		return 3
	case codereview.SeverityLow:
		return 2
	case codereview.SeverityInfo:
		return 1
	default:
		return 0
	}
}

// GenerateJSON generates a JSON report (placeholder)
func (r *Reporter) GenerateJSON() string {
	// This would be implemented using encoding/json
	return ""
}

// GenerateHTML generates an HTML report (placeholder)
func (r *Reporter) GenerateHTML() string {
	// This would be implemented with proper HTML formatting
	return ""
}
