package models

import "strings"

// Graph represents the commit graph structure
type Graph struct {
	Commits    []*Commit          `json:"commits"`
	Branches   map[string][]*Commit `json:"branches"`
	Edges      []*Edge            `json:"edges"`
	Columns    []string           `json:"columns"` // branch names for columns
}

// Edge represents a connection between commits
type Edge struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Type   string `json:"type"`   // "parent", "merge"
	Column int    `json:"column"` // column index for visualization
}

// GraphRow represents a single row in the graph visualization
type GraphRow struct {
	Hash      string `json:"hash"`
	ShortHash string `json:"short_hash"`
	Author    string `json:"author"`
	Date      string `json:"date"`
	Message   string `json:"message"`
	Line      string `json:"line"` // ASCII art line
	Branches  []string `json:"branches,omitempty"`
	IsMerge   bool   `json:"is_merge"`
}

// NewGraph creates a new graph
func NewGraph() *Graph {
	return &Graph{
		Commits:  []*Commit{},
		Branches: make(map[string][]*Commit),
		Edges:    []*Edge{},
		Columns:  []string{},
	}
}

// AddCommit adds a commit to the graph
func (g *Graph) AddCommit(commit *Commit) {
	g.Commits = append(g.Commits, commit)
}

// AddEdge adds an edge to the graph
func (g *Graph) AddEdge(from, to, edgeType string, column int) {
	g.Edges = append(g.Edges, &Edge{
		From:   from,
		To:     to,
		Type:   edgeType,
		Column: column,
	})
}

// TruncateMessage truncates a commit message to specified length
func TruncateMessage(msg string, maxLen int) string {
	// Remove newlines
	msg = strings.ReplaceAll(msg, "\n", " ")
	msg = strings.TrimSpace(msg)
	
	if len(msg) <= maxLen {
		return msg
	}
	return msg[:maxLen-3] + "..."
}

// FormatDate formats a commit date
func FormatDate(t string) string {
	// Return short date format
	if len(t) >= 10 {
		return t[:10]
	}
	return t
}
