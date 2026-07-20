package fixes

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// LinterFixer applies fixes from built-in linters
type LinterFixer struct {
}

// NewLinterFixer creates a new linter fixer
func NewLinterFixer() *LinterFixer {
	return &LinterFixer{}
}

// Fix applies a linter-based fix
func (lf *LinterFixer) Fix(ctx context.Context, req FixRequest) (*FixResult, error) {
	result := &FixResult{
		FindingID: req.Finding.ID,
		Strategy:  FixStrategyLinter,
	}

	// Read the file
	content, err := os.ReadFile(req.Finding.File)
	if err != nil {
		result.Error = fmt.Errorf("failed to read file: %w", err)
		return result, result.Error
	}

	result.Original = string(content)

	// Apply common fixes based on rule type
	fixed := lf.applyCommonFix(string(content), req.Finding)
	if fixed == "" {
		// No built-in fix available
		result.Success = false
		result.Error = fmt.Errorf("no linter fix available for rule: %s", req.Finding.Rule)
		return result, result.Error
	}

	result.Fixed = fixed
	result.Success = true
	result.CommitMessage = fmt.Sprintf("fix: apply %s suggestion", req.Finding.Rule)

	// Write the fixed content
	if err := os.WriteFile(req.Finding.File, []byte(fixed), 0644); err != nil {
		result.Success = false
		result.Error = fmt.Errorf("failed to write fixed file: %w", err)
		return result, result.Error
	}

	return result, nil
}

// applyCommonFix applies common, pattern-based fixes
func (lf *LinterFixer) applyCommonFix(content string, finding codereview.Finding) string {
	lines := strings.Split(content, "\n")

	if finding.Line < 1 || finding.Line > len(lines) {
		return ""
	}

	line := lines[finding.Line-1]

	// Common fixes
	switch {
	case strings.Contains(finding.Rule, "unused"):
		// Remove unused variable/import
		if strings.HasPrefix(strings.TrimSpace(line), "import") {
			lines = append(lines[:finding.Line-1], lines[finding.Line:]...)
		}

	case strings.Contains(finding.Rule, "trailing-comma") || strings.Contains(finding.Rule, "comma-dangle"):
		// Add or remove trailing commas
		if strings.HasSuffix(strings.TrimSpace(line), ",") {
			lines[finding.Line-1] = strings.TrimSuffix(line, ",")
		}

	case strings.Contains(finding.Rule, "semi"):
		// Add or remove semicolons
		trimmed := strings.TrimSpace(line)
		if strings.HasSuffix(trimmed, ";") {
			lines[finding.Line-1] = strings.TrimSuffix(line, ";")
		} else if !strings.HasSuffix(trimmed, "{") && !strings.HasSuffix(trimmed, "}") {
			lines[finding.Line-1] = line + ";"
		}

	case strings.Contains(finding.Rule, "no-console"):
		// Remove console statements
		if strings.Contains(line, "console.") {
			lines = append(lines[:finding.Line-1], lines[finding.Line:]...)
		}

	case strings.Contains(finding.Rule, "whitespace") || strings.Contains(finding.Rule, "indent"):
		// Fix indentation
		indent := lf.detectIndentation(line)
		trimmed := strings.TrimLeft(line, " \t")
		lines[finding.Line-1] = indent + trimmed

	case strings.Contains(finding.Rule, "space"):
		// Remove extra spaces
		lines[finding.Line-1] = strings.Join(strings.Fields(line), " ")

	default:
		return ""
	}

	return strings.Join(lines, "\n")
}

// detectIndentation detects the indentation level of a line
func (lf *LinterFixer) detectIndentation(line string) string {
	for i, ch := range line {
		if ch != ' ' && ch != '\t' {
			return line[:i]
		}
	}
	return ""
}
