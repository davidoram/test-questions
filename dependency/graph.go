package dependency

import (
	"fmt"
	"slices"
)

// Build a graph structure starting with by reading a package.json.
// Then read each of its "dependencies". Use a resolver to do that job.
// The resolver is only called upon if dependency hasn't alrrady been read
// So if A depends on B, then A->B, and B depends on C then B->C, so the load
// order is  C, B, A
// If A depends on B, and B depends on C, and D, and D depends on C then
// A -> B -> C
//      +-> D -> C
// order is  C, D, B, A

type Graph struct {
	nodes []Node
}

func NewGraph() *Graph {
	return &Graph{nodes: []Node{}}
}

func (g *Graph) AddNode(name string) {
	if slices.ContainsFunc(g.nodes, func(v Node) bool { return v.name == name }) {
		return
	}
	g.nodes = append(g.nodes, NewNode(name))
}

func (g *Graph) AddEdge(from, to string) error {
	var fromIdx, toIdx int
	if fromIdx = slices.IndexFunc(g.nodes, func(v Node) bool { return v.name == from }); fromIdx == -1 {
		return fmt.Errorf("Graph doesn't contain %s node", from)
	}
	if toIdx = slices.IndexFunc(g.nodes, func(v Node) bool { return v.name == to }); toIdx == -1 {
		return fmt.Errorf("Graph doesn't contain %s node", to)
	}
	g.nodes[fromIdx].edges = append(g.nodes[fromIdx].edges, &g.nodes[toIdx])
	return nil
}

const (
	unvisited = iota
	visiting
	visited
)

func (g *Graph) IsDirectedAcyclicGraph() bool {

	// Turn graph into map string[string]
	graph := map[string][]string{}
	for _, node := range g.nodes {
		for _, neighbor := range node.edges {
			if graph[node.name] == nil {
				graph[node.name] = []string{}
			}
			graph[node.name] = append(graph[node.name], neighbor.name)
		}
	}

	// track state
	state := map[string]int{}

	var visit func(string) bool
	visit = func(node string) bool {
		if state[node] == visiting {
			return false
		}
		if state[node] == visited {
			return true
		}
		state[node] = visiting
		for _, next := range graph[node] {
			if !visit(next) {
				return false
			}
		}
		state[node] = visited
		return true
	}

	// Iterate over all nodes
	for node := range graph {
		if !visit(node) {
			return false
		}
	}
	return true
}

type Node struct {
	name string
	// neigbour is a directed vertex from this node
	edges []*Node
}

func NewNode(name string) Node {
	return Node{name: name, edges: []*Node{}}
}

type TopologicalSort struct {
	g       *Graph
	visited []string
	stack   []string
}

// depth first search
func (ts *TopologicalSort) dfs(v Node) {
	ts.visited = append(ts.visited, v.name)
	for _, neighbour := range v.edges {
		if !slices.Contains(ts.visited, neighbour.name) {
			ts.dfs(*neighbour)
		}
	}
	ts.stack = append(ts.stack, v.name)
}

func (g *Graph) TopologicalSort() ([]string, error) {
	if !g.IsDirectedAcyclicGraph() {
		return nil, fmt.Errorf("Not a DAG")
	}
	ts := TopologicalSort{g: g}
	for _, v := range ts.g.nodes {
		if !slices.Contains(ts.visited, v.name) {
			ts.dfs(v)
		}
	}
	return ts.stack, nil
}
