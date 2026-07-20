package linters

import (
	"context"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// Linter defines the interface for external linters
type Linter interface {
	// Name returns the name of the linter
	Name() string

	// Version returns the version of the linter
	Version() string

	// SupportedLanguages returns the list of languages this linter supports
	SupportedLanguages() []string

	// IsAvailable checks if the linter is installed and available
	IsAvailable(ctx context.Context) bool

	// Lint analyzes the given files and returns findings
	Lint(ctx context.Context, files []string, config map[string]interface{}) ([]LintFinding, error)

	// CanAutoFix returns whether this linter can auto-fix issues
	CanAutoFix() bool

	// AutoFix applies automatic fixes and returns the fixed files
	AutoFix(ctx context.Context, files []string, config map[string]interface{}) (map[string]string, error)

	// GetConfig returns the default configuration for this linter
	GetConfig() map[string]interface{}
}

// LintFinding represents a finding from a linter
type LintFinding struct {
	File        string
	Line        int
	Column      int
	EndLine     int
	EndColumn   int
	Rule        string
	Message     string
	Severity    string // error, warning, info
	Code        string
	Suggestion  string
	AutoFixable bool
	SourceLint  string // Name of the linter that reported this
}

// ConvertToFinding converts a linter finding to a code review finding
func (lf LintFinding) ConvertToFinding(linterName string) *codereview.Finding {
	severity := codereview.SeverityLow
	switch lf.Severity {
	case "error":
		severity = codereview.SeverityHigh
	case "warning":
		severity = codereview.SeverityMedium
	case "info":
		severity = codereview.SeverityInfo
	}

	return &codereview.Finding{
		File:        lf.File,
		Line:        lf.Line,
		Column:      lf.Column,
		EndLine:     lf.EndLine,
		EndColumn:   lf.EndColumn,
		Type:        "linter",
		Rule:        lf.Rule,
		Severity:    severity,
		Title:       lf.Rule,
		Description: lf.Message,
		Message:     lf.Message,
		Code:        lf.Code,
		Suggestion:  lf.Suggestion,
		AutoFix:     lf.AutoFixable,
		Tags:        []string{"linter", linterName},
	}
}
