package linters

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Pylint linter for Python
type Pylint struct {
	name    string
	version string
}

// NewPylint creates a new Pylint linter instance
func NewPylint() *Pylint {
	return &Pylint{
		name:    "pylint",
		version: "2.0.0+",
	}
}

func (p *Pylint) Name() string {
	return p.name
}

func (p *Pylint) Version() string {
	return p.version
}

func (p *Pylint) SupportedLanguages() []string {
	return []string{"python"}
}

func (p *Pylint) IsAvailable(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "pylint", "--version")
	return cmd.Run() == nil
}

// pylintResult represents a finding from pylint
type pylintResult struct {
	Path    string `json:"path"`
	Module  string `json:"module"`
	ObjType string `json:"type"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	EndLine int    `json:"endLine"`
	EndCol  int    `json:"endColumn"`
	Message string `json:"message"`
	Symbol  string `json:"symbol"`
	MsgID   string `json:"message-id"`
}

func (p *Pylint) Lint(ctx context.Context, files []string, config map[string]interface{}) ([]LintFinding, error) {
	if !p.IsAvailable(ctx) {
		return nil, fmt.Errorf("pylint is not installed")
	}

	// Filter for Python files
	var pyFiles []string
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file), ".py") {
			pyFiles = append(pyFiles, file)
		}
	}

	if len(pyFiles) == 0 {
		return nil, nil
	}

	args := []string{"--output-format=json"}

	if configPath, ok := config["config"].(string); ok {
		args = append(args, "--rcfile", configPath)
	}

	args = append(args, pyFiles...)

	cmd := exec.CommandContext(ctx, "pylint", args...)
	output, _ := cmd.CombinedOutput()

	var results []pylintResult
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, fmt.Errorf("failed to parse pylint output: %w", err)
	}

	var findings []LintFinding
	for _, result := range results {
		severity := "info"
		if strings.HasPrefix(result.MsgID, "E") {
			severity = "error"
		} else if strings.HasPrefix(result.MsgID, "W") {
			severity = "warning"
		}

		findings = append(findings, LintFinding{
			File:       result.Path,
			Line:       result.Line,
			Column:     result.Column,
			EndLine:    result.EndLine,
			EndColumn:  result.EndCol,
			Rule:       result.Symbol,
			Message:    result.Message,
			Severity:   severity,
			SourceLint: "pylint",
		})
	}

	return findings, nil
}

func (p *Pylint) CanAutoFix() bool {
	return false // pylint doesn't have built-in auto-fix
}

func (p *Pylint) AutoFix(ctx context.Context, files []string, config map[string]interface{}) (map[string]string, error) {
	return nil, fmt.Errorf("pylint does not support auto-fix")
}

func (p *Pylint) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"disable": []string{
			"missing-docstring",
			"too-many-arguments",
		},
		"max-line-length": 100,
	}
}
