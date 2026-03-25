package git

import (
	"strconv"
	"strings"
	"time"

	"github.com/h1s97x/gh-repo-visualize/internal/models"
)

// LogOptions represents options for git log
type LogOptions struct {
	Limit  int
	Branch string
	Author string
	Since  string
	Until  string
}

// parseCommits parses git log output into Commit structs
func parseCommits(output string) []*models.Commit {
	var commits []*models.Commit
	
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		commit := parseCommitLine(line)
		if commit != nil {
			commits = append(commits, commit)
		}
	}
	
	return commits
}

// parseCommitLine parses a single commit line
// Format: %H|%h|%an|%ae|%at|%s|%P
func parseCommitLine(line string) *models.Commit {
	parts := strings.SplitN(line, "|", 7)
	if len(parts) < 6 {
		return nil
	}
	
	timestamp, _ := strconv.ParseInt(parts[4], 10, 64)
	
	commit := &models.Commit{
		Hash:      parts[0],
		ShortHash: parts[1],
		Author:    parts[2],
		Email:     parts[3],
		Date:      time.Unix(timestamp, 0),
		Message:   parts[5],
		Parents:   []string{},
	}
	
	// Parse parents if present
	if len(parts) > 6 && parts[6] != "" {
		parents := strings.Fields(parts[6])
		commit.Parents = parents
	}
	
	return commit
}

// GetCommitGraph returns a graph representation of commits
func (c *Client) GetCommitGraph(opts *LogOptions) (*models.Graph, error) {
	commits, err := c.GetCommits(opts)
	if err != nil {
		return nil, err
	}
	
	graph := models.NewGraph()
	
	// Build commit hash map for quick lookup
	commitMap := make(map[string]*models.Commit)
	for _, commit := range commits {
		commitMap[commit.Hash] = commit
		graph.AddCommit(commit)
	}
	
	// Build edges
	for _, commit := range commits {
		for i, parent := range commit.Parents {
			edgeType := "parent"
			if i > 0 {
				edgeType = "merge"
			}
			graph.AddEdge(commit.Hash, parent, edgeType, 0)
		}
	}
	
	return graph, nil
}
