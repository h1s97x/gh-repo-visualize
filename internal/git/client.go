package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/h1s97x/gh-repo-visualize/internal/errors"
	"github.com/h1s97x/gh-repo-visualize/internal/models"
)

// Client wraps git command operations
type Client struct {
	repoPath string
}

// NewClient creates a new git client
func NewClient(repoPath string) *Client {
	if repoPath == "" {
		repoPath, _ = os.Getwd()
	}
	return &Client{repoPath: repoPath}
}

// IsGitRepo checks if the path is a git repository
func (c *Client) IsGitRepo() bool {
	gitDir := filepath.Join(c.repoPath, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		return true
	}
	// Check if we're in a git repo using git command
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = c.repoPath
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

// runGit executes a git command and returns the output
func (c *Client) runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = c.repoPath
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return "", errors.GitError(strings.Join(args, " "), stderr.String(), err)
	}
	
	return stdout.String(), nil
}

// GetCurrentBranch returns the current branch name
func (c *Client) GetCurrentBranch() (string, error) {
	output, err := c.runGit("branch", "--show-current")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// GetBranches returns all branch names
func (c *Client) GetBranches() ([]string, error) {
	output, err := c.runGit("branch", "--list", "--all", "--format=%(refname:short)")
	if err != nil {
		return nil, err
	}
	
	branches := []string{}
	for _, b := range strings.Split(output, "\n") {
		b = strings.TrimSpace(b)
		if b != "" {
			branches = append(branches, b)
		}
	}
	return branches, nil
}

// GetCommits retrieves commits with optional filters
func (c *Client) GetCommits(opts *LogOptions) ([]*models.Commit, error) {
	if !c.IsGitRepo() {
		return nil, errors.NotGitRepo(c.repoPath)
	}

	args := []string{"log", "--pretty=format:%H|%h|%an|%ae|%at|%s|%P", "--all"}
	
	// Add limit
	if opts != nil && opts.Limit > 0 {
		args = append(args, fmt.Sprintf("-%d", opts.Limit))
	}
	
	// Add branch filter
	if opts != nil && opts.Branch != "" {
		// Remove --all and add specific branch
		args = args[:3] // Keep base args without --all
		args = append(args, opts.Branch)
	}
	
	// Add author filter
	if opts != nil && opts.Author != "" {
		args = append(args, "--author="+opts.Author)
	}
	
	output, err := c.runGit(args...)
	if err != nil {
		return nil, err
	}
	
	if output == "" {
		return nil, errors.NoCommits(opts.Branch)
	}
	
	return parseCommits(output), nil
}

// GetCommitStats returns statistics about commits
func (c *Client) GetCommitStats(opts *LogOptions) (*models.Stats, error) {
	commits, err := c.GetCommits(opts)
	if err != nil {
		return nil, err
	}
	
	stats := models.NewStats()
	stats.Calculate(commits)
	return stats, nil
}
