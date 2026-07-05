package dependency

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGraph(t *testing.T) {
	g := NewGraph()
	require.Empty(t, g.nodes)

	g.AddNode("a")
	g.AddNode("b")
	g.AddEdge("a", "b")
}

func TestIsDirectedAcyclicGraph(t *testing.T) {
	g := NewGraph()
	require.Empty(t, g.nodes)

	g.AddNode("a")
	g.AddNode("b")
	g.AddEdge("a", "b")

	require.True(t, g.IsDirectedAcyclicGraph())

	g.AddEdge("b", "a")
	require.False(t, g.IsDirectedAcyclicGraph())
}

func TestIsDirectedAcyclicGraph_self_loop(t *testing.T) {
	g := NewGraph()
	g.AddNode("a")
	g.AddEdge("a", "a")

	require.False(t, g.IsDirectedAcyclicGraph())
}

func TestTopologicalSort_straight(t *testing.T) {
	g := NewGraph()
	g.AddNode("a")
	g.AddNode("b")
	g.AddNode("c")
	g.AddEdge("a", "b")
	g.AddEdge("b", "c")
	nodes, err := g.TopologicalSort()
	require.Nil(t, err)
	require.Equal(
		t,
		[]string{"c", "b", "a"},
		nodes,
	)
}

func TestTopologicalSort_one(t *testing.T) {
	g := NewGraph()
	g.AddNode("a")
	nodes, err := g.TopologicalSort()
	require.Nil(t, err)
	require.Equal(
		t,
		[]string{"a"},
		nodes,
	)
}

func TestTopologicalSort_two(t *testing.T) {
	g := NewGraph()
	g.AddNode("a")
	g.AddNode("b")
	nodes, err := g.TopologicalSort()
	require.Nil(t, err)
	require.Equal(
		t,
		[]string{"a", "b"},
		nodes,
	)
}

func TestTopologicalSort_branch(t *testing.T) {
	g := NewGraph()
	g.AddNode("a")
	g.AddNode("b")
	g.AddNode("c")
	g.AddNode("d")

	g.AddEdge("a", "b")
	g.AddEdge("a", "c")
	g.AddEdge("c", "d")
	nodes, err := g.TopologicalSort()
	require.Nil(t, err)
	require.Equal(
		t,
		[]string{"b", "d", "c", "a"},
		nodes,
	)
}

func TestTopologicalSort_branchjoin(t *testing.T) {
	g := NewGraph()
	g.AddNode("a")
	g.AddNode("b")
	g.AddNode("c")
	g.AddNode("d")

	g.AddEdge("a", "b")
	g.AddEdge("a", "c")
	g.AddEdge("c", "d")
	g.AddEdge("b", "d")
	nodes, err := g.TopologicalSort()
	require.Nil(t, err)
	require.Equal(
		t,
		[]string{"d", "b", "c", "a"},
		nodes,
	)
}
