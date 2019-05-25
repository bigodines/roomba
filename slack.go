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
	// SlackSvc is the Slack service layer
	SlackSvc struct {
		channelID string
		repos     map[string]bool
		user      string
		webhook   string
		client    *http.Client
		countdown map[string]string
	}

	// Entry represents a Roomba Slack entity
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
	// Nickname used for slack messages
	roombaUser = "Roomba"
	// Date format for reminders and countdowns
	layoutISO = "2006-01-02"
)

// Create a new Slack Service that can talk to and from Slack
func NewSlackSvc(appConfig config.Config) (SlackSvc, error) {
	return SlackSvc{
		webhook:   appConfig.Webhook,
		repos:     appConfig.Repos,
		countdown: appConfig.Countdown,
		channelID: appConfig.ChannelID,
		user:      roombaUser,
		client:    &http.Client{},
	}, nil
}

// Parse, filter and report results into slack channel
func (s *SlackSvc) Report(results []Record) error {
	// first check for messages/countdown/reminders from config
	// TODO: wg.add()
	reminders := s.GetMessages()

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
	prs := make([]string, 0)
	for _, entry := range relevant {
		line := entry.ToString()
		if len(line) > 0 {
			prs = append(prs, line)
		}
	}

	log.Debug().Msgf("%+v", prs)
	err := s.SendMessage(reminders, prs)
	if err != nil {
		return err
	}

	return nil
}

// GetMessages return active countdowns in a report friendly format
func (s *SlackSvc) GetMessages() []string {
	msgs := make([]string, 0)
	if len(s.countdown) < 1 {
		return msgs
	}
	// append countdowns in the future
	for k, v := range s.countdown {
		d, err := time.Parse(layoutISO, k)
		if err != nil {
			continue
		}
		daysUntil := int64(time.Until(d).Hours() / 24)
		if daysUntil > 0 {
			msgs = append(msgs, fmt.Sprintf("%s is *%d* days away!", v, daysUntil))
		}
	}

	return msgs
}

// Send individual slack message to configured slack channel
func (s *SlackSvc) SendMessage(reminders, prs []string) error {
	attachments := make([]map[string]interface{}, 1)
	attachments[0] = map[string]interface{}{"text": fmt.Sprintf("Howdy! Here's a list of *%d* PRs waiting to be reviewed and merged:", len(prs))}
	for _, v := range prs {
		attachments = append(attachments, map[string]interface{}{"text": v})
	}

	if len(reminders) > 0 {
		for _, v := range reminders {
			attachments = append(attachments, map[string]interface{}{"text": v})
		}
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

	resp, err := s.client.Post(s.webhook, "application/json", bytes.NewReader(payload))
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
