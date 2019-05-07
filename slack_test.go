package main

import (
	"testing"

	"github.com/bigodines/roomba/config"
	"github.com/stretchr/testify/assert"
)

var (
	fakeConfig = config.Config{
		Webhook: "http://foo.bar.baz/123",
		Repos: map[string]bool{
			"repo1": true,
			"repo2": true,
			"repo3": true,
		},
		ChannelID:    "123",
		Organization: "bigodines",
	}
)

func TestConstructor(t *testing.T) {
	s, _ := NewSlackSvc(fakeConfig)
	assert.Equal(t, "http://foo.bar.baz/123", s.webhook)
}
