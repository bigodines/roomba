package roomba

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/olcolabs/roomba/config"
	"github.com/rs/zerolog"
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

func init() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

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
		att := res["attachments"].([]interface{})
		a := att[1].(map[string]interface{})
		assert.True(t, strings.Contains(a["text"].(string), "boom!"))

	}))

	s := SlackSvc{
		channelID: "chanID",
		user:      "Roomba",
		webhook:   testServer.URL,
		client:    &http.Client{},
	}

	fakeReport := ReportPayload{
		ChannelID: "123",
		Datetime:  time.Now(),
		PRs: []PullRequest{
			PullRequest{
				Title:     "boom!",
				Author:    "bigo",
				Permalink: "http://www.example.com",
			},
		},
	}

	err := s.SendMessage(fakeReport)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestReportCallback(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "POST" {
			t.Errorf("Expected POST")
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

		var req ReportPayload
		if err = json.Unmarshal(body, &req); err != nil {
			t.Error(err)
		}

		assert.Equal(t, "000ABC", req.ChannelID)

	}))

	s := SlackSvc{
		channelID:      "000ABC",
		user:           "Roomba",
		reportCallback: testServer.URL,
		client:         &http.Client{},
	}

	fakeReport := ReportPayload{
		ChannelID: "000ABC",
		Datetime:  time.Now(),
		PRs: []PullRequest{
			PullRequest{
				Title:     "boom!",
				Author:    "bigo",
				Permalink: "http://www.example.com",
			},
		},
	}

	err := s.ReportCallback(fakeReport)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetMessages(t *testing.T) {
	future := "2039-01-02"
	s, _ := NewSlackSvc(config.Config{
		Countdown: map[string]string{
			future: "This project shohuld be dead",
		},
	})

	d, _ := time.Parse(layoutISO, future)
	expectedDays := int64(time.Until(d).Hours() / 24)
	res := s.GetMessages()
	assert.Equal(t, 1, len(res))
	assert.True(t, strings.Contains(res[0].Text, fmt.Sprintf("%d", expectedDays)))
	assert.True(t, strings.Contains(res[0].Text, "This project shohuld be dead"))

	past := "1999-01-05"
	s, _ = NewSlackSvc(config.Config{
		Countdown: map[string]string{
			past: "This shouldn't be displayed",
		},
	})

	res = s.GetMessages()
	assert.Equal(t, 0, len(res))
}
