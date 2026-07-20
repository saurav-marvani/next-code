package fixes

import (
	"context"
	"fmt"
	"sync"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// FixStrategy defines how a finding can be fixed
type FixStrategy string

const (
	FixStrategyLinter   FixStrategy = "linter"      // Linter built-in fix
	FixStrategyAI       FixStrategy = "ai"           // AI-powered fix
	FixStrategyManual   FixStrategy = "manual"       // Manual fix required
	FixStrategySkip     FixStrategy = "skip"         // Cannot be fixed
)

// FixRequest represents a request to fix code issues
type FixRequest struct {
	Finding     codereview.Finding
	Strategy    FixStrategy
	Language    string
	Context     string // Additional context for LLM
	Temperature float32 // 0-1 for AI responses
}

// FixResult represents the result of a fix attempt
type FixResult struct {
	FindingID      string
	Original       string
	Fixed          string
	Explanation    string
	Strategy       FixStrategy
	Success        bool
	Error          error
	Confidence     float32 // 0-1 confidence score
	AIGenerated    bool
	CommitMessage  string
}

// Engine orchestrates code fixing
type Engine struct {
	linterFixer *LinterFixer
	aiFixer     *AIFixer
	mu          sync.RWMutex
}

// NewEngine creates a new fix engine
func NewEngine(aiFixer *AIFixer) *Engine {
	return &Engine{
		linterFixer: NewLinterFixer(),
		aiFixer:     aiFixer,
	}
}

// Fix applies a fix to a finding
func (e *Engine) Fix(ctx context.Context, req FixRequest) (*FixResult, error) {
	result := &FixResult{
		FindingID: req.Finding.ID,
		Strategy:  req.Strategy,
	}

	// Choose fix strategy
	switch req.Strategy {
	case FixStrategyLinter:
		return e.linterFixer.Fix(ctx, req)

	case FixStrategyAI:
		if e.aiFixer == nil {
			result.Success = false
			result.Error = fmt.Errorf("AI fixer not configured")
			return result, result.Error
		}
		return e.aiFixer.Fix(ctx, req)

	case FixStrategyManual:
		result.Success = false
		result.Error = fmt.Errorf("manual fix required")
		return result, nil

	case FixStrategySkip:
		result.Success = false
		result.Error = fmt.Errorf("fix skipped")
		return result, nil

	default:
		// Auto-select best strategy
		if req.Finding.AutoFix {
			return e.linterFixer.Fix(ctx, req)
		}
		if e.aiFixer != nil {
			return e.aiFixer.Fix(ctx, req)
		}
		result.Success = false
		result.Error = fmt.Errorf("no fix strategy available")
		return result, result.Error
	}
}

// FixMultiple applies fixes to multiple findings in parallel
func (e *Engine) FixMultiple(ctx context.Context, requests []FixRequest) ([]*FixResult, error) {
	results := make([]*FixResult, len(requests))
	var wg sync.WaitGroup

	for i, req := range requests {
		wg.Add(1)
		go func(idx int, r FixRequest) {
			defer wg.Done()
			result, _ := e.Fix(ctx, r)
			results[idx] = result
		}(i, req)
	}

	wg.Wait()
	return results, nil
}

// BatchCommit generates a commit message for a batch of fixes
func (e *Engine) BatchCommit(fixes []*FixResult) string {
	if len(fixes) == 0 {
		return ""
	}

	aiFixCount := 0
	linterFixCount := 0
	categories := make(map[string]int)

	for _, fix := range fixes {
		if fix.Success {
			if fix.AIGenerated {
				aiFixCount++
			} else {
				linterFixCount++
			}
			categories[fix.Strategy]++
		}
	}

	message := fmt.Sprintf("fix: apply code quality improvements\n\n")
	message += fmt.Sprintf("Applied %d automatic fixes:\n", len(fixes))

	if linterFixCount > 0 {
		message += fmt.Sprintf("- %d linter-based fixes\n", linterFixCount)
	}
	if aiFixCount > 0 {
		message += fmt.Sprintf("- %d AI-powered fixes\n", aiFixCount)
	}

	return message
}

// SelectStrategy intelligently selects the best fix strategy
func (e *Engine) SelectStrategy(ctx context.Context, finding codereview.Finding) FixStrategy {
	// Prefer linter fixes for style/simple issues
	if finding.AutoFix && (finding.Type == "style" || finding.Type == "linter") {
		return FixStrategyLinter
	}

	// Use AI for complex issues if available
	if e.aiFixer != nil && finding.Severity == codereview.SeverityHigh {
		return FixStrategyAI
	}

	// Default to manual for critical security issues
	if finding.Severity == codereview.SeverityCritical && finding.Type == "security" {
		return FixStrategyManual
	}

	// Prefer linter for available fixes
	if finding.AutoFix {
		return FixStrategyLinter
	}

	// Fall back to AI if available
	if e.aiFixer != nil {
		return FixStrategyAI
	}

	return FixStrategyManual
}
