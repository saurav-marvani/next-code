package policies

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// PolicyParser parses Markdown policy files
type PolicyParser struct {
	content string
}

// NewPolicyParser creates a new policy parser
func NewPolicyParser(content string) *PolicyParser {
	return &PolicyParser{
		content: content,
	}
}

// ParseMarkdownPolicy parses a Markdown policy file
func (p *PolicyParser) ParseMarkdownPolicy() ([]*codereview.PolicyRule, error) {
	rules := make([]*codereview.PolicyRule, 0)

	// Split by rule sections (## Rule Name)
	sections := strings.Split(p.content, "\n## ")

	for _, section := range sections {
		if strings.TrimSpace(section) == "" {
			continue
		}

		rule, err := p.parseRuleSection(section)
		if err != nil {
			continue // Skip malformed rules
		}

		if rule != nil {
			rules = append(rules, rule)
		}
	}

	return rules, nil
}

// parseRuleSection parses a single rule section
func (p *PolicyParser) parseRuleSection(section string) (*codereview.PolicyRule, error) {
	lines := strings.Split(section, "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty section")
	}

	rule := &codereview.PolicyRule{
		Enabled: true,
		Tags:    make([]string, 0),
	}

	// First line is the rule name
	rule.Name = strings.TrimSpace(lines[0])
	if rule.Name == "" {
		return nil, fmt.Errorf("missing rule name")
	}

	// Generate ID from name
	rule.ID = strings.ToLower(strings.ReplaceAll(rule.Name, " ", "_"))

	// Parse remaining lines
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// Parse metadata
		if strings.HasPrefix(line, "**Severity:**") {
			severity := strings.TrimSpace(strings.TrimPrefix(line, "**Severity:**"))
			rule.Severity = parseSeverity(severity)
		}

		if strings.HasPrefix(line, "**Pattern:**") {
			pattern := strings.TrimSpace(strings.TrimPrefix(line, "**Pattern:**"))
			// Remove markdown code formatting
			pattern = strings.Trim(pattern, "`")
			rule.Pattern = pattern
		}

		if strings.HasPrefix(line, "**Description:**") {
			desc := strings.TrimSpace(strings.TrimPrefix(line, "**Description:**"))
			rule.Description = desc
		}

		if strings.HasPrefix(line, "**Remediation:**") {
			rem := strings.TrimSpace(strings.TrimPrefix(line, "**Remediation:**"))
			rule.Remediation = rem
		}

		if strings.HasPrefix(line, "**Message:**") {
			msg := strings.TrimSpace(strings.TrimPrefix(line, "**Message:**"))
			rule.Message = msg
		}

		if strings.HasPrefix(line, "**Enabled:**") {
			enabled := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(line, "**Enabled:**")))
			rule.Enabled = enabled == "true" || enabled == "yes"
		}

		if strings.HasPrefix(line, "**Tags:**") {
			tags := strings.TrimSpace(strings.TrimPrefix(line, "**Tags:**"))
			rule.Tags = strings.Split(tags, ",")
			for i, tag := range rule.Tags {
				rule.Tags[i] = strings.TrimSpace(tag)
			}
		}
	}

	// Validate pattern if provided
	if rule.Pattern != "" {
		if _, err := regexp.Compile(rule.Pattern); err != nil {
			return nil, fmt.Errorf("invalid regex pattern: %v", err)
		}
	}

	return rule, nil
}

// parseSeverity converts string to Severity
func parseSeverity(severity string) codereview.Severity {
	switch strings.ToLower(severity) {
	case "critical":
		return codereview.SeverityCritical
	case "high":
		return codereview.SeverityHigh
	case "medium":
		return codereview.SeverityMedium
	case "low":
		return codereview.SeverityLow
	case "info":
		return codereview.SeverityInfo
	default:
		return codereview.SeverityMedium
	}
}

// BuildDefaultPolicies returns default Markdown policies
func BuildDefaultPolicies() string {
	return `# NextCode Review Policies

## SQL Injection Prevention
**Severity:** critical
**Pattern:** \`(?i)(\+|concat|format)\s*.*(?:query|sql|execute)\`
**Description:** Detects potential SQL injection vulnerabilities
**Message:** Potential SQL injection: use parameterized queries
**Remediation:** Always use prepared statements or parameterized queries
**Enabled:** true
**Tags:** security, database

## Hardcoded Credentials
**Severity:** critical
**Pattern:** \`(?i)(password|api_?key|secret|token)\s*=\s*["'][^"']*["']\`
**Description:** Detects hardcoded credentials in source code
**Message:** Hardcoded credentials detected in code
**Remediation:** Move credentials to environment variables or secrets manager
**Enabled:** true
**Tags:** security, credentials

## XSS Prevention
**Severity:** high
**Pattern:** \`(?i)innerHTML\s*=|dangerouslySetInnerHTML\`
**Description:** Detects potential XSS vulnerabilities
**Message:** Potential XSS vulnerability: sanitize user input
**Remediation:** Use safe APIs or sanitize all user input before rendering
**Enabled:** true
**Tags:** security, xss

## N+1 Query Problem
**Severity:** high
**Pattern:** \`(?i)for\s+.*\{[^}]*(?:query|find|select)\`
**Description:** Detects N+1 query patterns in code
**Message:** Potential N+1 query problem
**Remediation:** Use batch queries or eager loading
**Enabled:** true
**Tags:** performance, database

## Function Length
**Severity:** medium
**Pattern:** \`(?i)function.*(?:lines|length)\s*>\s*50\`
**Description:** Warns about overly long functions
**Message:** Function exceeds 50 lines - consider breaking it down
**Remediation:** Extract logic into smaller functions
**Enabled:** true
**Tags:** style, maintainability

## Error Handling
**Severity:** high
**Pattern:** \`(?i)(?:err|error)\s*:=\`
**Description:** Detects missing error handling
**Message:** Error handling required
**Remediation:** Always check and handle errors from function calls
**Enabled:** true
**Tags:** correctness, reliability
`
}
