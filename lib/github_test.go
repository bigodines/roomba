package roomba

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

func TestGetLabelType(t *testing.T) {
	assert.Equal(t, LabelType(NeedsOneReview), GetLabelType("needs 1 review"))
	assert.Equal(t, LabelType(NeedsOneReview), GetLabelType("NEEDS 1 REVIEW"))
	assert.Equal(t, LabelType(NeedsOneReview), GetLabelType("NeeDS 1 rEviEW"))

	assert.Equal(t, LabelType(NeedsTwoReviews), GetLabelType("needs 2 reviews"))
	assert.Equal(t, LabelType(NeedsTwoReviews), GetLabelType("NEEDS 2 REVIEWS"))
	assert.Equal(t, LabelType(NeedsTwoReviews), GetLabelType("NeeDS 2 rEviEWs"))

	assert.Equal(t, LabelType(READY), GetLabelType("ready"))
	assert.Equal(t, LabelType(READY), GetLabelType("READY"))
	assert.Equal(t, LabelType(READY), GetLabelType("ReAdY"))

	assert.Equal(t, LabelType(WIP), GetLabelType("WIP"))
	assert.Equal(t, LabelType(WIP), GetLabelType("wip"))
	assert.Equal(t, LabelType(WIP), GetLabelType("wIp"))

	assert.Equal(t, LabelType(UNKNOWN), GetLabelType("Prince"))
	assert.Equal(t, LabelType(UNKNOWN), GetLabelType("of"))
	assert.Equal(t, LabelType(UNKNOWN), GetLabelType("Egypt"))
}

func TestGetLabelTypeSetAllValid(t *testing.T) {
	stringArray := StringArrayToLabelsArrayHelper([]string{"needs 1 review", "needs 2 reviews", "ready", "wip"})
	assert.Equal(t, GetLabelTypeSet(stringArray),
		LabelTypeSet{NeedsOneReview: true, NeedsTwoReviews: true, READY: true, WIP: true})
}

func TestGetLabelTypeSetEmpty(t *testing.T) {
	stringArray := StringArrayToLabelsArrayHelper([]string{})
	assert.Equal(t, GetLabelTypeSet(stringArray), LabelTypeSet{})
}

func TestGetLabelTypeSetAllInvalid(t *testing.T) {
	stringArray := StringArrayToLabelsArrayHelper([]string{"prince", "of", "egypt"})
	assert.Equal(t, GetLabelTypeSet(stringArray), LabelTypeSet{UNKNOWN: true})
}

func TestGetLabelTypeSetMixed(t *testing.T) {
	stringArray := StringArrayToLabelsArrayHelper([]string{"needs 1 review", "might need stuff", "ready", "best PR"})
	assert.Equal(t, GetLabelTypeSet(stringArray), LabelTypeSet{NeedsOneReview: true, UNKNOWN: true, READY: true})
}

func TestHasValidLabel(t *testing.T) {
	assert.True(t, HasValidLabel(LabelTypeSet{WIP: true}))
	assert.True(t, HasValidLabel(LabelTypeSet{NeedsOneReview: true}))
	assert.True(t, HasValidLabel(LabelTypeSet{NeedsTwoReviews: true}))
	assert.True(t, HasValidLabel(LabelTypeSet{READY: true}))
	assert.True(t, HasValidLabel(LabelTypeSet{NeedsOneReview: true, UNKNOWN: true, READY: true}))
	assert.True(t, HasValidLabel(LabelTypeSet{NeedsOneReview: true, UNKNOWN: true}))
	assert.True(t, HasValidLabel(LabelTypeSet{UNKNOWN: true, READY: true}))
	assert.True(t, HasValidLabel(LabelTypeSet{NeedsOneReview: true, UNKNOWN: true, READY: true}))
	assert.False(t, HasValidLabel(LabelTypeSet{}))
	assert.False(t, HasValidLabel(LabelTypeSet{UNKNOWN: true}))
}
func TestGetLabelTypeSetFromRepositoryAllValid(t *testing.T) {
	repository := StringArrayToRepositoryHelper([]string{"needs 1 review", "needs 2 reviews", "ready", "wip"})
	assert.Equal(t, GetLabelTypeSetFromRepository(repository),
		LabelTypeSet{NeedsOneReview: true, NeedsTwoReviews: true, READY: true, WIP: true})
}

func TestGetLabelTypeSetFromRepositoryEmpty(t *testing.T) {
	repository := StringArrayToRepositoryHelper([]string{})
	assert.Equal(t, GetLabelTypeSetFromRepository(repository), LabelTypeSet{})
}

func TestGetLabelTypeSetFromRepositoryAllInvalid(t *testing.T) {
	repository := StringArrayToRepositoryHelper([]string{"prince", "of", "egypt"})
	assert.Equal(t, GetLabelTypeSetFromRepository(repository), LabelTypeSet{UNKNOWN: true})
}

func TestGetLabelTypeSetFromRepositoryMixed(t *testing.T) {
	repository := StringArrayToRepositoryHelper([]string{"needs 1 review", "might need stuff", "ready", "best PR"})
	assert.Equal(t, GetLabelTypeSetFromRepository(repository), LabelTypeSet{NeedsOneReview: true, UNKNOWN: true, READY: true})
}

func StringArrayToLabelsArrayHelper(stringArray []string) Labels {
	labelEdges := make([]LabelEdge, 0)
	for _, stringValue := range stringArray {
		labelEdges = append(labelEdges, LabelEdge{Node: LabelNode{Name: stringValue,}})
	}
	return Labels{Edges: labelEdges}
}

func StringArrayToRepositoryHelper(stringArray []string) Repository {
	nodesArray := make([]Node, 0)
	for _, stringValue := range stringArray {
		x := Node{Name:stringValue}
		nodesArray = append(nodesArray,x)
	}
	return Repository{Labels: struct{ Nodes []Node }{Nodes: nodesArray}}
}
