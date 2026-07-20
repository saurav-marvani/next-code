package codereview

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sauravmarvani/nextcode/internal/vcs"
)

// Analyzer is the main code review analysis engine
type Analyzer struct {
	scanners map[ScannerType]Scanner
	policies *PolicyEngine
	mu       sync.RWMutex
}

// Scanner is the interface for analysis scanners
type Scanner interface {
	GetType() ScannerType
	Scan(ctx context.Context, req *ReviewRequest) ([]Finding, error)
	IsEnabled() bool
	SetEnabled(bool)
}

// PolicyEngine enforces review policies
type PolicyEngine struct {
	rules    map[string]*PolicyRule
	teamData map[string]interface{}
}

// NewAnalyzer creates a new code review analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		scanners: make(map[ScannerType]Scanner),
		policies: &PolicyEngine{
			rules:    make(map[string]*PolicyRule),
			teamData: make(map[string]interface{}),
		},
	}
}

// RegisterScanner registers a new scanner
func (a *Analyzer) RegisterScanner(scanner Scanner) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.scanners[scanner.GetType()] = scanner
}

// Analyze performs a complete code review analysis
func (a *Analyzer) Analyze(ctx context.Context, req *ReviewRequest) (*ReviewResult, error) {
	start := time.Now()

	if req == nil {
		return nil, fmt.Errorf("review request is required")
	}

	if req.PullRequest == nil {
		return nil, fmt.Errorf("pull request information is required")
	}

	result := &ReviewResult{
		ID:        generateID(),
		RequestID: generateID(),
		Timestamp: time.Now(),
		PRNumber:  req.PullRequest.Number,
		PRTitle:   req.PullRequest.Title,
		Findings:  []Finding{},
		Suggestions: []Suggestion{},
	}

	// Run all enabled scanners in parallel
	findings, err := a.runScanners(ctx, req)
	if err != nil {
		return nil, err
	}

	result.Findings = findings

	// Apply team policies
	result.Findings = a.policies.ApplyPolicies(result.Findings)

	// Generate suggestions from findings
	result.Suggestions = a.generateSuggestions(result.Findings)

	// Calculate statistics
	result.Statistics = calculateStats(req, result.Findings)

	// Generate summary
	result.Summary = generateSummary(result)

	// Calculate metrics
	result.Metrics = calculateMetrics(result, time.Since(start))

	return result, nil
}

// runScanners executes all enabled scanners
func (a *Analyzer) runScanners(ctx context.Context, req *ReviewRequest) ([]Finding, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var findings []Finding
	var mu sync.Mutex
	var wg sync.WaitGroup
	var scanErr error

	for _, scanner := range a.scanners {
		if !scanner.IsEnabled() {
			continue
		}

		// Check if scanner type is requested
		if len(req.Scanners) > 0 {
			requested := false
			for _, t := range req.Scanners {
				if t == scanner.GetType() {
					requested = true
					break
				}
			}
			if !requested {
				continue
			}
		}

		wg.Add(1)
		go func(s Scanner) {
			defer wg.Done()

			scanFindings, err := s.Scan(ctx, req)
			if err != nil {
				// Log error but continue with other scanners
				fmt.Printf("Scanner %s error: %v\n", s.GetType(), err)
				return
			}

			mu.Lock()
			findings = append(findings, scanFindings...)
			mu.Unlock()
		}(scanner)
	}

	wg.Wait()

	if scanErr != nil {
		return nil, scanErr
	}

	return findings, nil
}

// generateSuggestions converts findings into actionable suggestions
func (a *Analyzer) generateSuggestions(findings []Finding) []Suggestion {
	suggestions := make([]Suggestion, 0)

	for _, finding := range findings {
		if finding.Suggestion != "" {
			suggestion := Suggestion{
				ID:          generateID(),
				FindingID:   finding.ID,
				Type:        "inline-comment",
				Title:       finding.Title,
				Description: finding.Description,
				Body:        formatSuggestionBody(finding),
				Priority:    priorityFromSeverity(finding.Severity),
				CreatedAt:   time.Now(),
			}
			suggestions = append(suggestions, suggestion)
		}
	}

	return suggestions
}

