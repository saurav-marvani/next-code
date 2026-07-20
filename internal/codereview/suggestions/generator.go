package suggestions

import (
	"fmt"
	"strings"
	"time"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// Generator creates suggestions from findings
type Generator struct {
	templates map[string]string
}

// NewGenerator creates a new suggestion generator
func NewGenerator() *Generator {
	return &Generator{
		templates: initializeTemplates(),
	}
}

// Generate creates suggestions from findings
func (g *Generator) Generate(findings []codereview.Finding) []codereview.Suggestion {
	suggestions := make([]codereview.Suggestion, 0)

	for _, finding := range findings {
		suggestion := g.generateForFinding(finding)
		if suggestion != nil {
			suggestions = append(suggestions, *suggestion)
		}
	}

	return suggestions
}

// generateForFinding creates a suggestion for a single finding
func (g *Generator) generateForFinding(finding codereview.Finding) *codereview.Suggestion {
	if finding.Suggestion == "" {
		return nil
	}

	suggestion := &codereview.Suggestion{
		ID:          fmt.Sprintf("sugg_%d", time.Now().UnixNano()),
		FindingID:   finding.ID,
		Type:        determineSuggestionType(finding),
		Title:       finding.Title,
		Description: finding.Description,
		Body:        formatSuggestionBody(finding),
		Priority:    priorityFromSeverity(finding.Severity),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return suggestion
}

// determineSuggestionType determines the type of suggestion
func determineSuggestionType(finding codereview.Finding) string {
	if finding.AutoFix {
		return "auto-fix"
	}

	switch finding.Type {
	case "security":
		return "security-fix"
	case "performance":
		return "optimization"
	case "style":
		return "refactoring"
	default:
		return "inline-comment"
	}
}

// formatSuggestionBody creates markdown-formatted suggestion body
func formatSuggestionBody(finding codereview.Finding) string {
	body := fmt.Sprintf("### %s\n\n", finding.Title)

	if finding.Severity != "" {
		body += fmt.Sprintf("**Severity:** %s\n", finding.Severity)
	}

	body += fmt.Sprintf("**Location:** %s:%d\n\n", finding.File, finding.Line)

	body += fmt.Sprintf("**Issue:** %s\n\n", finding.Description)

	body += fmt.Sprintf("**Current Code:**\n```\n%s\n```\n\n", finding.Code)

	if finding.Suggestion != "" {
		body += fmt.Sprintf("**Suggestion:**\n%s\n\n", finding.Suggestion)
	}

	if finding.Reference != "" {
		body += fmt.Sprintf("[Learn More](%s)\n", finding.Reference)
	}

	return body
}

// priorityFromSeverity maps severity to priority
func priorityFromSeverity(severity codereview.Severity) int {
	switch severity {
	case codereview.SeverityCritical:
		return 10
	case codereview.SeverityHigh:
		return 8
	case codereview.SeverityMedium:
		return 5
	case codereview.SeverityLow:
		return 2
	case codereview.SeverityInfo:
		return 1
	default:
		return 5
	}
}

// initializeTemplates initializes suggestion templates
func initializeTemplates() map[string]string {
	return map[string]string{
		"sql_injection": "Use parameterized queries to prevent SQL injection:\n\n**Before:**\n```sql\nquery := fmt.Sprintf(\"SELECT * FROM users WHERE id = %d\", userID)\n```\n\n**After:**\n```sql\nquery := \"SELECT * FROM users WHERE id = ?\"\ndb.Query(query, userID)\n```",

		"xss_vulnerability": "Sanitize user input before rendering:\n\n**Before:**\n```javascript\ndiv.innerHTML = userContent;\n```\n\n**After:**\n```javascript\nconst div = document.createElement('div');\ndiv.textContent = userContent;\n```",

		"hardcoded_credentials": "Move credentials to environment variables:\n\n**Before:**\n```go\nconst apiKey = \"sk_live_abc123xyz\"\n```\n\n**After:**\n```go\napiKey := os.Getenv(\"API_KEY\")\n```",

		"n_plus_one_query": "Use batch queries or eager loading:\n\n**Before:**\n```go\nfor _, id := range userIDs {\n  user := db.Query(\"SELECT * FROM users WHERE id = ?\", id)\n}\n```\n\n**After:**\n```go\nusers := db.Query(\"SELECT * FROM users WHERE id IN (?, ?, ?)\" userIDs...)\n```",

		"excessive_nesting": "Extract nested blocks into separate functions to improve readability",

		"magic_number": "Define a named constant instead of using a magic number",

		"missing_error_check": "Always check errors returned from functions",
	}
}
