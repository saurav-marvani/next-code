package codereview

import (
	"time"

	"github.com/sauravmarvani/nextcode/internal/vcs"
)

// ReviewRequest represents a code review analysis request
type ReviewRequest struct {
	// Source information
	PullRequest *vcs.PullRequest
	Files       []vcs.FileDiff
	Commits     []vcs.CommitInfo

	// Configuration
	Policies   map[string]interface{}
	Scanners   []ScannerType
	Context    string // Project context for LLM
	TeamStyle  map[string]interface{}
}

// ReviewResult contains the complete code review analysis
type ReviewResult struct {
	ID          string
	RequestID   string
	Timestamp   time.Time
	PRNumber    int
	PRTitle     string
	Statistics  ReviewStats
	Findings    []Finding
	Suggestions []Suggestion
	Summary     string
	Metrics     ReviewMetrics
}

// ReviewStats contains high-level review statistics
type ReviewStats struct {
	TotalFiles         int
	TotalChanges       int
	AddedLines         int
	RemovedLines       int
	FindingCount       int
	CriticalCount      int
	HighCount          int
	MediumCount        int
	LowCount           int
	InfoCount          int
	AverageSeverity    float64
	ComplexityScore    float64
	RiskScore          float64
}

// Finding represents a code issue found during review
type Finding struct {
	ID          string
	File        string
	Line        int
	EndLine     int
	Column      int
	EndColumn   int
	Type        string // "security", "performance", "style", "correctness", "coverage"
	Rule        string // Specific rule that was violated
	Severity    Severity
	Title       string
	Description string
	Message     string
	Code        string // The actual code snippet
	Suggestion  string // How to fix it
	Remediation string // Detailed remediation steps
	Reference   string // Link to documentation
	Tags        []string
	Auto Fix    bool // Whether auto-fix is available
}

// Severity levels for findings
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

// ScannerType represents the type of scanner to run
type ScannerType string

const (
	ScannerSecurity     ScannerType = "security"
	ScannerPerformance  ScannerType = "performance"
	ScannerCorrectness  ScannerType = "correctness"
	ScannerStyle        ScannerType = "style"
	ScannerCoverage     ScannerType = "coverage"
)

// Suggestion represents a suggested code improvement
type Suggestion struct {
	ID          string
	FindingID   string
	Type        string // "inline-comment", "auto-fix", "refactoring", "test"
	Title       string
	Description string
	Body        string       // Markdown formatted suggestion
	Diff        string       // Unified diff of the suggestion
	Priority    int          // 1-10, 10 highest
	Accepted    bool
	Feedback    string       // User feedback
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ReviewMetrics contains metrics about the review
type ReviewMetrics struct {
	ReviewDuration    int // seconds
	AverageComments   float64
	SuggestionRate    float64
	AcceptanceRate    float64
	TechnicalDebt     float64 // 0-100
	TestCoverage      float64 // percentage
	ComplexityChange  float64
	SecurityRisk      float64 // 0-100
	PerformanceRisk   float64 // 0-100
}

// TeamLearning represents learned team patterns
type TeamLearning struct {
	ID           string
	TeamID       string
	Pattern      string // e.g., "preferred_error_handling"
	Occurrences  int
	Confidence   float64 // 0-1
	LastUpdated  time.Time
	Data         map[string]interface{}
}

// PolicyRule represents a code review policy rule
type PolicyRule struct {
	ID          string
	Name        string
	Description string
	Enabled     bool
	Severity    Severity
	Pattern     string // regex or matcher pattern
	Message     string
	Remediation string
	Tags        []string
}

// AnalysisContext contains context for analysis
type AnalysisContext struct {
	Language        string
	Framework       string
	ProjectType     string // web, api, cli, library, etc
	Dependencies    []string
	TestingFramework string
	LintingTools    []string
}
