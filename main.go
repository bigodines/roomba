package main

import (
	"context"
	"os"

	"github.com/bigodines/roomgo/config"
	"github.com/rs/zerolog/log"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	conf, err := config.Load(getEnv())
	if err != nil {
		panic(err)
	}

	// Github auth
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)

	httpClient := oauth2.NewClient(context.Background(), src)
	ghClient := githubv4.NewClient(httpClient)

	slackSvc, err := NewSlackSvc(conf)
	if err != nil {
		log.Fatal().Err(err).Msg("Can create slack service")
	}

	// GraphQL query
	var q struct {
		Search Search `graphql:"search(query:$query, type:ISSUE, first:30)"`
	}
	vars := map[string]interface{}{
		"query": githubv4.String("is:pr is:open user:gametimesf"),
	}

	// results gets mapped into `q`
	err = ghClient.Query(context.Background(), &q, vars)
	if err != nil {
		log.Error().Err(err).Msg("Failed to reach github")
		return
	}

	if len(q.Search.Edges) < 1 {
		return
	}

	// parse results and report to slack
	err = slackSvc.Report(q.Search.Edges)
	if err != nil {
		log.Error().Err(err).Msg("Failed to issue PullRequest report")
		return
	}
}

// helper function to figure environment
func getEnv() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}
	return env
}
