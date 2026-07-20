package vcs

import (
	"context"
	"io"
)

// Client defines the interface for interacting with version control systems
type Client interface {
	// Authentication
	Authenticate(ctx context.Context, token string) error
	IsAuthenticated(ctx context.Context) (bool, error)

	// Repository operations
	GetRepositoryInfo(ctx context.Context, owner, repo string) (*RepositoryInfo, error)

	// Pull Request operations
	GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error)
	ListPullRequestFiles(ctx context.Context, owner, repo string, number int) ([]FileDiff, error)
	GetPullRequestDiff(ctx context.Context, owner, repo string, number int) (io.ReadCloser, error)

	// Commits
	GetCommit(ctx context.Context, owner, repo, sha string) (*CommitInfo, error)
	ListCommits(ctx context.Context, owner, repo, ref string, limit int) ([]CommitInfo, error)

	// Comments and Reviews
	PostComment(ctx context.Context, owner, repo string, number int, body string) (*Comment, error)
	PostReview(ctx context.Context, owner, repo string, number int, review *ReviewRequest) (*Review, error)
	GetReviews(ctx context.Context, owner, repo string, number int) ([]Review, error)
	GetComments(ctx context.Context, owner, repo string, number int) ([]Comment, error)

	// Branches and refs
	GetBranch(ctx context.Context, owner, repo, branch string) (*BranchInfo, error)

	// Webhooks (optional, may not be supported by all platforms)
	RegisterWebhook(ctx context.Context, owner, repo string, webhook *Webhook) (string, error)

	// Platform identifier
	GetPlatform() Platform
}

// ReviewRequest represents a review submission request
type ReviewRequest struct {
	Body     string
	Event    string // "APPROVE", "REQUEST_CHANGES", "COMMENT"
	Comments []*ReviewComment
}

// ReviewComment represents a comment on a specific line
type ReviewComment struct {
	Path     string
	Line     int
	Body     string
}

// BranchInfo contains information about a branch
type BranchInfo struct {
	Name      string
	SHA       string
	Protected bool
}

// Webhook contains webhook configuration
type Webhook struct {
	URL    string
	Events []string
	Active bool
}

// NewClient creates a new VCS client for the specified platform
func NewClient(platform Platform) Client {
	switch platform {
	case GitHub:
		return newGitHubClient()
	case GitLab:
		return newGitLabClient()
	case Bitbucket:
		return newBitbucketClient()
	case Azure:
		return newAzureClient()
	default:
		return nil
	}
}

// Stub implementations for future development
func newGitHubClient() Client {
	return &gitHubClient{}
}

func newGitLabClient() Client {
	return &gitLabClient{}
}

func newBitbucketClient() Client {
	return &bitbucketClient{}
}

func newAzureClient() Client {
	return &azureClient{}
}

// Placeholder structs for each platform
type gitHubClient struct{}
type gitLabClient struct{}
type bitbucketClient struct{}
type azureClient struct{}

// Implement Client interface stubs
func (c *gitHubClient) GetPlatform() Platform { return GitHub }
func (c *gitLabClient) GetPlatform() Platform { return GitLab }
func (c *bitbucketClient) GetPlatform() Platform { return Bitbucket }
func (c *azureClient) GetPlatform() Platform { return Azure }

// Stub implementations return errors for now
var notImplementedErr = NewError("not_implemented", "Feature not yet implemented")
