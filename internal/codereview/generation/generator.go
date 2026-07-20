package generation

import (
	"context"
	"fmt"
	"strings"

	"nextcode.io/fantasy"
)

// Generator creates tests and documentation
type Generator struct {
	model fantasy.LanguageModel
}

// NewGenerator creates a new generator
func NewGenerator(model fantasy.LanguageModel) *Generator {
	return &Generator{model: model}
}

// GenerateUnitTest generates unit tests for a function
func (g *Generator) GenerateUnitTest(ctx context.Context, language string, code string, functionName string) (string, error) {
	prompt := fmt.Sprintf(`Generate comprehensive unit tests for the following %s function:

\`\`\`%s
%s
\`\`\`

Requirements:
- Include happy path tests
- Include edge case tests
- Include error handling tests
- Use best practices for %s testing
- Include test fixtures if needed

Generate ONLY the test code:`, language, language, code, language)

	return g.callModel(ctx, prompt)
}

// GenerateDocstring generates documentation for a function
func (g *Generator) GenerateDocstring(ctx context.Context, language string, code string, functionName string) (string, error) {
	prompt := fmt.Sprintf(`Generate comprehensive documentation/docstring for the following %s function:

\`\`\`%s
%s
\`\`\`

Include:
- Function description
- Parameters/arguments with types
- Return value description
- Exceptions/errors it might raise
- Usage examples
- Related functions if any

Format for %s:`, language, language, code, language)

	return g.callModel(ctx, prompt)
}

// GenerateREADME generates a README for a file or module
func (g *Generator) GenerateREADME(ctx context.Context, files map[string]string, projectName string) (string, error) {
	fileList := ""
	for name, content := range files {
		fileList += fmt.Sprintf("\n%s:\n```\n%s\n```\n", name, truncate(content, 200))
	}

	prompt := fmt.Sprintf(`Generate a professional README for the following %s project:

Files:
%s

Include:
- Project description
- Features
- Installation instructions
- Usage examples
- API documentation
- Contributing guidelines
- License

Make it clear and beginner-friendly.`, projectName, fileList)

	return g.callModel(ctx, prompt)
}

// GenerateAPIDocumentation generates API documentation
func (g *Generator) GenerateAPIDocumentation(ctx context.Context, endpoints []map[string]interface{}) (string, error) {
	endpointDesc := ""
	for i, ep := range endpoints {
		endpointDesc += fmt.Sprintf("\n%d. %s %s", i+1, ep["method"], ep["path"])
	}

	prompt := fmt.Sprintf(`Generate comprehensive API documentation for these endpoints:

%s

For each endpoint include:
- Method and path
- Description
- Parameters
- Request body example
- Response example
- Error codes
- Authentication requirements

Format as Markdown with clear sections.`, endpointDesc)

	return g.callModel(ctx, prompt)
}

// GenerateTypeDefinitions generates TypeScript/Rust type definitions from examples
func (g *Generator) GenerateTypeDefinitions(ctx context.Context, language string, dataExamples []string) (string, error) {
	examples := strings.Join(dataExamples, "\n\n")

	prompt := fmt.Sprintf(`Generate %s type definitions/interfaces for the following data examples:

Examples:
%s

Create:
- Type definitions for all data structures
- Include proper typing for nested objects
- Use appropriate %s idioms and conventions
- Include JSDoc/comments where helpful

Generate ONLY the type definitions:`, language, examples, language)

	return g.callModel(ctx, prompt)
}

// GenerateErrorHandling generates error handling code
func (g *Generator) GenerateErrorHandling(ctx context.Context, language string, functionCode string) (string, error) {
	prompt := fmt.Sprintf(`Add comprehensive error handling to this %s function:

\`\`\`%s
%s
\`\`\`

Add:
- Input validation
- Error checks for all operations
- Proper error logging
- Error recovery where possible
- Error messages with context

Generate the improved code:`, language, language, functionCode)

	return g.callModel(ctx, prompt)
}

// callModel calls the LLM
func (g *Generator) callModel(ctx context.Context, prompt string) (string, error) {
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

	response, err := g.model.Complete(ctx, messages, &fantasy.CompleteOptions{
		MaxTokens:   3000,
		Temperature: 0.3, // Lower temperature for more consistent output
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

// truncate truncates string to max length
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
