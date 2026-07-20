package codereview

import (
	"fmt"
	"regexp"
)

// ApplyPolicies applies team policies to findings
func (p *PolicyEngine) ApplyPolicies(findings []Finding) []Finding {
	filtered := make([]Finding, 0)

	for _, finding := range findings {
		// Check if this finding matches any policy rules
		shouldInclude := true

		for _, rule := range p.rules {
			if rule.Enabled && p.matchesRule(finding, rule) {
				// Apply rule modifications
				finding.Severity = rule.Severity
				break
			}
		}

		if shouldInclude {
			filtered = append(filtered, finding)
		}
	}

	return filtered
}

// matchesRule checks if a finding matches a policy rule
func (p *PolicyEngine) matchesRule(finding Finding, rule *PolicyRule) bool {
	if rule.Pattern == "" {
		return false
	}

	// Try to match the pattern against the finding
	matched, err := regexp.MatchString(rule.Pattern, finding.Message)
	if err != nil {
		return false
	}

	return matched
}

// AddRule adds a new policy rule
func (p *PolicyEngine) AddRule(rule *PolicyRule) error {
	if rule.ID == "" {
		return fmt.Errorf("rule ID is required")
	}

	if rule.Pattern != "" {
		// Validate regex pattern
		if _, err := regexp.Compile(rule.Pattern); err != nil {
			return fmt.Errorf("invalid regex pattern: %v", err)
		}
	}

	p.rules[rule.ID] = rule
	return nil
}

// RemoveRule removes a policy rule
func (p *PolicyEngine) RemoveRule(id string) {
	delete(p.rules, id)
}

// GetRule retrieves a policy rule
func (p *PolicyEngine) GetRule(id string) (*PolicyRule, bool) {
	rule, exists := p.rules[id]
	return rule, exists
}

// ListRules returns all policy rules
func (p *PolicyEngine) ListRules() []*PolicyRule {
	rules := make([]*PolicyRule, 0, len(p.rules))
	for _, rule := range p.rules {
		rules = append(rules, rule)
	}
	return rules
}

// DefaultPolicies returns a set of default policy rules
func DefaultPolicies() map[string]*PolicyRule {
	return map[string]*PolicyRule{
		"security_sql_injection": {
			ID:          "security_sql_injection",
			Name:        "SQL Injection Detection",
			Description: "Detects potential SQL injection vulnerabilities",
			Enabled:     true,
			Severity:    SeverityCritical,
			Pattern:     `(?i)(sql|query).*[\+\*].*(?:user|input|request)`,
			Message:     "Potential SQL injection vulnerability detected",
			Remediation: "Use parameterized queries or prepared statements",
			Tags:        []string{"security", "sql"},
		},
		"security_hardcoded_credentials": {
			ID:          "security_hardcoded_credentials",
			Name:        "Hardcoded Credentials",
			Description: "Detects hardcoded credentials in code",
			Enabled:     true,
			Severity:    SeverityCritical,
			Pattern:     `(?i)(password|api_key|secret|token)\s*=\s*["\']`,
			Message:     "Hardcoded credentials detected",
			Remediation: "Move credentials to environment variables or secrets manager",
			Tags:        []string{"security", "credentials"},
		},
		"style_function_length": {
			ID:          "style_function_length",
			Name:        "Function Length",
			Description: "Warns about overly long functions",
			Enabled:     true,
			Severity:    SeverityMedium,
			Pattern:     `function.*(?:lines|length)\s*>\s*50`,
			Message:     "Function exceeds recommended length",
			Remediation: "Consider breaking the function into smaller units",
			Tags:        []string{"style", "maintainability"},
		},
		"performance_nested_loops": {
			ID:          "performance_nested_loops",
			Name:        "Nested Loops",
			Description: "Warns about deeply nested loops",
			Enabled:     true,
			Severity:    SeverityMedium,
			Pattern:     `for\s+.*for\s+.*for`,
			Message:     "Deeply nested loops detected",
			Remediation: "Consider refactoring to reduce loop nesting depth",
			Tags:        []string{"performance"},
		},
	}
}

// LoadPoliciesFromConfig loads policies from configuration
func (p *PolicyEngine) LoadPoliciesFromConfig(config map[string]interface{}) error {
	// This would parse policies from configuration format
	// For now, we'll load default policies
	defaults := DefaultPolicies()

	for _, rule := range defaults {
		p.rules[rule.ID] = rule
	}

	return nil
}

// SetTeamPreference sets a team preference that affects policy application
func (p *PolicyEngine) SetTeamPreference(key string, value interface{}) {
	p.teamData[key] = value
}

// GetTeamPreference retrieves a team preference
func (p *PolicyEngine) GetTeamPreference(key string) (interface{}, bool) {
	val, exists := p.teamData[key]
	return val, exists
}
