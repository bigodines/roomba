package main

import "github.com/shurcooL/githubv4"

type (
	// Search GraphQL query
	Search struct {
		Edges []Record
	}
	// Represents SearchResultItemConnection in Github API v4
	Record struct {
		Node struct {
			PullRequest struct {
				Author struct {
					Login string
				}
				Labels struct {
					Edges []struct {
						Node struct {
							Name string
						}
					}
				} `graphql:"labels(first:3)"`
				HeadRepository struct {
					Name string
				}
				UpdatedAt githubv4.DateTime
				Permalink string
				Title     string
			} `graphql:"... on PullRequest"`
		}
	}
)
