package main

import "time"

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
				Labels         Labels `graphql:"labels(first:3)"`
				HeadRepository struct {
					Name string
				}
				UpdatedAt time.Time
				Permalink string
				Title     string
			} `graphql:"... on PullRequest"`
		}
	}
	Labels struct {
		Edges []struct {
			Node struct {
				Name string
			}
		}
	}
)

func PrintableLabels(labels Labels) string {
	return ""
}
