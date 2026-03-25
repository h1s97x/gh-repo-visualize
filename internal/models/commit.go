package models

import "time"

// Commit represents a single Git commit
type Commit struct {
	Hash      string    `json:"hash"`
	ShortHash string    `json:"short_hash"`
	Author    string    `json:"author"`
	Email     string    `json:"email"`
	Date      time.Time `json:"date"`
	Message   string    `json:"message"`
	Parents   []string  `json:"parents,omitempty"`
	Branch    string    `json:"branch,omitempty"`
}

// CommitList represents a collection of commits
type CommitList struct {
	Commits  []*Commit `json:"commits"`
	Count    int       `json:"count"`
	Branch   string    `json:"branch,omitempty"`
	RepoPath string    `json:"repo_path"`
}

// NewCommitList creates a new commit list
func NewCommitList() *CommitList {
	return &CommitList{
		Commits: []*Commit{},
		Count:   0,
	}
}

// Add adds a commit to the list
func (cl *CommitList) Add(commit *Commit) {
	cl.Commits = append(cl.Commits, commit)
	cl.Count = len(cl.Commits)
}

// FilterByAuthor filters commits by author
func (cl *CommitList) FilterByAuthor(author string) *CommitList {
	result := NewCommitList()
	for _, c := range cl.Commits {
		if c.Author == author {
			result.Add(c)
		}
	}
	return result
}

// FilterByBranch filters commits by branch
func (cl *CommitList) FilterByBranch(branch string) *CommitList {
	result := NewCommitList()
	for _, c := range cl.Commits {
		if c.Branch == branch {
			result.Add(c)
		}
	}
	return result
}

// Limit returns the first n commits
func (cl *CommitList) Limit(n int) *CommitList {
	result := NewCommitList()
	for i, c := range cl.Commits {
		if i >= n {
			break
		}
		result.Add(c)
	}
	return result
}
