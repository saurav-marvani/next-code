package integrations

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/sauravmarvani/nextcode/internal/vcs"
)

// GitLabClient implements the VCS client for GitLab
type GitLabClient struct {
	token   string
	baseURL string
}

// NewGitLabClient creates a new GitLab client
func NewGitLabClient(token string, baseURL string) *GitLabClient {
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}
	return &GitLabClient{
		token:   token,
		baseURL: baseURL,
	}
}

// Authenticate verifies the token
func (c *GitLabClient) Authenticate(ctx context.Context, token string) error {
	if token == "" {
		return vcs.ErrAuthentication
	}
	c.token = token
	// TODO: Verify token by making an API call
	return nil
}

// IsAuthenticated checks if authenticated
func (c *GitLabClient) IsAuthenticated(ctx context.Context) (bool, error) {
	return c.token != "", nil
}

// GetRepositoryInfo retrieves repository information
func (c *GitLabClient) GetRepositoryInfo(ctx context.Context, owner, repo string) (*vcs.RepositoryInfo, error) {
	// TODO: Implement using go-gitlab
	return &vcs.RepositoryInfo{
		Owner:         owner,
		Name:          repo,
		URL:           fmt.Sprintf("%s/%s/%s", c.baseURL, owner, repo),
		DefaultBranch: "main",
		Private:       false,
	}, nil
}

// GetPullRequest retrieves an MR (Merge Request in GitLab)
func (c *GitLabClient) GetPullRequest(ctx context.Context, owner, repo string, number int) (*vcs.PullRequest, error) {
	// TODO: Implement using go-gitlab
	return &vcs.PullRequest{
		ID:     fmt.Sprintf("%d", number),
		Number: number,
		Title:  "Sample MR",
		State:  "opened",
		Source: "feature-branch",
		Target: "main",
	}, nil
}

// ListPullRequestFiles retrieves files in an MR
func (c *GitLabClient) ListPullRequestFiles(ctx context.Context, owner, repo string, number int) ([]vcs.FileDiff, error) {
	// TODO: Implement using go-gitlab
	return []vcs.FileDiff{}, nil
}

// GetPullRequestDiff retrieves the full diff of an MR
func (c *GitLabClient) GetPullRequestDiff(ctx context.Context, owner, repo string, number int) (io.ReadCloser, error) {
	// TODO: Implement using go-gitlab
	return nil, vcs.ErrNotImplemented
}

// GetCommit retrieves a commit
func (c *GitLabClient) GetCommit(ctx context.Context, owner, repo, sha string) (*vcs.CommitInfo, error) {
	// TODO: Implement using go-gitlab
	return &vcs.CommitInfo{
		SHA: sha,
	}, nil
}

// ListCommits retrieves commits from a branch
func (c *GitLabClient) ListCommits(ctx context.Context, owner, repo, ref string, limit int) ([]vcs.CommitInfo, error) {
	// TODO: Implement using go-gitlab
	return []vcs.CommitInfo{}, nil
}

// PostComment posts a comment on an MR
func (c *GitLabClient) PostComment(ctx context.Context, owner, repo string, number int, body string) (*vcs.Comment, error) {
	// TODO: Implement using go-gitlab
	return &vcs.Comment{
		ID:   "1",
		Body: body,
	}, nil
}

// PostReview submits an approval or change request
func (c *GitLabClient) PostReview(ctx context.Context, owner, repo string, number int, review *vcs.ReviewRequest) (*vcs.Review, error) {
	// TODO: Implement using go-gitlab
	return &vcs.Review{
		ID:    "1",
		State: review.Event,
	}, nil
}

// GetReviews retrieves approvals/reviews for an MR
func (c *GitLabClient) GetReviews(ctx context.Context, owner, repo string, number int) ([]vcs.Review, error) {
	// TODO: Implement using go-gitlab
	return []vcs.Review{}, nil
}

// GetComments retrieves comments on an MR
func (c *GitLabClient) GetComments(ctx context.Context, owner, repo string, number int) ([]vcs.Comment, error) {
	// TODO: Implement using go-gitlab
	return []vcs.Comment{}, nil
}

// GetBranch retrieves branch information
func (c *GitLabClient) GetBranch(ctx context.Context, owner, repo, branch string) (*vcs.BranchInfo, error) {
	// TODO: Implement using go-gitlab
	return &vcs.BranchInfo{
		Name: branch,
	}, nil
}

// RegisterWebhook registers a webhook
func (c *GitLabClient) RegisterWebhook(ctx context.Context, owner, repo string, webhook *vcs.Webhook) (string, error) {
	// TODO: Implement using go-gitlab
	return "webhook_id", nil
}

// GetPlatform returns the platform type
func (c *GitLabClient) GetPlatform() vcs.Platform {
	return vcs.GitLab
}

// ParseGitLabURL parses a GitLab MR URL
func ParseGitLabURL(url string) (owner, repo string, mrNumber int, err error) {
	// Pattern: https://gitlab.com/owner/repo/-/merge_requests/123
	pattern := regexp.MustCompile(`gitlab\.com/([^/]+)/([^/]+)/-/merge_requests/(\d+)`)
	matches := pattern.FindStringSubmatch(url)

	if len(matches) < 4 {
		return "", "", 0, fmt.Errorf("invalid GitLab MR URL: %s", url)
	}

	owner = matches[1]
	repo = matches[2]
	mrNumber, _ = strconv.Atoi(matches[3])

	return owner, repo, mrNumber, nil
}

// FormatGitLabURL formats owner, repo, and number into a URL
func FormatGitLabURL(owner, repo string, number int) string {
	return fmt.Sprintf("https://gitlab.com/%s/%s/-/merge_requests/%d", owner, repo, number)
}
