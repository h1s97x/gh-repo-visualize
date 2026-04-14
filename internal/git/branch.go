package git

import (
	"fmt"
	"strings"

	"github.com/h1s97x/gh-repo-visualize/internal/errors"
	"github.com/h1s97x/gh-repo-visualize/internal/models"
)

// BranchInfo represents information about a branch
type BranchInfo struct {
	Name      string
	IsHead    bool
	IsRemote  bool
	PointsTo  string // commit hash
	Timestamp int64  // commit timestamp the branch points to
}

// GetBranchGraph retrieves branch graph information
func (c *Client) GetBranchGraph(limit int) (*models.BranchGraph, error) {
	if !c.IsGitRepo() {
		return nil, errors.NotGitRepo(c.repoPath)
	}

	// Get current branch
	currentBranch, _ := c.GetCurrentBranch()

	// Get all branches with their info
	output, err := c.runGit("for-each-ref", "--format=%(refname:short)|%(objectname)|%(upstream:short)|%(creatordate:unix)", "--sort=-creatordate", "refs/heads/", "refs/remotes/")
	if err != nil {
		return nil, err
	}

	graph := &models.BranchGraph{
		Branches:      make([]*models.Branch, 0),
		CurrentBranch: currentBranch,
	}

	// Track unique branches
	seenBranches := make(map[string]bool)

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			continue
		}

		name := parts[0]
		hash := parts[1]
		upstream := ""
		if len(parts) > 2 {
			upstream = parts[2]
		}

		// Skip duplicates
		if seenBranches[name] {
			continue
		}
		seenBranches[name] = true

		isRemote := strings.HasPrefix(name, "origin/") || strings.Contains(name, "/")

		graph.Branches = append(graph.Branches, &models.Branch{
			Name:      name,
			IsCurrent: name == currentBranch,
			IsRemote:  isRemote,
			PointsTo:  hash,
			Upstream: upstream,
		})
	}

	// Limit branches if specified
	if limit > 0 && len(graph.Branches) > limit {
		graph.Branches = graph.Branches[:limit]
	}

	return graph, nil
}

// GetBranchCommits gets commits for a specific branch with branch tracking
func (c *Client) GetBranchCommits(branch string, limit int) ([]*models.Commit, error) {
	args := []string{"log", "--pretty=format:%H|%h|%an|%ae|%at|%s|%P", branch}

	if limit > 0 {
		args = append(args, fmt.Sprintf("-%d", limit))
	}

	output, err := c.runGit(args...)
	if err != nil {
		return nil, err
	}

	if output == "" {
		return nil, errors.NoCommits(branch)
	}

	return parseCommits(output), nil
}