// AnalyzeLocalDiff analyzes a local diff/branch
func (a *Analyzer) AnalyzeLocalDiff(ctx context.Context, diffContent string, policies map[string]interface{}) (*ReviewResult, error) {
	// Parse diff
	parser := vcs.NewDiffParser(nil)
	_ = parser // TODO: implement local diff parsing

	// Create synthetic PR
	pr := &vcs.PullRequest{
		ID:          "local",
		Title:       "Local Changes",
		State:       "open",
		Author:      vcs.User{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	req := &ReviewRequest{
		PullRequest: pr,
		Policies:    policies,
	}

	return a.Analyze(ctx, req)
}

// Helper functions

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func priorityFromSeverity(severity Severity) int {
	switch severity {
	case SeverityCritical:
		return 10
	case SeverityHigh:
		return 8
	case SeverityMedium:
		return 5
	case SeverityLow:
		return 2
	case SeverityInfo:
		return 1
	default:
		return 5
	}
}

func formatSuggestionBody(finding Finding) string {
	return fmt.Sprintf("**%s**: %s\n\n%s", finding.Title, finding.Description, finding.Suggestion)
}

func calculateStats(req *ReviewRequest, findings []Finding) ReviewStats {
	stats := ReviewStats{
		TotalFiles: len(req.Files),
		TotalChanges: 0,
	}

	for _, file := range req.Files {
		stats.TotalChanges += file.Changes
		stats.AddedLines += file.Additions
		stats.RemovedLines += file.Deletions
	}

	for _, finding := range findings {
		stats.FindingCount++
		switch finding.Severity {
		case SeverityCritical:
			stats.CriticalCount++
		case SeverityHigh:
			stats.HighCount++
		case SeverityMedium:
			stats.MediumCount++
		case SeverityLow:
			stats.LowCount++
		case SeverityInfo:
			stats.InfoCount++
		}
	}

	// Calculate average severity
	if stats.FindingCount > 0 {
		total := stats.CriticalCount*5 + stats.HighCount*4 + stats.MediumCount*3 + stats.LowCount*2 + stats.InfoCount*1
		stats.AverageSeverity = float64(total) / float64(stats.FindingCount)
	}

	return stats
}

func generateSummary(result *ReviewResult) string {
	summary := fmt.Sprintf("PR #%d: %s\n\n", result.PRNumber, result.PRTitle)
	summary += fmt.Sprintf("**Review Summary:**\n")
	summary += fmt.Sprintf("- Files: %d\n", result.Statistics.TotalFiles)
	summary += fmt.Sprintf("- Changes: +%d -%d\n", result.Statistics.AddedLines, result.Statistics.RemovedLines)
	summary += fmt.Sprintf("- Issues Found: %d\n", result.Statistics.FindingCount)

	if result.Statistics.CriticalCount > 0 {
		summary += fmt.Sprintf("  - Critical: %d\n", result.Statistics.CriticalCount)
	}
	if result.Statistics.HighCount > 0 {
		summary += fmt.Sprintf("  - High: %d\n", result.Statistics.HighCount)
	}

	return summary
}

func calculateMetrics(result *ReviewResult, duration time.Duration) ReviewMetrics {
	return ReviewMetrics{
		ReviewDuration: int(duration.Seconds()),
		SecurityRisk:   calculateSecurityRisk(result),
		PerformanceRisk: calculatePerformanceRisk(result),
	}
}

func calculateSecurityRisk(result *ReviewResult) float64 {
	risk := 0.0
	for _, finding := range result.Findings {
		if finding.Type == "security" {
			switch finding.Severity {
			case SeverityCritical:
				risk += 25
			case SeverityHigh:
				risk += 15
			case SeverityMedium:
				risk += 8
			case SeverityLow:
				risk += 2
			}
		}
	}
	if risk > 100 {
		risk = 100
	}
	return risk
}

func calculatePerformanceRisk(result *ReviewResult) float64 {
	risk := 0.0
	for _, finding := range result.Findings {
		if finding.Type == "performance" {
			switch finding.Severity {
			case SeverityCritical:
				risk += 20
			case SeverityHigh:
				risk += 12
			case SeverityMedium:
				risk += 6
			case SeverityLow:
				risk += 2
			}
		}
	}
	if risk > 100 {
		risk = 100
	}
	return risk
}
