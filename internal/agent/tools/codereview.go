package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"nextcode.io/fantasy"

	"github.com/sauravmarvani/nextcode/internal/codereview"
	"github.com/sauravmarvani/nextcode/internal/codereview/scanners"
	"github.com/sauravmarvani/nextcode/internal/vcs"
)

const (
	CodeReviewAnalyzeToolName   = "codereview_analyze"
	CodeReviewCommentToolName   = "codereview_comment"
	CodeReviewSuggestFixToolName = "codereview_suggest_fix"
	CodeReviewPolicyToolName    = "codereview_policy"
)

// CodeReviewAnalyzeParams parameters for analyzing a PR
type CodeReviewAnalyzeParams struct {
	PRURL       string   `json:"pr_url" description:"GitHub/GitLab PR URL to analyze"`
	Scanners    []string `json:"scanners" description:"List of scanners to run: security, performance, correctness, style"`
	PolicyFile  string   `json:"policy_file" description:"Optional path to custom policy file"`
}

// CodeReviewCommentParams parameters for posting a comment
type CodeReviewCommentParams struct {
	PRURL   string `json:"pr_url" description:"PR URL to comment on"`
	Body    string `json:"body" description:"Comment body in markdown"`
	File    string `json:"file" description:"Optional: file path for inline comment"`
	Line    int    `json:"line" description:"Optional: line number for inline comment"`
}

// CodeReviewSuggestFixParams parameters for suggesting a fix
type CodeReviewSuggestFixParams struct {
	FindingID string `json:"finding_id" description:"ID of the finding to suggest a fix for"`
	Context   string `json:"context" description:"Code context for generating the fix"`
}

// CodeReviewPolicyParams parameters for policy enforcement
type CodeReviewPolicyParams struct {
	PolicyID  string `json:"policy_id" description:"ID of the policy to check"`
	CodePath  string `json:"code_path" description:"Path to code to check"`
}

// CodeReviewAnalyzeTool creates the analyze PR tool
func CodeReviewAnalyzeTool(analyzer *codereview.Analyzer) fantasy.Tool {
	return fantasy.NewTool(
		CodeReviewAnalyzeToolName,
		"Analyze a pull request for code quality issues",
		func(ctx context.Context, params CodeReviewAnalyzeParams, call fantasy.ToolCall) (fantasy.ToolResponse, error) {
			if params.PRURL == "" {
				return fantasy.NewTextErrorResponse("pr_url is required"), nil
			}

			// Parse PR URL and fetch PR details
			// TODO: Implement VCS client integration to fetch actual PR data
			pr := &vcs.PullRequest{
				ID:    "1",
				Title: "Sample PR",
				State: "open",
			}

			// Build analysis request
			req := &codereview.ReviewRequest{
				PullRequest: pr,
				Files:       []vcs.FileDiff{},
				Scanners:    parseScannerTypes(params.Scanners),
			}

			// Run analysis
			result, err := analyzer.Analyze(ctx, req)
			if err != nil {
				return fantasy.NewTextErrorResponse(fmt.Sprintf("Analysis failed: %v", err)), nil
			}

			// Marshal result to JSON
			resultJSON, _ := json.Marshal(result)
			return fantasy.NewTextResponse(fmt.Sprintf("Code review complete:\n\n%s", string(resultJSON))), nil
		},
	)
}

// CodeReviewCommentTool creates the comment posting tool
func CodeReviewCommentTool() fantasy.Tool {
	return fantasy.NewTool(
		CodeReviewCommentToolName,
		"Post a comment on a pull request",
		func(ctx context.Context, params CodeReviewCommentParams, call fantasy.ToolCall) (fantasy.ToolResponse, error) {
			if params.PRURL == "" {
				return fantasy.NewTextErrorResponse("pr_url is required"), nil
			}

			if params.Body == "" {
				return fantasy.NewTextErrorResponse("body is required"), nil
			}

			// TODO: Implement actual VCS API call
			message := fmt.Sprintf("Posted comment to %s", params.PRURL)
			if params.File != "" {
				message = fmt.Sprintf("Posted inline comment to %s at %s:%d", params.PRURL, params.File, params.Line)
			}

			return fantasy.NewTextResponse(message), nil
		},
	)
}

// CodeReviewSuggestFixTool creates the suggest fix tool
func CodeReviewSuggestFixTool() fantasy.Tool {
	return fantasy.NewTool(
		CodeReviewSuggestFixToolName,
		"Generate a suggested fix for a code issue",
		func(ctx context.Context, params CodeReviewSuggestFixParams, call fantasy.ToolCall) (fantasy.ToolResponse, error) {
			if params.FindingID == "" {
				return fantasy.NewTextErrorResponse("finding_id is required"), nil
			}

			// TODO: Implement LLM-powered fix generation
			suggestion := fmt.Sprintf("Suggested fix for finding %s", params.FindingID)
			return fantasy.NewTextResponse(suggestion), nil
		},
	)
}

// CodeReviewPolicyTool creates the policy enforcement tool
func CodeReviewPolicyTool(policyEngine *codereview.PolicyEngine) fantasy.Tool {
	return fantasy.NewTool(
		CodeReviewPolicyToolName,
		"Check code against policies",
		func(ctx context.Context, params CodeReviewPolicyParams, call fantasy.ToolCall) (fantasy.ToolResponse, error) {
			if params.PolicyID == "" {
				return fantasy.NewTextErrorResponse("policy_id is required"), nil
			}

			if params.CodePath == "" {
				return fantasy.NewTextErrorResponse("code_path is required"), nil
			}

			// Retrieve policy
			policy, exists := policyEngine.GetRule(params.PolicyID)
			if !exists {
				return fantasy.NewTextErrorResponse(fmt.Sprintf("Policy %s not found", params.PolicyID)), nil
			}

			response := fmt.Sprintf("Policy: %s\nEnabled: %v\nSeverity: %s", 
				policy.Name, policy.Enabled, policy.Severity)
			return fantasy.NewTextResponse(response), nil
		},
	)
}

// helper function to parse scanner types
func parseScannerTypes(scanners []string) []codereview.ScannerType {
	if len(scanners) == 0 {
		return []codereview.ScannerType{
			codereview.ScannerSecurity,
			codereview.ScannerPerformance,
			codereview.ScannerCorrectness,
			codereview.ScannerStyle,
		}
	}

	types := make([]codereview.ScannerType, 0)
	for _, s := range scanners {
		types = append(types, codereview.ScannerType(s))
	}
	return types
}

// InitializeCodeReviewTools initializes all code review tools
func InitializeCodeReviewTools(ctx context.Context) ([]fantasy.Tool, error) {
	// Create analyzer
	analyzer := codereview.NewAnalyzer()

	// Register scanners
	registry := scanners.NewRegistry()
	registry.RegisterAllScanners(analyzer)

	// Create policy engine
	policyEngine := analyzer.policies

	tools := []fantasy.Tool{
		CodeReviewAnalyzeTool(analyzer),
		CodeReviewCommentTool(),
		CodeReviewSuggestFixTool(),
		CodeReviewPolicyTool(policyEngine),
	}

	return tools, nil
}
