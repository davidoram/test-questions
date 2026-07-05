package dependency

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFileResolver(t *testing.T) {
	resolver, err := NewFileResolver("data")
	require.NoError(t, err)

	require.Equal(t, "data", resolver.path)
	require.Empty(t, resolver.Packages())
}

func TestFileResolverResolve(t *testing.T) {
	resolver, err := NewFileResolver("data")
	require.NoError(t, err)

	pkg, err := resolver.Resolve("app-basic")
	require.NoError(t, err)

	require.Equal(t, "app-basic", pkg.Name())
	require.Equal(t, []string{"alpha", "beta"}, pkg.DependencyNames())
	require.Equal(t, []JSPackage{pkg}, resolver.Packages())
}

func TestFileResolverResolveUsesCache(t *testing.T) {
	resolver, err := NewFileResolver("data")
	require.NoError(t, err)

	pkg, err := resolver.Resolve("alpha")
	require.NoError(t, err)

	resolver.path = filepath.Join("data", "missing")
	cachedPkg, err := resolver.Resolve("alpha")
	require.NoError(t, err)
	require.Equal(t, pkg, cachedPkg)
}

func TestFileResolverResolveMissingPackage(t *testing.T) {
	resolver, err := NewFileResolver("data")
	require.NoError(t, err)

	pkg, err := resolver.Resolve("missing")

	require.Error(t, err)
	require.Empty(t, pkg)
}

func TestMakeGraphFromJS(t *testing.T) {
	tests := []struct {
		start        string
		isDAG        bool
		dependencies []string
	}{
		{
			start:        "app-basic",
			isDAG:        true,
			dependencies: []string{"delta", "charlie", "alpha", "beta", "app-basic"},
		},
		{
			start:        "app-shared",
			isDAG:        true,
			dependencies: []string{"delta", "charlie", "alpha", "echo", "app-shared"},
		},
		{
			start:        "app-dev",
			isDAG:        true,
			dependencies: []string{"delta", "charlie", "alpha", "tools", "app-dev"},
		},
		{
			start: "app-cycle",
			isDAG: false,
		},
	}

	for _, test := range tests {
		t.Run(test.start, func(t *testing.T) {
			graph, err := MakeGraphFromJS("data", test.start)
			require.NoError(t, err)
			require.Equal(t, test.isDAG, graph.IsDirectedAcyclicGraph())

			dependencies, err := graph.TopologicalSort()
			if !test.isDAG {
				require.Error(t, err)
				require.Nil(t, dependencies)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.dependencies, dependencies)
		})
	}
}

func TestMakeGraphFromJSMissingPath(t *testing.T) {
	graph, err := MakeGraphFromJS(filepath.Join("data", "missing"), "app-basic")

	require.Error(t, err)
	require.Nil(t, graph)
}
