package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sauravmarvani/nextcode/internal/codereview"
)

// Bot handles Slack integration for NextCode
type Bot struct {
	webhookURL string
	channel    string
	client     *http.Client
}

// NewBot creates a new Slack bot
func NewBot(webhookURL, channel string) *Bot {
	return &Bot{
		webhookURL: webhookURL,
		channel:    channel,
		client:     &http.Client{},
	}
}

// MessageBlock represents a Slack message block
type MessageBlock struct {
	Type string      `json:"type"`
	Text *TextBlock  `json:"text,omitempty"`
	Elements []interface{} `json:"elements,omitempty"`
}

// TextBlock represents Slack text formatting
type TextBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// SlackMessage represents a complete Slack message
type SlackMessage struct {
	Channel   string         `json:"channel,omitempty"`
	Text      string         `json:"text"`
	Blocks    []interface{}  `json:"blocks,omitempty"`
	ThreadTs  string         `json:"thread_ts,omitempty"`
}

// PostReviewStart announces the start of a review
func (b *Bot) PostReviewStart(ctx context.Context, prNumber int, prTitle string, author string) error {
	message := SlackMessage{
		Channel: b.channel,
		Text:    fmt.Sprintf("Starting code review of PR #%d", prNumber),
		Blocks: []interface{}{
			map[string]interface{}{
				"type": "header",
				"text": map[string]interface{}{
					"type": "plain_text",
					"text": fmt.Sprintf("🔍 Code Review: PR #%d", prNumber),
				},
			},
			map[string]interface{}{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*%s* by @%s", prTitle, author),
				},
			},
		},
	}

	return b.sendMessage(ctx, message)
}

// PostReviewComplete posts the final review summary
func (b *Bot) PostReviewComplete(ctx context.Context, result *codereview.ReviewResult) error {
	severity := "✅"
	if result.Statistics.CriticalCount > 0 {
		severity = "🔴"
	} else if result.Statistics.HighCount > 0 {
		severity = "🟠"
	}

	text := fmt.Sprintf("%s Review Complete - PR #%d: %d findings", severity, result.PRNumber, result.Statistics.FindingCount)

	blocks := []interface{}{
		map[string]interface{}{
			"type": "header",
			"text": map[string]interface{}{
				"type": "plain_text",
				"text": text,
			},
		},
		map[string]interface{}{
			"type": "section",
			"fields": []interface{}{
				map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*Critical:* %d", result.Statistics.CriticalCount),
				},
				map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*High:* %d", result.Statistics.HighCount),
				},
				map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*Medium:* %d", result.Statistics.MediumCount),
				},
				map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*Low:* %d", result.Statistics.LowCount),
				},
			},
		},
	}

	// Add top issues
	if len(result.Findings) > 0 {
		issueText := "*Top Issues:*\n"
		for i, finding := range result.Findings {
			if i >= 3 {
				issueText += fmt.Sprintf("... and %d more", len(result.Findings)-3)
				break
			}
			issueText += fmt.Sprintf("• %s (%s)\n", finding.Rule, finding.Severity)
		}

		blocks = append(blocks, map[string]interface{}{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": issueText,
			},
		})
	}

	message := SlackMessage{
		Channel: b.channel,
		Text:    text,
		Blocks:  blocks,
	}

	return b.sendMessage(ctx, message)
}

// PostFindingThread posts a finding in a thread
func (b *Bot) PostFindingThread(ctx context.Context, finding codereview.Finding, threadTs string) error {
	emoji := "🔴"
	switch finding.Severity {
	case codereview.SeverityHigh:
		emoji = "🟠"
	case codereview.SeverityMedium:
		emoji = "🟡"
	case codereview.SeverityLow:
		emoji = "🔵"
	}

	text := fmt.Sprintf("%s %s: %s", emoji, finding.Rule, finding.Message)

	message := SlackMessage{
		Channel:  b.channel,
		Text:     text,
		ThreadTs: threadTs,
		Blocks: []interface{}{
			map[string]interface{}{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*%s*: %s\n_File:_ %s:%d", finding.Title, finding.Message, finding.File, finding.Line),
				},
			},
		},
	}

	return b.sendMessage(ctx, message)
}

// PostSlashCommand handles slash commands
func (b *Bot) HandleSlashCommand(ctx context.Context, command string) (string, error) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid command")
	}

	cmd := parts[0]
	switch cmd {
	case "/review":
		return "Starting code review...", nil

	case "/fix":
		return "Applying auto-fixes...", nil

	case "/status":
		return "NextCode is running and ready to review", nil

	default:
		return "Unknown command. Try: /review, /fix, /status", nil
	}
}

// sendMessage sends a message to Slack
func (b *Bot) sendMessage(ctx context.Context, message SlackMessage) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", b.webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("slack returned status %d: %s", resp.StatusCode, body)
	}

	return nil
}

// NotifyTeam sends a team notification
func (b *Bot) NotifyTeam(ctx context.Context, title string, message string) error {
	slackMsg := SlackMessage{
		Channel: b.channel,
		Text:    title,
		Blocks: []interface{}{
			map[string]interface{}{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*%s*\n%s", title, message),
				},
			},
		},
	}

	return b.sendMessage(ctx, slackMsg)
}
