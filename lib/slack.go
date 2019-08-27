package roomba

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go4.org/sort"

	"github.com/olcolabs/roomba/config"
	"github.com/rs/zerolog/log"
)

type (
	// SlackSvc is the Slack service layer
	SlackSvc struct {
		channelID      string
		repos          map[string]bool
		user           string
		webhook        string
		client         *http.Client
		countdown      map[string]string
		reportCallback string
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
		webhook:        appConfig.Webhook,
		repos:          appConfig.Repos,
		countdown:      appConfig.Countdown,
		channelID:      appConfig.ChannelID,
		reportCallback: appConfig.ReportCallback,
		user:           roombaUser,
		client:         &http.Client{},
	}, nil
}

// Parse, filter and report results into slack channel
func (s *SlackSvc) Report(results []Record) error {
	// first check for messages/countdown/reminders from config
	// TODO: wg.add()
	reminders := s.GetMessages()

	relevant := make([]PullRequest, 0)
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
		relevant = append(relevant, PullRequest{
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

	report := ReportPayload{
		ChannelID: s.channelID,
		Datetime:  time.Now(),
		PRs:       relevant,
		Reminders: reminders,
	}

	log.Debug().Msgf("%+v", report)
	err := s.SendMessage(report)
	if err != nil {
		return err
	}

	if len(s.reportCallback) > 0 {
		err = s.ReportCallback(report)
		if err != nil {
			return err
		}
	}
	return nil
}

// Send individual slack message to configured slack channel
func (s *SlackSvc) SendMessage(report ReportPayload) error {
	attachments := make([]map[string]interface{}, 1)
	attachments[0] = map[string]interface{}{"text": fmt.Sprintf("Howdy! Here's a list of *%d* PRs waiting to be reviewed and merged:", len(report.PRs))}
	for _, v := range report.PRs {
		attachments = append(attachments, map[string]interface{}{"text": v.ToString()})
	}

	if len(report.Reminders) > 0 {
		for _, v := range report.Reminders {
			attachments = append(attachments, map[string]interface{}{"text": v.Text})
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

	defer resp.Body.Close()

	log.Info().Msgf("Message successfully sent to channel %s", s.channelID)
	return nil
}

func (s *SlackSvc) ReportCallback(report ReportPayload) error {
	payload, err := json.Marshal(report)
	if err != nil {
		log.Error().Err(err).Msg("Failed to serialize callback payload")
		return err
	}

	resp, err := s.client.Post(s.reportCallback, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to post payload to callback endpoint")
		return errors.New("Failed to post")
	}
	defer resp.Body.Close()

	log.Info().Str("callback_url", s.reportCallback).Msgf("Report callback sent")
	return nil
}

// GetMessages return active countdowns in a report friendly format
func (s *SlackSvc) GetMessages() []Reminder {
	msgs := make([]Reminder, 0)
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
			msgs = append(msgs, Reminder{
				Date: d,
				Text: fmt.Sprintf("Friendly Reminder: \"%s\" is *%d* days away!", v, daysUntil),
			})
		}
	}

	return msgs
}
