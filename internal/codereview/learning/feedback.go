package learning

import (
	"sync"
	"time"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// FeedbackCollector collects team feedback on suggestions
type FeedbackCollector struct {
	feedback map[string]*SuggestionFeedback
	mu       sync.RWMutex
}

// SuggestionFeedback represents feedback on a suggestion
type SuggestionFeedback struct {
	SuggestionID string
	FindingID    string
	Helpful      bool // true = helpful, false = not helpful
	Accepted     bool // whether fix was accepted
	Notes        string
	Timestamp    time.Time
	Author       string
}

// NewFeedbackCollector creates a new feedback collector
func NewFeedbackCollector() *FeedbackCollector {
	return &FeedbackCollector{
		feedback: make(map[string]*SuggestionFeedback),
	}
}

// RecordFeedback records feedback on a suggestion
func (fc *FeedbackCollector) RecordFeedback(feedback *SuggestionFeedback) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	feedback.Timestamp = time.Now()
	fc.feedback[feedback.SuggestionID] = feedback
}

// GetFeedback retrieves feedback for a suggestion
func (fc *FeedbackCollector) GetFeedback(suggestionID string) (*SuggestionFeedback, bool) {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	feedback, exists := fc.feedback[suggestionID]
	return feedback, exists
}

// GetAllFeedback returns all recorded feedback
func (fc *FeedbackCollector) GetAllFeedback() []*SuggestionFeedback {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	feedback := make([]*SuggestionFeedback, 0, len(fc.feedback))
	for _, f := range fc.feedback {
		feedback = append(feedback, f)
	}
	return feedback
}

// CalculateAcceptanceRate calculates the suggestion acceptance rate
func (fc *FeedbackCollector) CalculateAcceptanceRate() float64 {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	if len(fc.feedback) == 0 {
		return 0
	}

	accepted := 0
	for _, f := range fc.feedback {
		if f.Accepted {
			accepted++
		}
	}

	return float64(accepted) / float64(len(fc.feedback)) * 100
}

// CalculateHelpfulnessRate calculates how helpful suggestions are
func (fc *FeedbackCollector) CalculateHelpfulnessRate() float64 {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	if len(fc.feedback) == 0 {
		return 0
	}

	helpful := 0
	for _, f := range fc.feedback {
		if f.Helpful {
			helpful++
		}
	}

	return float64(helpful) / float64(len(fc.feedback)) * 100
}

// TeamFeedbackSummary provides team-level feedback metrics
type TeamFeedbackSummary struct {
	TotalFeedback      int
	AcceptanceRate     float64
	HelpfulnessRate    float64
	MostHelpfulRules   []string
	LeastHelpfulRules  []string
	CommonPatterns     []string
}

// GetTeamSummary generates a summary of team feedback
func (fc *FeedbackCollector) GetTeamSummary() *TeamFeedbackSummary {
	summary := &TeamFeedbackSummary{
		TotalFeedback:     len(fc.feedback),
		AcceptanceRate:    fc.CalculateAcceptanceRate(),
		HelpfulnessRate:   fc.CalculateHelpfulnessRate(),
		MostHelpfulRules:  make([]string, 0),
		LeastHelpfulRules: make([]string, 0),
	}

	return summary
}
