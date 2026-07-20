package summary

import (
	"context"
	"fmt"
	"strings"

	"github.com/sauravmarvani/nextcode/internal/vcs"
	"nextcode.io/fantasy"
)

// Summarizer generates human-readable summaries of code changes
type Summarizer struct {
	model fantasy.LanguageModel
}

// NewSummarizer creates a new diff summarizer
func NewSummarizer(model fantasy.LanguageModel) *Summarizer {
	return &Summarizer{
		model: model,
	}
}

// SummaryResult contains the generated summary
type SummaryResult struct {
	Title       string            // One-line summary
	Description string            // Detailed description
	Impact      string            // What changed and why
	Risks       []string          // Potential risks
	Highlights  []string          // Key highlights
	Files       map[string]string // File-level summaries
	Diagram     string            // ASCII architecture diagram
	Walkthrough string            // Step-by-step walkthrough
}

// Summarize generates a TL;DR summary of the PR changes
func (s *Summarizer) Summarize(ctx context.Context, pr *vcs.PullRequest, files []vcs.FileDiff, commits []vcs.CommitInfo) (*SummaryResult, error) {
	result := &SummaryResult{
		Files: make(map[string]string),
	}

	// Generate title and description
	titleAndDesc, err := s.generateTitleAndDescription(ctx, pr, files, commits)
	if err != nil {
		return nil, err
	}
	result.Title = titleAndDesc.Title
	result.Description = titleAndDesc.Description

	// Generate impact analysis
	impact, err := s.generateImpactAnalysis(ctx, files)
	if err != nil {
		return nil, err
	}
	result.Impact = impact

	// Generate risk assessment
	risks, err := s.generateRisks(ctx, files)
	if err != nil {
		return nil, err
	}
	result.Risks = risks

	// Generate highlights
	highlights, err := s.generateHighlights(ctx, files)
	if err != nil {
		return nil, err
	}
	result.Highlights = highlights

	// Generate file-level summaries
	for _, file := range files {
		fileSummary, err := s.summarizeFile(ctx, file)
		if err == nil {
			result.Files[file.Path] = fileSummary
		}
	}

	// Generate architecture diagram
	if hasArchitectureChanges(files) {
		diagram, err := s.generateArchitectureDiagram(ctx, files)
		if err == nil {
			result.Diagram = diagram
		}
	}

	// Generate walkthrough
	walkthrough, err := s.generateWalkthrough(ctx, pr, files)
	if err == nil {
		result.Walkthrough = walkthrough
	}

	return result, nil
}

// titleDescResult holds title and description
type titleDescResult struct {
	Title       string
	Description string
}

// generateTitleAndDescription generates a title and description for the PR
func (s *Summarizer) generateTitleAndDescription(ctx context.Context, pr *vcs.PullRequest, files []vcs.FileDiff, commits []vcs.CommitInfo) (*titleDescResult, error) {
	prompt := fmt.Sprintf(`Analyze this PR and generate a concise summary.

PR Title: %s
PR Description: %s

Files Changed: %d
Total Changes: %d lines added, %d lines removed

File List:
%s

Commit Messages:
%s

Generate:
1. A one-line title (max 60 chars) summarizing the main change
2. A 2-3 line description explaining what this PR does

Format:
Title: [your title]
Description: [your description]`,
		pr.Title,
		pr.Description,
		len(files),
		countAddedLines(files),
		countRemovedLines(files),
		formatFileList(files),
		formatCommits(commits))

	response, err := s.callModel(ctx, prompt)
	if err != nil {
		return nil, err
	}

	return parseResponse(response), nil
}

// generateImpactAnalysis generates an impact analysis
func (s *Summarizer) generateImpactAnalysis(ctx context.Context, files []vcs.FileDiff) (string, error) {
	prompt := fmt.Sprintf(`Analyze the impact of these code changes in 2-3 sentences.

Files Changed: %d
Added Lines: %d
Removed Lines: %d

Files:
%s

Focus on: What functionality changed, what users will notice, what systems are affected.`,
		len(files),
		countAddedLines(files),
		countRemovedLines(files),
		formatFileList(files))

	return s.callModel(ctx, prompt)
}

// generateRisks generates a risk assessment
func (s *Summarizer) generateRisks(ctx context.Context, files []vcs.FileDiff) ([]string, error) {
	prompt := fmt.Sprintf(`Identify 3-5 potential risks or concerns with these changes:

Files:
%s

For each risk, provide a brief statement (max 1 line each).`,
		formatFileList(files))

	response, err := s.callModel(ctx, prompt)
	if err != nil {
		return nil, err
	}

	risks := strings.Split(response, "\n")
	var cleanRisks []string
	for _, risk := range risks {
		risk = strings.TrimSpace(risk)
		if risk != "" && !strings.HasPrefix(risk, "-") {
			cleanRisks = append(cleanRisks, risk)
		}
	}

	return cleanRisks, nil
}

