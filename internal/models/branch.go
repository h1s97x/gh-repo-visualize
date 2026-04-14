package models

// Branch represents a git branch
type Branch struct {
	Name      string
	IsCurrent bool
	IsRemote  bool
	PointsTo  string
	Upstream  string
}

// BranchGraph represents the branch visualization data
type BranchGraph struct {
	Branches      []*Branch
	CurrentBranch string
}

// TruncateBranchName truncates branch name for display
func TruncateBranchName(name string, maxLen int) string {
	if len(name) <= maxLen {
		return name
	}
	return name[:maxLen-3] + "..."
}
