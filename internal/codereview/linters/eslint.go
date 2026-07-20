package linters

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ESLint linter for JavaScript/TypeScript
type ESLint struct {
	name    string
	version string
}

// NewESLint creates a new ESLint linter instance
func NewESLint() *ESLint {
	return &ESLint{
		name:    "eslint",
		version: "8.0.0+",
	}
}

func (e *ESLint) Name() string {
	return e.name
}

func (e *ESLint) Version() string {
	return e.version
}

func (e *ESLint) SupportedLanguages() []string {
	return []string{"javascript", "typescript", "jsx", "tsx"}
}

func (e *ESLint) IsAvailable(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "eslint", "--version")
	return cmd.Run() == nil
}

// eslintResult represents the JSON output from eslint
type eslintResult struct {
	FilePath string `json:"filePath"`
	Messages []struct {
		Line    int    `json:"line"`
		Column  int    `json:"column"`
		EndLine int    `json:"endLine"`
		EndCol  int    `json:"endColumn"`
		Message string `json:"message"`
		RuleID  string `json:"ruleId"`
		Severity int   `json:"severity"` // 1 = warning, 2 = error
	} `json:"messages"`
}

func (e *ESLint) Lint(ctx context.Context, files []string, config map[string]interface{}) ([]LintFinding, error) {
	if !e.IsAvailable(ctx) {
		return nil, fmt.Errorf("eslint is not installed")
	}

	args := []string{"--format", "json"}
	
	// Add custom config if provided
	if configPath, ok := config["config"].(string); ok {
		args = append(args, "--config", configPath)
	}

	// Filter for JS/TS files
	var jsFiles []string
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file))
		if ext == ".js" || ext == ".ts" || ext == ".jsx" || ext == ".tsx" {
			jsFiles = append(jsFiles, file)
		}
	}

	if len(jsFiles) == 0 {
		return nil, nil
	}

	args = append(args, jsFiles...)

	cmd := exec.CommandContext(ctx, "eslint", args...)
	output, _ := cmd.CombinedOutput()

	var results []eslintResult
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, fmt.Errorf("failed to parse eslint output: %w", err)
	}

	var findings []LintFinding
	for _, result := range results {
		for _, msg := range result.Messages {
			severity := "warning"
			if msg.Severity == 2 {
				severity = "error"
			}

			findings = append(findings, LintFinding{
				File:        result.FilePath,
				Line:        msg.Line,
				Column:      msg.Column,
				EndLine:     msg.EndLine,
				EndColumn:   msg.EndCol,
				Rule:        msg.RuleID,
				Message:     msg.Message,
				Severity:    severity,
				AutoFixable: strings.Contains(msg.Message, "autofix"),
				SourceLint:  "eslint",
			})
		}
	}

	return findings, nil
}

func (e *ESLint) CanAutoFix() bool {
	return true
}

func (e *ESLint) AutoFix(ctx context.Context, files []string, config map[string]interface{}) (map[string]string, error) {
	if !e.IsAvailable(ctx) {
		return nil, fmt.Errorf("eslint is not installed")
	}

	args := []string{"--fix"}

	if configPath, ok := config["config"].(string); ok {
		args = append(args, "--config", configPath)
	}

	// Filter for JS/TS files
	var jsFiles []string
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file))
		if ext == ".js" || ext == ".ts" || ext == ".jsx" || ext == ".tsx" {
			jsFiles = append(jsFiles, file)
		}
	}

	if len(jsFiles) == 0 {
		return nil, nil
	}

	args = append(args, jsFiles...)
	cmd := exec.CommandContext(ctx, "eslint", args...)
	if err := cmd.Run(); err != nil {
		// eslint returns non-zero on errors, but might have fixed some
	}

	// Read fixed files
	result := make(map[string]string)
	for _, file := range jsFiles {
		content, err := os.ReadFile(file)
		if err == nil {
			result[file] = string(content)
		}
	}

	return result, nil
}

func (e *ESLint) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"extends": "eslint:recommended",
		"rules": map[string]interface{}{
			"no-console": "warn",
			"no-unused-vars": "warn",
		},
	}
}
