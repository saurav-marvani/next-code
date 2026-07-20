package vcs

import (
	"time"
)

// VCS Platform types
type Platform string

const (
	GitHub    Platform = "github"
	GitLab    Platform = "gitlab"
	Bitbucket Platform = "bitbucket"
	Azure     Platform = "azure"
)

// FileDiff represents a file change in a PR/MR
type FileDiff struct {
	Path      string // File path in repository
	OldPath   string // Previous path (for renames)
	Status    FileStatus
	Additions int
	Deletions int
	Changes   int
	Patch     string // Raw unified diff
	OldBlob   string // Previous file content
	NewBlob   string // New file content
}

// FileStatus represents the type of change to a file
type FileStatus string

const (
	FileAdded    FileStatus = "added"
	FileModified FileStatus = "modified"
	FileDeleted  FileStatus = "deleted"
	FileRenamed  FileStatus = "renamed"
	FileCopied   FileStatus = "copied"
)

// PullRequest represents a PR/MR from a VCS platform
type PullRequest struct {
	ID          string
	Number      int
	Title       string
	Description string
	State       string // "open", "closed", "merged"
	Source      string // Source branch name
	Target      string // Target branch name
	Author      User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	MergedAt    *time.Time
	Commits     int
	Comments    int
	Reviews     int
}

// User represents a user/author
type User struct {
	ID    string
	Name  string
	Email string
	Login string
}

// Comment represents a review comment
type Comment struct {
	ID        string
	Path      string
	Line      int
	Body      string
	Author    User
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Review represents a code review
type Review struct {
	ID    string
	State string // "approved", "changes_requested", "commented"
	Body  string
	User  User
	Date  time.Time
}

// DiffHunk represents a section of a diff
type DiffHunk struct {
	OldStart  int
	OldLines  int
	NewStart  int
	NewLines  int
	Header    string
	Lines     []DiffLine
}

// DiffLine represents a single line in a diff
type DiffLine struct {
	Type      LineType
	OldLineNo int
	NewLineNo int
	Content   string
}

// LineType represents the type of change for a line
type LineType string

const (
	LineAdded   LineType = "added"
	LineRemoved LineType = "removed"
	LineContext LineType = "context"
)

// CommitInfo represents a commit
type CommitInfo struct {
	SHA     string
	Message string
	Author  User
	Date    time.Time
	Files   []FileDiff
}

// RepositoryInfo contains metadata about a repository
type RepositoryInfo struct {
	Owner       string
	Name        string
	URL         string
	DefaultBranch string
	Description string
	Private     bool
}
