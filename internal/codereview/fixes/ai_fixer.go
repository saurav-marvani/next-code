package fixes

import (
	"context"
	"fmt"
	"os"
	"strings"

	"nextcode.io/fantasy"
)

// AIFixer applies AI-powered fixes using Fantasy SDK
type AIFixer struct {
	model   fantasy.LanguageModel
	maxTokens int
}

// NewAIFixer creates a new AI fixer
func NewAIFixer(model fantasy.LanguageModel) *AIFixer {
	return &AIFixer{
		model:     model,
		maxTokens: 2000,
	}
}

// Fix applies an AI-powered fix
func (af *AIFixer) Fix(ctx context.Context, req FixRequest) (*FixResult, error) {
	result := &FixResult{
		FindingID:   req.Finding.ID,
		Strategy:    FixStrategyAI,
		AIGenerated: true,
	}

	// Read the file
	content, err := os.ReadFile(req.Finding.File)
	if err != nil {
		result.Error = fmt.Errorf("failed to read file: %w", err)
		return result, result.Error
	}

	result.Original = string(content)

	// Build the prompt for the AI model
	prompt := af.buildPrompt(req, string(content))

	// Call Fantasy SDK
	messages := []fantasy.Message{
		{
			Role: fantasy.MessageRoleUser,
			Content: []fantasy.MessagePart{
				fantasy.TextPart{
					Text: prompt,
				},
			},
		},
	}

	response, err := af.model.Complete(ctx, messages, &fantasy.CompleteOptions{
		MaxTokens:   af.maxTokens,
		Temperature: req.Temperature,
	})
	if err != nil {
		result.Error = fmt.Errorf("AI model error: %w", err)
		return result, result.Error
	}

	if len(response.Message.Content) == 0 {
		result.Error = fmt.Errorf("empty response from AI model")
		return result, result.Error
	}

	// Extract the fixed code from the response
	textPart, ok := fantasy.AsMessagePart[fantasy.TextPart](response.Message.Content[0])
	if !ok {
		result.Error = fmt.Errorf("unexpected response format from AI model")
		return result, result.Error
	}

	fixedCode := af.extractCode(textPart.Text)
	if fixedCode == "" {
		result.Error = fmt.Errorf("AI model did not provide valid code fix")
		return result, result.Error
	}

	result.Fixed = fixedCode
	result.Explanation = af.extractExplanation(textPart.Text)
	result.Success = true
	result.Confidence = 0.85 // AI-generated fixes are moderately confident
	result.CommitMessage = fmt.Sprintf("fix: apply AI suggestion for %s", req.Finding.Rule)

	// Write the fixed content
	if err := os.WriteFile(req.Finding.File, []byte(fixedCode), 0644); err != nil {
		result.Success = false
		result.Error = fmt.Errorf("failed to write fixed file: %w", err)
		return result, result.Error
	}

	return result, nil
}

// buildPrompt creates a prompt for the AI model
func (af *AIFixer) buildPrompt(req FixRequest, content string) string {
	prompt := fmt.Sprintf(`You are an expert code reviewer and fixer. Your task is to fix the following code issue.

File: %s
Language: %s
Issue: %s
Rule: %s
Severity: %s

Current Code:
\`\`\`%s
%s
\`\`\`

Description of the issue:
%s

Please provide:
1. The corrected code (wrapped in code blocks with the language specified)
2. A brief explanation of what was wrong and how you fixed it

Provide ONLY the fixed code and explanation, nothing else.`, 
	req.Finding.File,
	req.Language,
	req.Finding.Title,
	req.Finding.Rule,
	req.Finding.Severity,
	req.Language,
	req.Finding.Code,
	req.Finding.Description)

	if req.Context != "" {
		prompt += fmt.Sprintf("\n\nAdditional context:\n%s", req.Context)
	}

	return prompt
}

// extractCode extracts the code block from the AI response
func (af *AIFixer) extractCode(response string) string {
	// Look for code blocks
	start := strings.Index(response, "```")
	if start == -1 {
		return ""
	}

	// Skip the language specifier
	start = strings.Index(response[start:], "\n")
	if start == -1 {
		return ""
	}
	start += start

	end := strings.Index(response[start:], "```")
	if end == -1 {
		return ""
	}
	end += start

	return strings.TrimSpace(response[start:end])
}

// extractExplanation extracts the explanation from the AI response
func (af *AIFixer) extractExplanation(response string) string {
	// Find content after the code block
	endCodeBlock := strings.LastIndex(response, "```")
	if endCodeBlock == -1 {
		return ""
	}

	explanation := strings.TrimSpace(response[endCodeBlock+3:])
	return explanation
}

// SetMaxTokens sets the maximum tokens for AI responses
func (af *AIFixer) SetMaxTokens(tokens int) {
	af.maxTokens = tokens
}

// Supports checks if AI fixer can handle a specific finding type
func (af *AIFixer) Supports(finding codereview.Finding) bool {
	// AI can handle most types except critical security issues (require manual review)
	if finding.Severity == codereview.SeverityCritical && finding.Type == "security" {
		return false
	}
	return true
}
