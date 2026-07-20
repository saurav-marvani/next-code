package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// ComprehensiveReport generates detailed reports
type ComprehensiveReport struct {
	result *codereview.ReviewResult
}

// NewComprehensiveReport creates a new comprehensive reporter
func NewComprehensiveReport(result *codereview.ReviewResult) *ComprehensiveReport {
	return &ComprehensiveReport{result: result}
}

// GenerateStandupReport generates a daily standup summary
func (cr *ComprehensiveReport) GenerateStandupReport(prs []*codereview.ReviewResult) string {
	report := fmt.Sprintf("# Daily Standup Report - %s\n\n", time.Now().Format("2006-01-02"))

	// Summary stats
	totalPRs := len(prs)
	totalFindings := 0
	criticalCount := 0
	highCount := 0

	for _, pr := range prs {
		totalFindings += pr.Statistics.FindingCount
		criticalCount += pr.Statistics.CriticalCount
		highCount += pr.Statistics.HighCount
	}

	report += fmt.Sprintf("## Summary\n")
	report += fmt.Sprintf("- PRs Reviewed: %d\n", totalPRs)
	report += fmt.Sprintf("- Total Findings: %d\n", totalFindings)
	report += fmt.Sprintf("- Critical Issues: %d\n", criticalCount)
	report += fmt.Sprintf("- High Priority: %d\n\n", highCount)

	// By PR breakdown
	report += "## Reviews by PR\n\n"
	for i, pr := range prs {
		report += fmt.Sprintf("%d. **%s** (#%d)\n", i+1, pr.PRTitle, pr.PRNumber)
		report += fmt.Sprintf("   - Findings: %d (🔴 %d critical, 🟠 %d high)\n",
			pr.Statistics.FindingCount,
			pr.Statistics.CriticalCount,
			pr.Statistics.HighCount)
		report += fmt.Sprintf("   - Files: %d, Changes: +%d/-%d\n\n",
			pr.Statistics.TotalFiles,
			pr.Statistics.AddedLines,
			pr.Statistics.RemovedLines)
	}

	return report
}

// GenerateSprintReview generates a sprint review report
func (cr *ComprehensiveReport) GenerateSprintReview(prs []*codereview.ReviewResult, duration string) string {
	report := fmt.Sprintf("# Sprint Review Report\n")
	report += fmt.Sprintf("Period: %s\n\n", duration)

	// Metrics
	report += "## Code Quality Metrics\n"
	avgRisk := 0.0
	avgDebt := 0.0
	avgCoverage := 0.0

	for _, pr := range prs {
		avgRisk += float64(pr.Metrics.SecurityRisk)
		avgDebt += float64(pr.Metrics.TechnicalDebt)
		avgCoverage += pr.Metrics.TestCoverage
	}

	if len(prs) > 0 {
		avgRisk /= float64(len(prs))
		avgDebt /= float64(len(prs))
		avgCoverage /= float64(len(prs))
	}

	report += fmt.Sprintf("- Average Security Risk: %.1f%%\n", avgRisk)
	report += fmt.Sprintf("- Average Technical Debt: %.1f%%\n", avgDebt)
	report += fmt.Sprintf("- Average Test Coverage: %.1f%%\n\n", avgCoverage)

	// Top issues
	report += "## Top Issues\n\n"
	allIssues := make(map[string]int)
	for _, pr := range prs {
		for _, finding := range pr.Findings {
			allIssues[finding.Rule]++
		}
	}

	// Sort and display top 10
	type issueCount struct {
		rule  string
		count int
	}
	var issues []issueCount
	for rule, count := range allIssues {
		issues = append(issues, issueCount{rule, count})
	}
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].count > issues[j].count
	})

	for i, ic := range issues {
		if i >= 10 {
			break
		}
		report += fmt.Sprintf("%d. %s: %d occurrences\n", i+1, ic.rule, ic.count)
	}

	return report
}

// GenerateTeamAnalytics generates team analytics
func (cr *ComprehensiveReport) GenerateTeamAnalytics(prs []*codereview.ReviewResult, teamMembers []string) string {
	report := fmt.Sprintf("# Team Analytics Report\n\n")

	// Member contributions
	report += "## Team Contributions\n\n"
	for _, member := range teamMembers {
		memberPRs := 0
		memberFindings := 0
		for _, pr := range prs {
			if strings.Contains(pr.PRTitle, member) {
				memberPRs++
				memberFindings += pr.Statistics.FindingCount
			}
		}
		report += fmt.Sprintf("- **%s**: %d PRs, %d findings\n", member, memberPRs, memberFindings)
	}
	report += "\n"

	// Improvement trends
	report += "## Trends\n"
	if len(prs) > 1 {
		avg1 := float64(prs[0].Statistics.FindingCount)
		avgLast := float64(prs[len(prs)-1].Statistics.FindingCount)
		trend := "📈 Increasing"
		if avgLast < avg1 {
			trend = "📉 Decreasing"
		}
		report += fmt.Sprintf("- Finding Trend: %s\n", trend)
	}

	return report
}

