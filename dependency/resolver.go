package dependency

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
)

type FileResolver struct {
	path  string
	cache map[string]JSPackage
}

func NewFileResolver(path string) (FileResolver, error) {
	resolver := FileResolver{
		path:  path,
		cache: map[string]JSPackage{},
	}

	return resolver, nil
}

func (r FileResolver) Resolve(name string) (JSPackage, error) {
	if pkg, ok := r.cache[name]; ok {
		return pkg, nil
	}

	data, err := os.ReadFile(filepath.Join(r.path, name, "package.json"))
	if err != nil {
		return JSPackage{}, err
	}

	var pkg JSPackage
	if err := json.Unmarshal(data, &pkg); err != nil {
		return JSPackage{}, err
	}

	r.cache[name] = pkg
	return pkg, nil
}

func (r FileResolver) Packages() []JSPackage {
	packages := make([]JSPackage, 0, len(r.cache))
	for _, pkg := range r.cache {
		packages = append(packages, pkg)
	}

	slices.SortFunc(packages, func(a, b JSPackage) int {
		if a.Name() < b.Name() {
			return -1
		}
		if a.Name() > b.Name() {
			return 1
		}
		return 0
	})

	return packages
}

func MakeGraphFromJS(path string, names ...string) (*Graph, error) {
	resolver, err := NewFileResolver(path)
	if err != nil {
		return nil, err
	}

	graph := NewGraph()
	seen := map[string]bool{}

	var addPackage func(name string) error
	addPackage = func(name string) error {
		if seen[name] {
			return nil
		}
		seen[name] = true

		pkg, err := resolver.Resolve(name)
		if err != nil {
			return err
		}

		graph.AddNode(pkg.Name())
		for _, dependency := range pkg.DependencyNames() {
			if err := addPackage(dependency); err != nil {
				return err
			}
			if err := graph.AddEdge(pkg.Name(), dependency); err != nil {
				return err
			}
		}
		return nil
	}

	for _, name := range names {
		if err := addPackage(name); err != nil {
			return nil, err
		}
	}

	return graph, nil
}
