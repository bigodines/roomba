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
	client := githubv4.NewClient(httpClient)

	{
		type Record struct {
			Node struct {
				PullRequest struct {
					Title string
				} `graphql:"... on PullRequest"`
			}
		}
		var q struct {
			Search struct {
				Edges []Record
			} `graphql:"search(query:$q, type:ISSUE, first:30)"`
		}
		variables := map[string]interface{}{
			"q": githubv4.String("user:gametimesf"),
		}
		err := client.Query(context.Background(), &q, variables)
		if err != nil {
			fmt.Printf("%+v", err)
			return
		}
		printJSON(q)
	}
}

// printJSON prints v as JSON encoded with indent to stdout. It panics on any error.
func printJSON(v interface{}) {
	w := json.NewEncoder(os.Stdout)
	w.SetIndent("", "   ")
	err := w.Encode(v)
	if err != nil {
		panic(err)
	}
}
