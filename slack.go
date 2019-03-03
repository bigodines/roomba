package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bigodines/roomgo/config"
	"github.com/rs/zerolog/log"
)

type (
	SlackSvc struct {
		channelID string
		repos     map[string]bool
		user      string
		webhook   string
	}

	Entry struct {
		Repository string
		Author     string
		UpdatedAt  time.Time
		Labels     string
		Title      string
		Permalink  string
	}
)

const (
	roomgoUser = "Roomba"
)

// Create a new Slack Service that can talk to and from Slack
func NewSlackSvc(webhookURL string, appConfig config.Config) (SlackSvc, error) {
	return SlackSvc{
		webhook:   webhookURL,
		repos:     appConfig.Repos,
		channelID: appConfig.ChannelID,
		user:      roomgoUser,
	}, nil
}

// Parse, filter and report github results into slack channel
func (s *SlackSvc) Report(results []Record) error {
	relevant := make([]*Entry, 0)
	// filter
	for _, v := range results {
		pr := v.Node.PullRequest
		_, exists := s.repos[pr.HeadRepository.Name]
		if !exists {
			// we don't care about this repository
			continue
		}
		// create and add a report entry
		l := PrintableLabels(pr.Labels)
		relevant = append(relevant, &Entry{
			Title:      pr.Title,
			Author:     pr.Author.Login,
			Permalink:  pr.Permalink,
			Repository: pr.HeadRepository.Name,
			Labels:     l,
			UpdatedAt:  pr.UpdatedAt,
		})
	}

	// report
	msg := make([]string, 0)
	for _, entry := range relevant {
		line := entry.ToString()
		if len(line) > 0 {
			msg = append(msg, line)
		}
	}
	// TODO: replace w log lib
	fmt.Printf("%+v", msg)
	err := s.SendMessage(strings.Join(msg[:], "\n"))
	if err != nil {
		return err
	}
	// TODO: remove
	return nil
}

// Send individual slack message to configured slack channel
func (s *SlackSvc) SendMessage(contents string) error {
	message := map[string]interface{}{
		"text":       "",
		"channel":    s.channelID,
		"username":   s.user,
		"icon_emoji": ":robot_face:",
		"attachments": []map[string]interface{}{
			{
				"fallback":    contents,
				"color":       "good",
				"author_name": "Roomba",
				"text":        contents,
			},
		},
	}

	payload, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to serialize Slack payload")
		return err
	}

	resp, err := http.Post(s.webhook, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to serialize Slack payload: %v", err)
		return err
	}

	resp.Body.Close()

	fmt.Printf("Message successfully sent to channel %s", s.channelID)
	return nil
}

// Printable format of an Entry
func (e *Entry) ToString() string {
	// TODO: improve format
	return fmt.Sprintf("[%s] %s - \"%s\" (%s)", e.Labels, e.Repository, e.Title, e.Author)
}
