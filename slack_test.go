package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func TestSendMessage(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "POST" {
			t.Errorf("Expected POST")
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

		var res map[string]interface{}
		if err = json.Unmarshal(body, &res); err != nil {
			t.Error(err)
		}

		assert.Equal(t, "chanID", res["channel"])
		assert.Equal(t, "Roomba", res["username"])

	}))

	s := SlackSvc{
		channelID: "chanID",
		user:      "Roomba",
		webhook:   testServer.URL,
		client:    &http.Client{},
	}

	s.SendMessage([]string{"hi"})
}
