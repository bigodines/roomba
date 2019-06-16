package roomba

import (
	"strings"
	"time"
)

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
		Edges []LabelEdge
	}

	LabelEdge struct {
		Node LabelNode
	}

	LabelNode struct {
		Name string
	}
)

// PrintableLabels Converts a list of labels into a printable string
func PrintableLabels(labels Labels) string {
	ll := make([]string, 0)

	if len(labels.Edges) > 0 {
		for _, edge := range labels.Edges {
			n := edge.Node.Name
			if len(n) > 0 {
				ll = append(ll, n)
			}
		}
	}

	return strings.Join(ll[:], ", ")
}
