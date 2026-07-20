package learning

import (
	"fmt"
	"sync"
	"time"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// PatternLearner learns team coding patterns
type PatternLearner struct {
	patterns   map[string]*codereview.TeamLearning
	ruleUsage  map[string]int
	feedback   *FeedbackCollector
	mu         sync.RWMutex
}

// NewPatternLearner creates a new pattern learner
func NewPatternLearner(feedback *FeedbackCollector) *PatternLearner {
	return &PatternLearner{
		patterns:   make(map[string]*codereview.TeamLearning),
		ruleUsage:  make(map[string]int),
		feedback:   feedback,
	}
}

// LearnFromResults learns patterns from review results
func (pl *PatternLearner) LearnFromResults(results []codereview.ReviewResult) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	for _, result := range results {
		// Track rule usage
		for _, finding := range result.Findings {
			pl.ruleUsage[finding.Rule]++
		}

		// Learn from acceptance patterns
		pl.learnAcceptancePatterns(result.Suggestions)
	}
}

// learnAcceptancePatterns learns which suggestions are accepted
func (pl *PatternLearner) learnAcceptancePatterns(suggestions []codereview.Suggestion) {
	for _, suggestion := range suggestions {
		feedback, exists := pl.feedback.GetFeedback(suggestion.ID)
		if !exists {
			continue
		}

		pattern := fmt.Sprintf("suggestion_%s_acceptance", suggestion.Type)
		learning := pl.patterns[pattern]

		if learning == nil {
			learning = &codereview.TeamLearning{
				ID:           generateID(),
				Pattern:      pattern,
				Occurrences:  0,
				Confidence:   0,
				LastUpdated:  time.Now(),
				Data:         make(map[string]interface{}),
			}
		}

		learning.Occurrences++
		if feedback.Accepted {
			learning.Data["accepted"].(int)++
		}

		if learning.Occurrences > 0 {
			accepted := learning.Data["accepted"].(int)
			learning.Confidence = float64(accepted) / float64(learning.Occurrences)
		}

		learning.LastUpdated = time.Now()
		pl.patterns[pattern] = learning
	}
}

// GetPattern retrieves a learned pattern
func (pl *PatternLearner) GetPattern(id string) (*codereview.TeamLearning, bool) {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	pattern, exists := pl.patterns[id]
	return pattern, exists
}

// ListPatterns returns all learned patterns
func (pl *PatternLearner) ListPatterns() []*codereview.TeamLearning {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	patterns := make([]*codereview.TeamLearning, 0, len(pl.patterns))
	for _, p := range pl.patterns {
		patterns = append(patterns, p)
	}
	return patterns
}

// AdjustSeverityByTeamLearning adjusts finding severity based on team patterns
func (pl *PatternLearner) AdjustSeverityByTeamLearning(findings []codereview.Finding) []codereview.Finding {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	adjusted := make([]codereview.Finding, len(findings))
	copy(adjusted, findings)

	for i, finding := range adjusted {
		// Check if this rule is commonly accepted
		pattern := fmt.Sprintf("rule_%s_acceptance", finding.Rule)
		if learning, exists := pl.patterns[pattern]; exists {
			// If confidence is very high that fixes are accepted, lower severity
			if learning.Confidence > 0.8 {
				adjusted[i].Severity = lowerSeverity(adjusted[i].Severity)
			}
		}
	}

	return adjusted
}

// lowerSeverity reduces severity by one level
func lowerSeverity(severity codereview.Severity) codereview.Severity {
	switch severity {
	case codereview.SeverityCritical:
		return codereview.SeverityHigh
	case codereview.SeverityHigh:
		return codereview.SeverityMedium
	case codereview.SeverityMedium:
		return codereview.SeverityLow
	case codereview.SeverityLow:
		return codereview.SeverityInfo
	default:
		return severity
	}
}

// GetMostAcceptedRules returns the most commonly accepted rules
func (pl *PatternLearner) GetMostAcceptedRules(limit int) []string {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	rules := make([]string, 0)
	// TODO: Sort by acceptance rate and return top N
	return rules
}

// GetLeastAcceptedRules returns the least commonly accepted rules
func (pl *PatternLearner) GetLeastAcceptedRules(limit int) []string {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	rules := make([]string, 0)
	// TODO: Sort by acceptance rate and return bottom N
	return rules
}

// helper function
func generateID() string {
	return fmt.Sprintf("learning_%d", time.Now().UnixNano())
}