// generateHighlights extracts key highlights
func (s *Summarizer) generateHighlights(ctx context.Context, files []vcs.FileDiff) ([]string, error) {
	prompt := fmt.Sprintf(`List 3-5 key highlights or improvements from these changes:

Files:
%s

Format as bullet points.`,
		formatFileList(files))

	response, err := s.callModel(ctx, prompt)
	if err != nil {
		return nil, err
	}

	highlights := strings.Split(response, "\n")
	var cleanHighlights []string
	for _, h := range highlights {
		h = strings.TrimSpace(h)
		if h != "" {
			cleanHighlights = append(cleanHighlights, h)
		}
	}

	return cleanHighlights, nil
}

// summarizeFile generates a summary for a single file
func (s *Summarizer) summarizeFile(ctx context.Context, file vcs.FileDiff) (string, error) {
	prompt := fmt.Sprintf(`Summarize this file change in 1 sentence:

File: %s
Status: %s
Added: %d lines
Removed: %d lines`,
		file.Path,
		file.Status,
		countFileAddedLines(file),
		countFileRemovedLines(file))

	return s.callModel(ctx, prompt)
}

// generateArchitectureDiagram generates an ASCII architecture diagram
func (s *Summarizer) generateArchitectureDiagram(ctx context.Context, files []vcs.FileDiff) (string, error) {
	prompt := fmt.Sprintf(`Generate a simple ASCII architecture diagram showing how these changed components interact:

Changed Files:
%s

Use simple ASCII boxes and arrows. Max 20 lines.`,
		formatFileList(files))

	return s.callModel(ctx, prompt)
}

// generateWalkthrough generates a step-by-step walkthrough
func (s *Summarizer) generateWalkthrough(ctx context.Context, pr *vcs.PullRequest, files []vcs.FileDiff) (string, error) {
	prompt := fmt.Sprintf(`Create a step-by-step walkthrough for understanding this PR:

Title: %s
Files Changed: %d

Files:
%s

Format as numbered steps (3-5 steps max).`,
		pr.Title,
		len(files),
		formatFileList(files))

	return s.callModel(ctx, prompt)
}

// callModel calls the LLM model
func (s *Summarizer) callModel(ctx context.Context, prompt string) (string, error) {
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

	response, err := s.model.Complete(ctx, messages, &fantasy.CompleteOptions{
		MaxTokens:   1000,
		Temperature: 0.5,
	})
	if err != nil {
		return "", err
	}

	if len(response.Message.Content) == 0 {
		return "", fmt.Errorf("empty response from model")
	}

	textPart, ok := fantasy.AsMessagePart[fantasy.TextPart](response.Message.Content[0])
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	return textPart.Text, nil
}

// Helper functions

func hasArchitectureChanges(files []vcs.FileDiff) bool {
	for _, file := range files {
		if strings.Contains(file.Path, "arch") || strings.Contains(file.Path, "structure") {
			return true
		}
	}
	return false
}

func countAddedLines(files []vcs.FileDiff) int {
	count := 0
	for _, f := range files {
		count += countFileAddedLines(f)
	}
	return count
}

func countRemovedLines(files []vcs.FileDiff) int {
	count := 0
	for _, f := range files {
		count += countFileRemovedLines(f)
	}
	return count
}

func countFileAddedLines(file vcs.FileDiff) int {
	return strings.Count(file.Patch, "\n+") - strings.Count(file.Patch, "\n+++")
}

func countFileRemovedLines(file vcs.FileDiff) int {
	return strings.Count(file.Patch, "\n-") - strings.Count(file.Patch, "\n---")
}

func formatFileList(files []vcs.FileDiff) string {
	var list []string
	for _, f := range files {
		list = append(list, fmt.Sprintf("- %s (%s)", f.Path, f.Status))
	}
	return strings.Join(list, "\n")
}

func formatCommits(commits []vcs.CommitInfo) string {
	var list []string
	for _, c := range commits {
		list = append(list, fmt.Sprintf("- %s", c.Message))
	}
	return strings.Join(list, "\n")
}

func parseResponse(response string) *titleDescResult {
	result := &titleDescResult{}
	lines := strings.Split(response, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "Title:") {
			result.Title = strings.TrimPrefix(line, "Title:")
			result.Title = strings.TrimSpace(result.Title)
		} else if strings.HasPrefix(line, "Description:") {
			result.Description = strings.TrimPrefix(line, "Description:")
			result.Description = strings.TrimSpace(result.Description)
		}
	}

	return result
}
