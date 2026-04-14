package visualize

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ColorScheme defines colors for different commit types
type ColorScheme struct {
	// Commit type colors
	Feat lipgloss.Style
	Fix  lipgloss.Style
	Docs lipgloss.Style
	Refactor lipgloss.Style
	Chore lipgloss.Style
	Ci lipgloss.Style
	Test lipgloss.Style
	Build lipgloss.Style
	Merge lipgloss.Style
	Default lipgloss.Style
	
	// UI colors
	Border lipgloss.Style
	Header lipgloss.Style
	Meta lipgloss.Style
	Bar lipgloss.Style
	BarEmpty lipgloss.Style
	
	// Enabled flag
	Enabled bool
}

// NewColorScheme creates a new color scheme
func NewColorScheme(enabled bool) *ColorScheme {
	if !enabled {
		// Return styles that don't apply colors
		return &ColorScheme{
			Enabled: false,
			Feat:    lipgloss.NewStyle(),
			Fix:     lipgloss.NewStyle(),
			Docs:    lipgloss.NewStyle(),
			Refactor: lipgloss.NewStyle(),
			Chore:   lipgloss.NewStyle(),
			Ci:      lipgloss.NewStyle(),
			Test:    lipgloss.NewStyle(),
			Build:   lipgloss.NewStyle(),
			Merge:   lipgloss.NewStyle(),
			Default: lipgloss.NewStyle(),
			Border:  lipgloss.NewStyle(),
			Header:  lipgloss.NewStyle(),
			Meta:    lipgloss.NewStyle(),
			Bar:     lipgloss.NewStyle(),
			BarEmpty: lipgloss.NewStyle(),
		}
	}
	
	return &ColorScheme{
		Enabled: true,
		
		// Commit type colors
		Feat:    lipgloss.NewStyle().Foreground(lipgloss.Color("#3B82F6")).Bold(true),      // Blue
		Fix:     lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E")).Bold(true),       // Green
		Docs:    lipgloss.NewStyle().Foreground(lipgloss.Color("#EAB308")).Bold(true),       // Yellow
		Refactor: lipgloss.NewStyle().Foreground(lipgloss.Color("#06B6D4")).Bold(true),      // Cyan
		Chore:   lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF")).Bold(true),       // Gray
		Ci:      lipgloss.NewStyle().Foreground(lipgloss.Color("#06B6D4")).Bold(true),        // Cyan
		Test:    lipgloss.NewStyle().Foreground(lipgloss.Color("#A855F7")).Bold(true),        // Purple
		Build:   lipgloss.NewStyle().Foreground(lipgloss.Color("#F97316")).Bold(true),       // Orange
		Merge:   lipgloss.NewStyle().Foreground(lipgloss.Color("#F97316")).Bold(true),       // Orange
		Default: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")),                  // White
		
		// UI colors
		Border:  lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")),
		Header:  lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true),
		Meta:    lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF")),
		Bar:     lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E")),
		BarEmpty: lipgloss.NewStyle().Foreground(lipgloss.Color("#374151")),
	}
}

// GetCommitStyle returns the appropriate style for a commit message
func (c *ColorScheme) GetCommitStyle(message string) lipgloss.Style {
	message = strings.TrimSpace(message)
	
	lowerMsg := strings.ToLower(message)
	
	switch {
	case strings.HasPrefix(lowerMsg, "feat"):
		return c.Feat
	case strings.HasPrefix(lowerMsg, "fix"):
		return c.Fix
	case strings.HasPrefix(lowerMsg, "docs"):
		return c.Docs
	case strings.HasPrefix(lowerMsg, "refactor"):
		return c.Refactor
	case strings.HasPrefix(lowerMsg, "chore"):
		return c.Chore
	case strings.HasPrefix(lowerMsg, "ci"):
		return c.Ci
	case strings.HasPrefix(lowerMsg, "test"):
		return c.Test
	case strings.HasPrefix(lowerMsg, "build"):
		return c.Build
	case strings.HasPrefix(lowerMsg, "merge") || strings.HasPrefix(lowerMsg, "pull"):
		return c.Merge
	default:
		return c.Default
	}
}

// RenderCommitMessage renders a commit message with color based on type
func (c *ColorScheme) RenderCommitMessage(message string) string {
	if !c.Enabled {
		return message
	}
	style := c.GetCommitStyle(message)
	return style.Render(message)
}

// RenderBar renders a progress bar with color
func (c *ColorScheme) RenderBar(filled, total, width int) string {
	if !c.Enabled || total == 0 {
		return ""
	}
	
	filledWidth := filled * width / total
	emptyWidth := width - filledWidth
	
	filledBar := strings.Repeat("█", filledWidth)
	emptyBar := strings.Repeat("░", emptyWidth)
	
	return c.Bar.Render(filledBar) + c.BarEmpty.Render(emptyBar)
}
