package linters

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// Manager orchestrates linter execution
type Manager struct {
	registry *Registry
	mu       sync.RWMutex
}

// NewManager creates a new linter manager
func NewManager() *Manager {
	return &Manager{
		registry: NewRegistry(),
	}
}

// RegisterDefaultLinters registers all built-in linters
func (m *Manager) RegisterDefaultLinters() error {
	linters := []Linter{
		NewESLint(),
		NewPylint(),
		NewGolangciLint(),
	}

	for _, linter := range linters {
		if err := m.registry.Register(linter); err != nil {
			return fmt.Errorf("failed to register linter %s: %w", linter.Name(), err)
		}
	}

	return nil
}

// DetectLanguages analyzes files to determine what languages are present
func (m *Manager) DetectLanguages(files []string) map[string][]string {
	languages := make(map[string][]string)

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file))
		var lang string

		switch ext {
		case ".js", ".jsx", ".ts", ".tsx":
			lang = "javascript"
		case ".py":
			lang = "python"
		case ".go":
			lang = "go"
		case ".java":
			lang = "java"
		case ".rb":
			lang = "ruby"
		case ".cs":
			lang = "csharp"
		case ".rs":
			lang = "rust"
		case ".php":
			lang = "php"
		default:
			continue
		}

		if lang != "" {
			languages[lang] = append(languages[lang], file)
		}
	}

	return languages
}

// RunAll runs all available linters for the detected languages
func (m *Manager) RunAll(ctx context.Context, files []string, config map[string]interface{}) ([]codereview.Finding, error) {
	languages := m.DetectLanguages(files)
	if len(languages) == 0 {
		return nil, nil
	}

	var findings []codereview.Finding
	var mu sync.Mutex
	var wg sync.WaitGroup

	for lang, langFiles := range languages {
		linters, err := m.registry.GetForLanguage(ctx, lang)
		if err != nil || len(linters) == 0 {
			continue
		}

		for _, linter := range linters {
			wg.Add(1)
			go func(l Linter, lf []string) {
				defer wg.Done()

				lintFindings, err := l.Lint(ctx, lf, config)
				if err != nil {
					return // Skip on error
				}

				mu.Lock()
				for _, lf := range lintFindings {
					findings = append(findings, *lf.ConvertToFinding(l.Name()))
				}
				mu.Unlock()
			}(linter, langFiles)
		}
	}

	wg.Wait()
	return findings, nil
}

// RunLinter runs a specific linter
func (m *Manager) RunLinter(ctx context.Context, linterName string, files []string, config map[string]interface{}) ([]codereview.Finding, error) {
	linter, err := m.registry.Get(linterName)
	if err != nil {
		return nil, err
	}

	if !linter.IsAvailable(ctx) {
		return nil, fmt.Errorf("linter %s is not available", linterName)
	}

	lintFindings, err := linter.Lint(ctx, files, config)
	if err != nil {
		return nil, err
	}

	var findings []codereview.Finding
	for _, lf := range lintFindings {
		findings = append(findings, *lf.ConvertToFinding(linterName))
	}

	return findings, nil
}

// AutoFix applies auto-fixes from available linters
func (m *Manager) AutoFix(ctx context.Context, files []string, config map[string]interface{}) (map[string]string, error) {
	languages := m.DetectLanguages(files)
	results := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for lang, langFiles := range languages {
		linters, err := m.registry.GetForLanguage(ctx, lang)
		if err != nil || len(linters) == 0 {
			continue
		}

		for _, linter := range linters {
			if !linter.CanAutoFix() {
				continue
			}

			wg.Add(1)
			go func(l Linter, lf []string) {
				defer wg.Done()

				fixed, err := l.AutoFix(ctx, lf, config)
				if err != nil {
					return
				}

				mu.Lock()
				for file, content := range fixed {
					results[file] = content
				}
				mu.Unlock()
			}(linter, langFiles)
		}
	}

	wg.Wait()
	return results, nil
}

// GetAvailableLinters returns all available linters
func (m *Manager) GetAvailableLinters(ctx context.Context) ([]Linter, error) {
	return m.registry.GetAvailable(ctx)
}

// ListLinters returns all registered linters
func (m *Manager) ListLinters(ctx context.Context) map[string]bool {
	return m.registry.List(ctx)
}
