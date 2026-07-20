package policies

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// PolicyLoader loads policies from files and directories
type PolicyLoader struct {
	searchPaths []string
}

// NewPolicyLoader creates a new policy loader
func NewPolicyLoader(searchPaths ...string) *PolicyLoader {
	return &PolicyLoader{
		searchPaths: searchPaths,
	}
}

// LoadFromFile loads policies from a markdown file
func (l *PolicyLoader) LoadFromFile(filePath string) ([]*codereview.PolicyRule, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read policy file: %v", err)
	}

	parser := NewPolicyParser(string(content))
	return parser.ParseMarkdownPolicy()
}

// LoadFromDirectory loads all policy files from a directory
func (l *PolicyLoader) LoadFromDirectory(dirPath string) ([]*codereview.PolicyRule, error) {
	rules := make([]*codereview.PolicyRule, 0)

	// Check if directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return rules, nil
	}

	// Read all .md files
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() || !filepath.Ext(file.Name()) == ".md" {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		fileRules, err := l.LoadFromFile(filePath)
		if err != nil {
			continue // Skip files with errors
		}

		rules = append(rules, fileRules...)
	}

	return rules, nil
}

// SearchAndLoad searches for policy files in known locations
func (l *PolicyLoader) SearchAndLoad(projectRoot string) ([]*codereview.PolicyRule, error) {
	rules := make([]*codereview.PolicyRule, 0)

	// Search paths
	searchPaths := []string{
		filepath.Join(projectRoot, ".nextcode-policy.md"),
		filepath.Join(projectRoot, ".code-review", "policies.md"),
		filepath.Join(projectRoot, "docs", "code-review-policies.md"),
	}

	// Add custom search paths
	searchPaths = append(searchPaths, l.searchPaths...)

	for _, path := range searchPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		fileRules, err := l.LoadFromFile(path)
		if err != nil {
			continue
		}

		rules = append(rules, fileRules...)
	}

	// If no policies found, use defaults
	if len(rules) == 0 {
		parser := NewPolicyParser(BuildDefaultPolicies())
		defaultRules, _ := parser.ParseMarkdownPolicy()
		return defaultRules, nil
	}

	return rules, nil
}

// LoadDefault loads the default policies
func (l *PolicyLoader) LoadDefault() []*codereview.PolicyRule {
	parser := NewPolicyParser(BuildDefaultPolicies())
	rules, _ := parser.ParseMarkdownPolicy()
	return rules
}

// SavePolicyFile saves rules to a markdown policy file
func (l *PolicyLoader) SavePolicyFile(rules []*codereview.PolicyRule, filePath string) error {
	content := buildPolicyMarkdown(rules)
	return ioutil.WriteFile(filePath, []byte(content), 0644)
}

// buildPolicyMarkdown builds markdown content from rules
func buildPolicyMarkdown(rules []*codereview.PolicyRule) string {
	content := "# NextCode Code Review Policies\n\n"

	for _, rule := range rules {
		content += fmt.Sprintf("## %s\n", rule.Name)

		if rule.Description != "" {
			content += fmt.Sprintf("**Description:** %s\n", rule.Description)
		}

		content += fmt.Sprintf("**Severity:** %s\n", rule.Severity)

		if rule.Pattern != "" {
			content += fmt.Sprintf("**Pattern:** `%s`\n", rule.Pattern)
		}

		if rule.Message != "" {
			content += fmt.Sprintf("**Message:** %s\n", rule.Message)
		}

		if rule.Remediation != "" {
			content += fmt.Sprintf("**Remediation:** %s\n", rule.Remediation)
		}

		content += fmt.Sprintf("**Enabled:** %v\n", rule.Enabled)

		if len(rule.Tags) > 0 {
			tags := ""
			for i, tag := range rule.Tags {
				if i > 0 {
					tags += ", "
				}
				tags += tag
			}
			content += fmt.Sprintf("**Tags:** %s\n", tags)
		}

		content += "\n"
	}

	return content
}

// MergePolicies merges multiple policy lists, with later ones taking precedence
func MergePolicies(policyLists ...[]*codereview.PolicyRule) []*codereview.PolicyRule {
	merged := make(map[string]*codereview.PolicyRule)

	for _, policies := range policyLists {
		for _, policy := range policies {
			merged[policy.ID] = policy
		}
	}

	result := make([]*codereview.PolicyRule, 0, len(merged))
	for _, policy := range merged {
		result = append(result, policy)
	}

	return result
}
