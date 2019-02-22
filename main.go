package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	flag.Parse()

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)

	httpClient := oauth2.NewClient(context.Background(), src)
	ghClient := githubv4.NewClient(httpClient)

	var q struct {
		Search struct {
			Edges []Record
		} `graphql:"search(query:$query, type:ISSUE, first:30)"`
	}
	vars := map[string]interface{}{
		"query": githubv4.String("is:pr is:open user:gametimesf"),
	}
	err := ghClient.Query(context.Background(), &q, vars)
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	// TODO: slack results
	printJSON(q)
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
