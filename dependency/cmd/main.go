package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/davidoram/test-questions/dependency"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "merge: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout *os.File) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: merge path file")
	}
	path := args[0]
	file := args[1]

	graph, err := dependency.MakeGraphFromJS(path, file)
	if err != nil {
		return err
	}
	if !graph.IsDirectedAcyclicGraph() {
		return fmt.Errorf("not a DAG")
	}

	sorted, err := graph.TopologicalSort()
	if err != nil {
		return err
	}

	fmt.Fprintf(stdout, "Build order for '%s' is \n - %s\n", file, strings.Join(sorted, "\n - "))

	return nil
}
