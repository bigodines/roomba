package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintableLabels(t *testing.T) {
	labels := NewTestLabel("foo", "bar")
	assert.Equal(t, "foo, bar", PrintableLabels(labels))

	labels = NewTestLabel()
	assert.Equal(t, "", PrintableLabels(labels))
}

func NewTestLabel(names ...string) Labels {
	edges := make([]LabelEdge, 0)
	for _, name := range names {
		edges = append(edges, LabelEdge{
			Node: LabelNode{
				Name: name,
			},
		})
	}
	return Labels{Edges: edges}
}
