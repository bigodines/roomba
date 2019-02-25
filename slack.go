package main

import "github.com/nlopes/slack"

type (
	SlackSvc struct {
		client interface{}
	}
)

// Create a new Slack Service that can talk to and from Slack
func NewSlackSvc(token string) (SlackSvc, error) {
	c := slack.New(token)
	return SlackSvc{
		client: c,
	}, nil
}

// Parse, filter and report github results into slack channel
func (ss *SlackSvc) Report(results []Record) error {
	// TODO: slack results
	for _, v := range results {
		printJSON(v)
	}
	return nil
}

// Send individual slack message to configured slack channel
func (s *SlackSvc) SendMessage(contents string) error {
	return nil
}