// GenerateMarkdownReport generates a markdown report
func (cr *ComprehensiveReport) GenerateMarkdownReport() string {
	if cr.result == nil {
		return ""
	}

	report := fmt.Sprintf("# Code Review Report - PR #%d\n", cr.result.PRNumber)
	report += fmt.Sprintf("**%s**\n\n", cr.result.PRTitle)

	// Summary
	report += "## Summary\n"
	report += cr.result.Summary + "\n\n"

	// Statistics
	report += "## Statistics\n"
	report += fmt.Sprintf("| Metric | Value |\n")
	report += fmt.Sprintf("|--------|-------|\n")
	report += fmt.Sprintf("| Files Changed | %d |\n", cr.result.Statistics.TotalFiles)
	report += fmt.Sprintf("| Total Changes | %d |\n", cr.result.Statistics.TotalChanges)
	report += fmt.Sprintf("| Added Lines | %d |\n", cr.result.Statistics.AddedLines)
	report += fmt.Sprintf("| Removed Lines | %d |\n", cr.result.Statistics.RemovedLines)
	report += fmt.Sprintf("| Findings | %d |\n", cr.result.Statistics.FindingCount)
	report += fmt.Sprintf("| Security Risk | %.1f%% |\n", cr.result.Metrics.SecurityRisk)
	report += "\n"

	// Severity breakdown
	report += "## Findings by Severity\n"
	report += fmt.Sprintf("- 🔴 Critical: %d\n", cr.result.Statistics.CriticalCount)
	report += fmt.Sprintf("- 🟠 High: %d\n", cr.result.Statistics.HighCount)
	report += fmt.Sprintf("- 🟡 Medium: %d\n", cr.result.Statistics.MediumCount)
	report += fmt.Sprintf("- 🔵 Low: %d\n", cr.result.Statistics.LowCount)
	report += fmt.Sprintf("- ℹ️ Info: %d\n\n", cr.result.Statistics.InfoCount)

	// Top findings
	if len(cr.result.Findings) > 0 {
		report += "## Key Findings\n\n"

		// Sort by severity
		sort.Slice(cr.result.Findings, func(i, j int) bool {
			severityOrder := map[codereview.Severity]int{
				codereview.SeverityCritical: 0,
				codereview.SeverityHigh:     1,
				codereview.SeverityMedium:   2,
				codereview.SeverityLow:      3,
				codereview.SeverityInfo:     4,
			}
			return severityOrder[cr.result.Findings[i].Severity] < severityOrder[cr.result.Findings[j].Severity]
		})

		// Show top 10
		limit := 10
		if len(cr.result.Findings) < limit {
			limit = len(cr.result.Findings)
		}

		for i := 0; i < limit; i++ {
			finding := cr.result.Findings[i]
			report += fmt.Sprintf("### %d. %s\n", i+1, finding.Title)
			report += fmt.Sprintf("**Severity:** %s | **File:** %s:%d\n", finding.Severity, finding.File, finding.Line)
			report += fmt.Sprintf("**Rule:** %s\n", finding.Rule)
			report += fmt.Sprintf("**Message:** %s\n", finding.Message)
			if finding.Suggestion != "" {
				report += fmt.Sprintf("**Suggestion:** %s\n", finding.Suggestion)
			}
			report += "\n"
		}

		if len(cr.result.Findings) > limit {
			report += fmt.Sprintf("... and %d more findings\n\n", len(cr.result.Findings)-limit)
		}
	}

	// Suggestions
	if len(cr.result.Suggestions) > 0 {
		report += "## Suggested Improvements\n\n"
		for i, sug := range cr.result.Suggestions {
			if i >= 5 {
				report += fmt.Sprintf("... and %d more suggestions\n", len(cr.result.Suggestions)-5)
				break
			}
			report += fmt.Sprintf("%d. %s\n", i+1, sug.Title)
		}
		report += "\n"
	}

	// Review metadata
	report += "---\n"
	report += fmt.Sprintf("*Generated: %s*\n", cr.result.Timestamp.Format("2006-01-02 15:04:05"))
	report += fmt.Sprintf("*Review ID: %s*\n", cr.result.ID)

	return report
}
