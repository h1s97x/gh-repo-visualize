package models

import (
	"fmt"
	"time"
)

// Stats represents commit statistics
type Stats struct {
	TotalCommits   int            `json:"total_commits"`
	Authors        []*AuthorStats `json:"authors"`
	DateRange      *DateRange     `json:"date_range,omitempty"`
	ByAuthor       map[string]int `json:"by_author,omitempty"`
	ByDate         map[string]int `json:"by_date,omitempty"`
	ByDayOfWeek    map[string]int `json:"by_day_of_week,omitempty"`
	ByHour         map[int]int    `json:"by_hour,omitempty"`
	AvgPerDay      float64        `json:"avg_per_day"`
	MostActiveDay  string         `json:"most_active_day,omitempty"`
	MostActiveHour int            `json:"most_active_hour,omitempty"`
}

// AuthorStats represents statistics for a single author
type AuthorStats struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Commits int    `json:"commits"`
	Percent string `json:"percent"`
	Latest  string `json:"latest_commit_date,omitempty"`
	First   string `json:"first_commit_date,omitempty"`
}

// DateRange represents a date range
type DateRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// NewStats creates a new stats object
func NewStats() *Stats {
	return &Stats{
		Authors:     []*AuthorStats{},
		ByAuthor:    make(map[string]int),
		ByDate:      make(map[string]int),
		ByDayOfWeek: make(map[string]int),
		ByHour:      make(map[int]int),
	}
}

// Calculate computes statistics from commits
func (s *Stats) Calculate(commits []*Commit) {
	s.TotalCommits = len(commits)
	if len(commits) == 0 {
		return
	}

	// Track first and last dates
	var firstDate, lastDate time.Time
	authorMap := make(map[string]*AuthorStats)

	for _, c := range commits {
		// Date tracking
		if firstDate.IsZero() || c.Date.Before(firstDate) {
			firstDate = c.Date
		}
		if lastDate.IsZero() || c.Date.After(lastDate) {
			lastDate = c.Date
		}

		// Author stats
		key := c.Author + " <" + c.Email + ">"
		if _, exists := authorMap[key]; !exists {
			authorMap[key] = &AuthorStats{
				Name:  c.Author,
				Email: c.Email,
			}
		}
		authorMap[key].Commits++

		// By date (day)
		dateKey := c.Date.Format("2006-01-02")
		s.ByDate[dateKey]++

		// By day of week
		dayKey := c.Date.Weekday().String()
		s.ByDayOfWeek[dayKey]++

		// By hour
		s.ByHour[c.Date.Hour()]++
	}

	// Convert author map to slice
	for _, as := range authorMap {
		as.Percent = calculatePercent(as.Commits, s.TotalCommits)
		s.Authors = append(s.Authors, as)
		s.ByAuthor[as.Name] = as.Commits
	}

	// Calculate average per day
	if !firstDate.IsZero() && !lastDate.IsZero() {
		days := lastDate.Sub(firstDate).Hours() / 24
		if days > 0 {
			s.AvgPerDay = float64(s.TotalCommits) / (days + 1)
		}
		s.DateRange = &DateRange{
			Start: firstDate.Format("2006-01-02"),
			End:   lastDate.Format("2006-01-02"),
		}
	}

	// Find most active day
	maxDay := 0
	for day, count := range s.ByDayOfWeek {
		if count > maxDay {
			maxDay = count
			s.MostActiveDay = day
		}
	}

	// Find most active hour
	maxHour := 0
	for hour, count := range s.ByHour {
		if count > maxHour {
			maxHour = count
			s.MostActiveHour = hour
		}
	}
}

func calculatePercent(value, total int) string {
	if total == 0 {
		return "0%"
	}
	percent := float64(value) / float64(total) * 100
	return fmt.Sprintf("%.1f%%", percent)
}
