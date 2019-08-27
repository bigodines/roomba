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
		ChannelID string        `json:"channel_id"`
		Datetime  time.Time     `json:"datetime"`
		PRs       []PullRequest `json:"prs"`
		Reminders []Reminder    `json:"reminders"`
	}

	Reminder struct {
		Date time.Time `json:"date"`
		Text string    `json:"text"`
	}

	// PullRequest represents the important fields for Roomba in a Pull Request
	PullRequest struct {
		Repository string    `json:"repository"`
		Author     string    `json:"author"`
		UpdatedAt  time.Time `json:"updated_at"`
		Labels     string    `json:"labels"`
		Title      string    `json:"title"`
		Permalink  string    `json:"permalink"`
	}
)

// Printable format of a PullRequest
func (e *PullRequest) ToString() string {
	age := humanize.Time(e.UpdatedAt)
	return fmt.Sprintf("*%s* | %s | %s\n\t [%s] \"<%s|%s>\"", e.Repository, e.Author, age, e.Labels, e.Permalink, e.Title)
}
