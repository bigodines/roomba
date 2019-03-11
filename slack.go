package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go4.org/sort"

	"github.com/bigodines/roomba/config"
	humanize "github.com/dustin/go-humanize"
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
	roombaUser = "Roomba"
)

// Create a new Slack Service that can talk to and from Slack
func NewSlackSvc(appConfig config.Config) (SlackSvc, error) {
	return SlackSvc{
		webhook:   appConfig.Webhook,
		repos:     appConfig.Repos,
		channelID: appConfig.ChannelID,
		user:      roombaUser,
	}, nil
}

// Parse, filter and report github results into slack channel
func (s *SlackSvc) Report(results []Record) error {
	relevant := make([]*Entry, 0)
	// filter relevant Pull Requests
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

	// Oldest first
	sort.Slice(relevant, func(i, j int) bool {
		return relevant[i].UpdatedAt.Before(relevant[j].UpdatedAt)
	})

	// Create Report
	msg := make([]string, 0)
	for _, entry := range relevant {
		line := entry.ToString()
		if len(line) > 0 {
			msg = append(msg, line)
		}
	}

	log.Debug().Msgf("%+v", msg)
	err := s.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

// Send individual slack message to configured slack channel
func (s *SlackSvc) SendMessage(contents []string) error {
	attachments := make([]map[string]interface{}, 1)
	attachments[0] = map[string]interface{}{"text": fmt.Sprintf("Howdy! Here's a list of *%d* PRs waiting to be reviewed and merged:", len(contents))}
	for _, v := range contents {
		attachments = append(attachments, map[string]interface{}{"text": v})
	}

	message := map[string]interface{}{
		"channel":     s.channelID,
		"username":    s.user,
		"icon_emoji":  ":robot_face:",
		"attachments": attachments,
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

	log.Info().Msgf("Message successfully sent to channel %s", s.channelID)
	return nil
}

// Printable format of an Entry
func (e *Entry) ToString() string {
	age := humanize.Time(e.UpdatedAt)
	return fmt.Sprintf("*%s* | %s | %s\n\t [%s] \"<%s|%s>\"", e.Repository, e.Author, age, e.Labels, e.Permalink, e.Title)
}
