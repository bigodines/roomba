package roomba

import (
	"fmt"
	"time"

	humanize "github.com/dustin/go-humanize"
)

type (
	// ReportPayload is the payload Roomba will post to the report callback URL (see config.go for details)
	// It contains a mix of Config{}, Record{} and runtime information
	ReportPayload struct {
		ChannelID string        `json:"channel_id,omitempty"`
		Datetime  time.Time     `json:"datetime,omitempty"`
		PRs       []PullRequest `json:"prs,omitempty"`
		Reminders []Reminder    `json:"reminders,omitempty"`
	}

	Reminder struct {
		Date time.Time `json:"date,omitempty"`
		Text string    `json:"text,omitempty"`
	}

	// PullRequest represents the important fields for Roomba in a Pull Request
	PullRequest struct {
		Repository string    `json:"repository,omitempty"`
		Author     string    `json:"author,omitempty"`
		UpdatedAt  time.Time `json:"updated_at,omitempty"`
		Labels     string    `json:"labels,omitempty"`
		Title      string    `json:"title,omitempty"`
		Permalink  string    `json:"permalink,omitempty"`
	}
)

// Printable format of a PullRequest
func (e *PullRequest) ToString() string {
	age := humanize.Time(e.UpdatedAt)
	return fmt.Sprintf("*%s* | %s | %s\n\t [%s] \"<%s|%s>\"", e.Repository, e.Author, age, e.Labels, e.Permalink, e.Title)
}
