package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/bigodines/roomgo/config"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	// TODO: support flags
	flag.Parse()

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
		panic(err)
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
		fmt.Printf("%+v", err)
		return
	}

	if len(q.Search.Edges) < 1 {
		return
	}

	// parse results and report to slack
	err = slackSvc.Report(q.Search.Edges)
	if err != nil {
		fmt.Printf("Failed to issue PullRequest report: (%s)\n", err)
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

// TODO: remove
// printJSON prints v as JSON encoded with indent to stdout. It panics on any error.
func printJSON(v interface{}) {
	w := json.NewEncoder(os.Stdout)
	w.SetIndent("", "   ")
	err := w.Encode(v)
	if err != nil {
		panic(err)
	}
}
