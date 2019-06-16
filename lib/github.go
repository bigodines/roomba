package roomba

import (
	"regexp"
	"strings"
	"time"
)

type (
	// Search GraphQL query
	Search struct {
		Edges []Record
	}
	Node struct{
		Name string
	}
	//Represents a Repository in Github API v4
	Repository struct {
		Labels struct {
			Nodes []Node
		} `graphql:"labels(first:20)"`
	}


	// Represents SearchResultItemConnection in Github API v4
	Record struct {
		Node struct {
			PullRequest struct {
				Author struct {
					Login string
				}
				Labels         Labels `graphql:"labels(first:10)"`
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

// GetLabelTypeSet Converts a list of labels into a set of LabelType
func GetLabelTypeSetFromRepository(repo Repository) LabelTypeSet {
	labelTypeSet := make(LabelTypeSet)
	if len(repo.Labels.Nodes) > 0 {
		for _, node := range repo.Labels.Nodes {
			labelTypeSet [GetLabelType(node.Name)] = true
		}
	}

	return labelTypeSet
}

// GetLabelTypeSet Converts a list of labels into a set of LabelType
func GetLabelTypeSet(labels Labels) LabelTypeSet {
	labelTypeSet := make(LabelTypeSet)

	if len(labels.Edges) > 0 {
		for _, edge := range labels.Edges {
			labelTypeSet [GetLabelType(edge.Node.Name)] = true
		}
	}

	return labelTypeSet
}

type LabelTypeSet map[LabelType]bool

type LabelType string

const (
	UNKNOWN LabelType = "UNKNOWN"
	NeedsOneReview = "NEEDS ONE REVIEW"
	NeedsTwoReviews = "NEEDS TWO REVIEWS"
	READY = "READY"
	WIP ="WIP"
)

var oneReviewRegex = regexp.MustCompile(`(?i)needs 1 review`)
var twoReviewRegex = regexp.MustCompile(`(?i)needs 2 reviews`)
var readyRegex = regexp.MustCompile(`(?i)ready`)
var wipRegex = regexp.MustCompile(`(?i)wip`)

// GetLabelType returns label type enum
func GetLabelType(label string) LabelType {
	switch {
	case oneReviewRegex.MatchString(label):
		return LabelType(NeedsOneReview)
	case twoReviewRegex.MatchString(label):
		return LabelType(NeedsTwoReviews)
	case readyRegex.MatchString(label):
		return LabelType(READY)
	case wipRegex.MatchString(label):
		return LabelType(WIP)
	default:
		return LabelType(UNKNOWN)
	}
}

// HasValidLabel returns whether or not a given LabelTypeSet has any valid labels in it
func HasValidLabel(set LabelTypeSet) bool {
	if _, ok := set[NeedsOneReview]; ok {
		return true
	} else if _, ok := set[NeedsTwoReviews]; ok {
		return true
	} else if _, ok := set[READY]; ok {
		return true
	} else if _, ok := set[WIP]; ok {
		return true
	} else {
		return false
	}
}
