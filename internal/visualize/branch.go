package visualize

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/h1s97x/gh-repo-visualize/internal/models"
)

// BranchGraphRenderer renders branch graphs
type BranchGraphRenderer struct {
	options RenderOptions
	colors  *ColorScheme
}

// NewBranchGraphRenderer creates a new branch graph renderer
func NewBranchGraphRenderer(opts RenderOptions) *BranchGraphRenderer {
	if opts.Width == 0 {
		opts.Width = 80
	}
	colors := NewColorScheme(opts.Color)
	return &BranchGraphRenderer{options: opts, colors: colors}
}

// Render renders the branch graph
func (r *BranchGraphRenderer) Render(graph *models.BranchGraph) string {
	if graph == nil || len(graph.Branches) == 0 {
		return "No branches to display"
	}

	var sb strings.Builder
	width := r.options.Width

	// Header
	border := r.colors.Border.Render("╭") + strings.Repeat("─", width-2) + r.colors.Border.Render("╮")
	sb.WriteString(border + "\n")

	title := " Branch Overview "
	if r.colors.Enabled {
		title = r.colors.Header.Render(title)
	}
	sb.WriteString(fmt.Sprintf("│%s│\n", centered(title, width-2)))

	borderMid := r.colors.Border.Render("├") + strings.Repeat("─", width-2) + r.colors.Border.Render("┤")
	sb.WriteString(borderMid + "\n")

	// Current branch info
	currentInfo := fmt.Sprintf(" Current Branch: %s ", graph.CurrentBranch)
	if r.colors.Enabled {
		currentInfo = r.colors.Feat.Render(currentInfo)
	}
	sb.WriteString(fmt.Sprintf("│%s│\n", padRight(currentInfo, width-2)))

	sb.WriteString(borderMid + "\n")

	// Branch list - local branches first
	sb.WriteString(r.renderBranchSection("Local Branches", getLocalBranches(graph.Branches)))
	sb.WriteString(borderMid + "\n")
	sb.WriteString(r.renderBranchSection("Remote Branches", getRemoteBranches(graph.Branches)))

	// Footer
	borderBottom := r.colors.Border.Render("╰") + strings.Repeat("─", width-2) + r.colors.Border.Render("╯")
	sb.WriteString(borderBottom + "\n")

	return sb.String()
}

// RenderASCII renders a simple ASCII tree of branches
func (r *BranchGraphRenderer) RenderASCII(graph *models.BranchGraph) string {
	if graph == nil || len(graph.Branches) == 0 {
		return "No branches to display"
	}

	var sb strings.Builder

	sb.WriteString("Branch Graph:\n")
	sb.WriteString(strings.Repeat("=", 40) + "\n\n")

	// Group branches
	local := getLocalBranches(graph.Branches)
	remote := getRemoteBranches(graph.Branches)

	sb.WriteString("Local Branches:\n")
	for _, branch := range local {
		marker := "  "
		nameStyle := ""
		if branch.IsCurrent {
			marker = "* "
			nameStyle = r.colors.Feat.Render(branch.Name)
		} else {
			nameStyle = r.colors.Meta.Render(branch.Name)
		}
		sb.WriteString(fmt.Sprintf("%s%s (%s)\n", marker, nameStyle, branch.PointsTo[:7]))
	}

	if len(remote) > 0 {
		sb.WriteString("\nRemote Branches:\n")
		for _, branch := range remote {
			nameStyle := r.colors.Meta.Render(branch.Name)
			tracking := ""
			if branch.Upstream != "" {
				tracking = r.colors.Fix.Render(" -> " + branch.Upstream)
			}
			sb.WriteString(fmt.Sprintf("  %s (%s)%s\n", nameStyle, branch.PointsTo[:7], tracking))
		}
	}

	return sb.String()
}

// RenderJSON renders as JSON
func (r *BranchGraphRenderer) RenderJSON(graph *models.BranchGraph) string {
	if graph == nil {
		return "{}"
	}

	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString(fmt.Sprintf(`  "current_branch": "%s",
  "branches": [
`, graph.CurrentBranch))

	for i, branch := range graph.Branches {
		sb.WriteString(fmt.Sprintf(`    {
      "name": "%s",
      "is_current": %t,
      "is_remote": %t,
      "hash": "%s",
      "upstream": "%s"
    }`,
			escapeJSON(branch.Name),
			branch.IsCurrent,
			branch.IsRemote,
			branch.PointsTo,
			branch.Upstream,
		))
		if i < len(graph.Branches)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("  ]\n")
	sb.WriteString("}\n")
	return sb.String()
}

// Helper functions

func getLocalBranches(branches []*models.Branch) []*models.Branch {
	var local []*models.Branch
	for _, b := range branches {
		if !b.IsRemote {
			local = append(local, b)
		}
	}
	return local
}

func getRemoteBranches(branches []*models.Branch) []*models.Branch {
	var remote []*models.Branch
	for _, b := range branches {
		if b.IsRemote {
			remote = append(remote, b)
		}
	}
	return remote
}

func (r *BranchGraphRenderer) renderBranchSection(title string, branches []*models.Branch) string {
	if len(branches) == 0 {
		return ""
	}

	var sb strings.Builder

	sectionTitle := " " + title + " "
	if r.colors.Enabled {
		sectionTitle = r.colors.Header.Render(sectionTitle)
	}
	sb.WriteString(fmt.Sprintf("│%s│\n", padRight(sectionTitle, r.options.Width-2)))

	for _, branch := range branches {
		marker := r.colors.Meta.Render("○")
		nameStyle := r.colors.Meta.Render(branch.Name)
		info := ""

		if branch.IsCurrent {
			marker = r.colors.Feat.Render("●")
			nameStyle = r.colors.Feat.Render(branch.Name)
		}

		if branch.Upstream != "" {
			info = r.colors.Fix.Render(" -> " + branch.Upstream)
		}

		hashStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E"))
		if r.colors.Enabled {
			hashStyle = hashStyle.Bold(true)
		}
		hash := hashStyle.Render(branch.PointsTo[:7])

		line := fmt.Sprintf("  %s %s %s%s", marker, nameStyle, hash, info)
		sb.WriteString(fmt.Sprintf("│%s│\n", padRight(line, r.options.Width-2)))
	}

	return sb.String()
}

func centered(s string, width int) string {
	pad := (width - len(s)) / 2
	if pad < 0 {
		pad = 0
	}
	return strings.Repeat(" ", pad) + s + strings.Repeat(" ", width-pad-len(s))
}
