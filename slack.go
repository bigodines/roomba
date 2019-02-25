package main

import (
	"fmt"
	"time"

	"github.com/bigodines/roomgo/config"
	"github.com/nlopes/slack"
)

type (
	SlackSvc struct {
		client    *slack.Client
		repos     map[string]bool
		channelID string
	}

	Entry struct {
		Repository string
		Author     string
		UpdatedAt  time.Time
		Labels     []string
		Title      string
		Permalink  string
	}
)

// Create a new Slack Service that can talk to and from Slack
func NewSlackSvc(token string, appConfig config.Config) (SlackSvc, error) {
	c := slack.New(token)
	return SlackSvc{
		client:    c,
		repos:     appConfig.Repos,
		channelID: appConfig.ChannelID,
	}, nil
}

// Parse, filter and report github results into slack channel
func (s *SlackSvc) Report(results []Record) error {
	relevant := make([]*Entry, 0)
	for _, v := range results {
		pr := v.Node.PullRequest
		_, exists := s.repos[pr.HeadRepository.Name]
		if exists {
			relevant = append(relevant, &Entry{
				Title:      pr.Title,
				Author:     pr.Author.Login,
				Permalink:  pr.Permalink,
				Repository: pr.HeadRepository.Name,
				// TODO: finish
			})
		}
		printJSON(v)
	}
	return nil
}

// Send individual slack message to configured slack channel
func (s *SlackSvc) SendMessage(contents string) error {
	attachment := slack.Attachment{
		Pretext: "some pretext",
		Text:    "some text",
		// Uncomment the following part to send a field too
		/*
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "a",
					Value: "no",
				},
			},
		*/
	}

	channelID, timestamp, err := s.client.PostMessage("CHANNEL_ID", slack.MsgOptionText("Some text", false), slack.MsgOptionAttachments(attachment))
	if err != nil {
		return err
	}

	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return nil
}
