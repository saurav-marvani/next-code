package linters

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// GolangciLint linter for Go
type GolangciLint struct {
	name    string
	version string
}

// NewGolangciLint creates a new golangci-lint linter instance
func NewGolangciLint() *GolangciLint {
	return &GolangciLint{
		name:    "golangci-lint",
		version: "1.50.0+",
	}
}

func (g *GolangciLint) Name() string {
	return g.name
}

func (g *GolangciLint) Version() string {
	return g.version
}

func (g *GolangciLint) SupportedLanguages() []string {
	return []string{"go"}
}

func (g *GolangciLint) IsAvailable(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "golangci-lint", "--version")
	return cmd.Run() == nil
}

// golangciResult represents a finding from golangci-lint
type golangciResult struct {
	Issues []struct {
		Text          string `json:"Text"`
		Pos           string `json:"Pos"`
		FromLinter    string `json:"FromLinter"`
		Replacement   *struct {
			NewLines []string `json:"NewLines"`
			Inline   *struct {
				StartCol int `json:"StartCol"`
				Length   int `json:"Length"`
				NewText  string `json:"NewText"`
			} `json:"Inline"`
		} `json:"Replacement"`
	} `json:"Issues"`
}

func (g *GolangciLint) Lint(ctx context.Context, files []string, config map[string]interface{}) ([]LintFinding, error) {
	if !g.IsAvailable(ctx) {
		return nil, fmt.Errorf("golangci-lint is not installed")
	}

	args := []string{"run", "--out-format", "json"}

	if configPath, ok := config["config"].(string); ok {
		args = append(args, "--config", configPath)
	}

	args = append(args, "./...")

	cmd := exec.CommandContext(ctx, "golangci-lint", args...)
	output, _ := cmd.CombinedOutput()

	var result golangciResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse golangci-lint output: %w", err)
	}

	var findings []LintFinding
	for _, issue := range result.Issues {
		// Parse position format: file:line:col
		parts := strings.Split(issue.Pos, ":")
		line, col := 0, 0
		file := ""
		if len(parts) >= 3 {
			file = parts[0]
			fmt.Sscanf(parts[1], "%d", &line)
			fmt.Sscanf(parts[2], "%d", &col)
		}

		severity := "warning"
		if strings.Contains(strings.ToLower(issue.FromLinter), "error") {
			severity = "error"
		}

		findings = append(findings, LintFinding{
			File:        file,
			Line:        line,
			Column:      col,
			Rule:        issue.FromLinter,
			Message:     issue.Text,
			Severity:    severity,
			AutoFixable: issue.Replacement != nil,
			SourceLint:  "golangci-lint",
		})
	}

	return findings, nil
}

func (g *GolangciLint) CanAutoFix() bool {
	return true
}

func (g *GolangciLint) AutoFix(ctx context.Context, files []string, config map[string]interface{}) (map[string]string, error) {
	return nil, fmt.Errorf("auto-fix for golangci-lint is handled by individual linters")
}

func (g *GolangciLint) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"linters": []string{
			"errcheck",
			"gosimple",
			"govet",
			"ineffassign",
			"staticcheck",
			"unused",
		},
	}
}
